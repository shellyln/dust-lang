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
func formalArgument() ParserFn {
	return Trans(
		FlatGroup(
			symbol(mnem.Symbol), sp0(),
			erase(Seq(":")), sp0(),
			typeExpr(),
		),
		TransFormalArgument,
	)
}

// For two-pass parsing of function definitions
func formalArgumentForPreScan() ParserFn {
	return Trans(
		FlatGroup(
			symbol(mnem.Symbol), sp0(),
			erase(Seq(":")), sp0(),
			typeExpr(),
		),
		TransFormalArgumentForPreScan,
	)
}

//
func formalArgumentForTyping() ParserFn {
	return Trans(
		FlatGroup(
			Zero(Ast{ClassName: clsz.Placeholder, Value: ""}),
			typeExpr(),
		),
		TransFormalArgumentForTyping,
	)
}

//
func lambdaExpressionInner(typing bool) ParserFn {
	return FlatGroup(
		erase(CharClass("|")), sp0(),
		Zero(Ast{ClassName: clsz.Placeholder}),
		Group(First(
			FlatGroup(
				erase(CharClass("|")), sp0(),
			),
			FlatGroup(
				// TODO: spread
				If(typing,
					formalArgumentForTyping(),
					formalArgument(),
				),
				ZeroOrMoreTimes(
					erase(CharClass(",")), sp0(),
					// TODO: spread
					If(typing,
						formalArgumentForTyping(),
						formalArgument(),
					),
				),
				ZeroOrOnce(
					erase(CharClass(",")), sp0(),
				),
				erase(CharClass("|")), sp0(),
			),
			If(typing,
				Unmatched(),
				Error(emsg.ParserErr00012),
			),
		)),
		First(
			FlatGroup(
				erase(Seq("->")), sp0(),
				typeExpr(),
			),
			Zero(Ast{ClassName: clsz.Placeholder}),
		),
	)
}

//
func lambdaExpression() ParserFn {
	return lambdaScope(Trans(
		FlatGroup(
			Trans(
				lambdaExpressionInner(false),
				TransDefineLambdaSymbols,
			),
			mustBracketMultipleStatement(),
		),
		TransLambdaExpression,
	))
}

//
func functionDeclarationInner(preScan bool) ParserFn {
	return FlatGroup(
		erase(Seq("fn")), sp1(),
		symbol(mnem.Symbol), sp0(),
		erase(CharClass("(")), sp0(),
		Group(First(
			FlatGroup(
				erase(CharClass(")")), sp0(),
			),
			FlatGroup(
				// TODO: spread
				If(preScan,
					formalArgumentForPreScan(),
					formalArgument(),
				),
				ZeroOrMoreTimes(
					erase(CharClass(",")), sp0(),
					// TODO: spread
					If(preScan,
						formalArgumentForPreScan(),
						formalArgument(),
					),
				),
				ZeroOrOnce(
					erase(CharClass(",")), sp0(),
				),
				erase(CharClass(")")), sp0(),
			),
			Error(emsg.ParserErr00013),
		)),
		First(
			FlatGroup(
				erase(Seq("->")), sp0(),
				typeExpr(),
			),
			Zero(Ast{ClassName: clsz.Placeholder}),
		),
	)
}

//
func functionDefinition() ParserFn {
	return funcScope(Trans(
		FlatGroup(
			Trans(
				functionDeclarationInner(false),
				TransDefineFuncSymbols,
			),
			mustBracketMultipleStatement(),
		),
		TransFuncStatement,
	))
}
