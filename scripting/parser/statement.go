package parser

import (
	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/dust-lang/scripting/rules"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func functionDefinitionScannerChild() ParserFn {
	return First(
		sp1(),
		FlatGroup(
			erase(CharClass("{")), sp0(),
			refFunctionDefinitionScannerChild(),
			erase(CharClass("}")), sp0(),
		),
		Any(),
	)
}

//
func refFunctionDefinitionScannerChild() ParserFn {
	return ref(functionDefinitionScannerChild)
}

// Two-pass parsing of function definitions
func functionDefinitionScanner() ParserFn {
	return LookAhead(
		ZeroOrMoreTimes(
			First(
				sp1(),
				FlatGroup(
					erase(CharClass("{")), sp0(),
					functionDefinitionScannerChild(),
					erase(CharClass("}")), sp0(),
				),
				Trans(
					functionDeclarationInner(true),
					TransPreScanFuncDefinitions,
				),
				erase(Any()),
			),
		),
	)
}

//
func refFunctionDefinitionScanner() ParserFn {
	return ref(functionDefinitionScanner)
}

//
func multipleStatementImpl(top bool) ParserFn {
	return Trans(
		FlatGroup(
			If(top,
				sp0(),
				FlatGroup(erase(CharClass("{")), sp0()),
			),
			// NOTE: Two-pass parsing of function definitions
			functionDefinitionScanner(),
			Group(
				First(
					FlatGroup(
						OneOrMoreTimes(
							refNestedStatement(),
						),
						If(top,
							LookAhead(End()),
							FlatGroup(erase(CharClass("}")), sp0()),
						),
					),
					Error(emsg.ParserErr00010),
				),
			),
		),
		TransMultipleStatement,
	)
}

//
func multipleStatement(top bool) ParserFn {
	return scope(multipleStatementImpl(top))
}

//
func mustBracketMultipleStatementImpl() ParserFn {
	return Trans(
		FlatGroup(
			erase(CharClass("{")), sp0(),
			// NOTE: Two-pass parsing of function definitions.
			//       At the present, functions can only be defined at the top level.
			// functionDefinitionScanner(),
			Group(
				First(
					FlatGroup(
						ZeroOrMoreTimes(
							refNestedStatement(),
						),
						erase(CharClass("}")), sp0(),
					),
					Error(emsg.ParserErr00010),
				),
			),
		),
		TransMultipleStatement,
	)
}

//
func mustBracketMultipleStatement() ParserFn {
	return scope(mustBracketMultipleStatementImpl())
}

//
func returnStatement() ParserFn {
	return Trans(
		FlatGroup(
			erase(Seq("return")),
			WordBoundary(), sp0(),
			First(
				expression(true),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
		),
		TransReturnStatement,
	)
}

//
func breakStatement() ParserFn {
	return Trans(
		FlatGroup(
			erase(Seq("break")),
			WordBoundary(), sp0(),
			First(
				FlatGroup(
					LookAheadN(Seq("returning"), WordBoundary()),
					erase(CharClass("'")),
					symbol(mnem.Symbol),
					WordBoundary(), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			First(
				FlatGroup(
					erase(Seq("returning")),
					WordBoundary(), sp0(),
					expression(true),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
		),
		TransBreakStatement,
	)
}

//
func continueStatement() ParserFn {
	return Trans(
		FlatGroup(
			erase(Seq("continue")),
			WordBoundary(), sp0(),
			First(
				FlatGroup(
					erase(CharClass("'")),
					symbol(mnem.Symbol),
					WordBoundary(), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
		),
		TransContinueStatement,
	)
}

// TODO:
func pubDefStatement() ParserFn {
	return FlatGroup(
		erase(Seq("pub")), sp1(),
		defExpression(),
	)
}

// TODO:
func pubFunctionDefinition() ParserFn {
	return FlatGroup(
		erase(Seq("pub")), sp1(),
		functionDefinition(),
	)
}

//
func nestedStatement() ParserFn {
	return FlatGroup(
		First(
			FlatGroup(
				First(
					// TODO: declare statement
					returnStatement(),
					breakStatement(),
					continueStatement(),
					doWhileExpression(),
					pubDefStatement(),
					defExpression(),
					expression(true),
				),
				LookAhead(
					First(
						CharClass(";", "}"),
						End(),
					),
				),
			),
			// NOTE: if/while/for expression as a statement
			ifExpression(),
			loopExpression(),
			whileExpression(),
			forInExpression(),
			tradForExpression(),
			pubFunctionDefinition(),
			functionDefinition(),
			multipleStatement(false),
		),
		ZeroOrMoreTimes(
			erase(CharClass(";")), sp0(),
		),
	)
}

//
func refNestedStatement() ParserFn {
	return ref(nestedStatement)
}

//
func topLevelStatement() ParserFn {
	return multipleStatement(true)
}
