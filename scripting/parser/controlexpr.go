package parser

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/dust-lang/scripting/rules"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func ifExpression() ParserFn {
	return scope(Trans(
		Group(
			FlatGroup(
				erase(Seq("if")),
				sp1(),
				Zero(Ast{
					ClassName: clsz.IfCondition,
					Type:      AstType_String,
					Value:     "if",
				}),
				First(
					FlatGroup(
						defExpression(),
						erase(Seq(";")), sp0(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				// TODO: add scope v
				refExpression(true),
				mustBracketMultipleStatement(),
				// TODO: add scope ^
			),
			ZeroOrMoreTimes(
				erase(Seq("else")),
				sp1(),
				erase(Seq("if")),
				sp1(),
				Zero(Ast{
					ClassName: clsz.IfCondition,
					Type:      AstType_String,
					Value:     "elseif",
				}),
				First(
					FlatGroup(
						defExpression(),
						erase(Seq(";")), sp0(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				// TODO: add scope v
				refExpression(true),
				mustBracketMultipleStatement(),
				// TODO: add scope ^
			),
			First(
				FlatGroup(
					erase(Seq("else")),
					sp0(),
					Zero(Ast{
						ClassName: clsz.IfConditionElse,
						Type:      AstType_String,
						Value:     "else",
					}),
					Zero(Ast{ClassName: clsz.Placeholder}),
					Zero(Ast{ClassName: clsz.Placeholder}),
					mustBracketMultipleStatement(),
				),
				FlatGroup(
					Zero(Ast{
						ClassName: clsz.IfConditionNop,
						Type:      AstType_String,
						Value:     "nop",
					}),
					Zero(Ast{ClassName: clsz.Placeholder}),
					Zero(Ast{ClassName: clsz.Placeholder}),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
			),
		),
		TransIfExpression,
	))
}

//
func loopExpression() ParserFn {
	return scope(Trans(
		FlatGroup(
			First(
				FlatGroup(
					erase(Seq("'")),
					symbol(mnem.Imm_str), sp0(),
					erase(Seq(":")), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			Group(
				erase(Seq("loop")),
				WordBoundary(), sp0(),
				First(
					FlatGroup(
						defExpression(),
						erase(Seq(";")), sp0(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				// TODO: add scope v
				mustBracketMultipleStatement(),
				// TODO: add scope ^
			),
		),
		TransLoopExpression,
	))
}

//
func whileExpression() ParserFn {
	return scope(Trans(
		FlatGroup(
			First(
				FlatGroup(
					erase(Seq("'")),
					symbol(mnem.Imm_str), sp0(),
					erase(Seq(":")), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			Group(
				erase(Seq("while")),
				WordBoundary(), sp0(),
				First(
					FlatGroup(
						defExpression(),
						erase(Seq(";")), sp0(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				// TODO: add scope v
				refExpression(true),
				mustBracketMultipleStatement(),
				// TODO: add scope ^
			),
		),
		TransWhileExpression,
	))
}

//
func doWhileExpression() ParserFn {
	return scope(Trans(
		FlatGroup(
			First(
				FlatGroup(
					erase(Seq("'")),
					symbol(mnem.Imm_str), sp0(),
					erase(Seq(":")), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			Group(
				erase(Seq("do")),
				WordBoundary(), sp0(),
				First(
					FlatGroup(
						defExpression(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				mustBracketMultipleStatement(),
				erase(Seq("while")),
				WordBoundary(), sp0(),
				// TODO: add scope v
				refExpression(true),
				// TODO: add scope ^
			),
		),
		TransDoWhileExpression,
	))
}

//
func forInExpression() ParserFn {
	return scope(Trans(
		FlatGroup(
			First(
				FlatGroup(
					erase(Seq("'")),
					symbol(mnem.Imm_str), sp0(),
					erase(Seq(":")), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			Group(
				Trans(
					FlatGroup(
						erase(Seq("for")),
						WordBoundary(), sp0(),
						First(
							FlatGroup(
								defExpression(),
								erase(Seq(";")), sp0(),
							),
							Zero(Ast{ClassName: clsz.Placeholder}),
						),
						symbol(mnem.Symbol), sp0(),
						WordBoundary(), sp0(),
						erase(Seq("in")),
						WordBoundary(), sp0(),
						First(
							FlatGroup(
								Zero(Ast{ClassName: clsz.Range}),
								refExpression(true),
								erase(CharClass("..")), sp0(),
								refExpression(true),
							),
							FlatGroup(
								Zero(Ast{ClassName: clsz.Placeholder}),
								refExpression(true),
								Zero(Ast{ClassName: clsz.Placeholder}),
							),
						),
					),
					TransDefForInSymbol,
				),
				mustBracketMultipleStatement(),
			),
		),
		TransForInExpression,
	))
}

//
func tradForExpression() ParserFn {
	return scope(Trans(
		FlatGroup(
			First(
				FlatGroup(
					erase(Seq("'")),
					symbol(mnem.Imm_str), sp0(),
					erase(Seq(":")), sp0(),
				),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
			Group(
				erase(Seq("for")),
				WordBoundary(), sp0(),
				First(
					defExpression(),
					refExpression(true),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				erase(Seq(";")), sp0(),
				First(
					refExpression(true),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				erase(Seq(";")), sp0(),
				First(
					refExpression(true),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				mustBracketMultipleStatement(),
			),
		),
		TransTradForExpression,
	))
}

// TODO: match {} statement
