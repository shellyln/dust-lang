package parser

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/dust-lang/scripting/rules"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func typeObjectMember() ParserFn {
	return FlatGroup(
		symbol(mnem.Symbol), WordBoundary(), sp0(),
		erase(CharClass(":")), sp0(),
		refTypeExpr(),
	)
}

//
func typeObjectExpr() ParserFn {
	return Trans(
		Group(
			erase(CharClass("{")), sp0(),
			Group(First(
				FlatGroup(
					erase(CharClass("}")), sp0(),
				),
				FlatGroup(
					typeObjectMember(),
					ZeroOrMoreTimes(
						erase(CharClass(",")), sp0(),
						typeObjectMember(),
					),
					erase(CharClass("}")), sp0(),
				),
			)),
		),
		TransTypeObjectExpr,
	)
}

//
func typeFuncExpr() ParserFn {
	return Trans(
		lambdaExpressionInner(true),
		TransTypeFuncExpr,
	)
}

//
func typeExpr() ParserFn {
	return First(
		Trans(
			FlatGroup(
				erase(CharClass("[")), sp0(),
				refTypeExpr(),
				erase(CharClass("]")), sp0(),
			),
			TransTypeSliceExpr,
		),
		ref(typeFuncExpr),
		Trans(
			FlatGroup(
				symbol(mnem.Symbol), WordBoundary(), sp0(),
				Group(
					erase(CharClass("<")), sp0(),
					refTypeExpr(),
					ZeroOrMoreTimes(
						erase(CharClass(",")), sp0(),
						refTypeExpr(),
					),
					ZeroOrOnce(
						erase(CharClass(",")), sp0(),
					),
					erase(CharClass(">")), sp0(),
				),
			),
			TransTypeGenericsExpr,
		),
		typeObjectExpr(),
		Trans(
			FlatGroup(
				symbol(mnem.Symbol), WordBoundary(), sp0(),
			),
			TransTypeExpr,
		),
	)
}

//
func refTypeExpr() ParserFn {
	return ref(typeExpr)
}
