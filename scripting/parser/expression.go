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
func unaryPostfixOp() ParserFn {
	return Trans(
		FlatGroup(CharClass("++", "--")),
		ChangeClassName(clsz.UnaryPostfixOp),
	)
}

//
func unaryPrefixOp() ParserFn {
	return Trans(
		FlatGroup(CharClass(
			"++", "--", "!", "~", "+", "-",
			// TODO: "&"
		), LookAheadN(Number())),
		ChangeClassName(clsz.UnaryPrefixOp),
	)
}

//
func binaryOp() ParserFn {
	return Trans(
		First(
			FlatGroup(CharClass("/"), LookAheadN(CharClass("/"))),
			CharClass(
				// NOTE: The operators are ordered by their string length,
				//       not by operator precedence.
				"**", "*", "%", "+", "-",
				"<<", ">>", ">>>",
				"<=", "<", ">=", ">",
				"===", "!==", "==", "!=",
				"|>",
				"&&", "||",
				"&", "^", "|",
				"=",
			),
		),
		ChangeClassName(clsz.BinaryOp),
	)
}

//
func binaryCommaOp() ParserFn {
	return Trans(
		CharClass(","),
		ChangeClassName(clsz.BinaryOp),
	)
}

//
func binaryTypeCastAsOp() ParserFn {
	return Trans(
		FlatGroup(
			WordBoundary(),
			Seq("as"),
			WordBoundary(),
			LookAhead(sp0(), typeExpr()),
		),
		ChangeClassName(clsz.TypeCast),
	)
}

//
func binaryDotOp() ParserFn {
	return Trans(
		FlatGroup(
			CharClass("."),
			LookAhead(sp0(), symbol(mnem.Symbol)),
		),
		ChangeClassName(clsz.Access),
	)
}

//
func ternaryTailExpr(enableComma bool) ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.TernaryOp,
			Type:      AstType_String,
			Value:     "?:",
		}),
		erase(CharClass("?")), sp0(), refExpression(enableComma),
		erase(CharClass(":")), sp0(), refExpression(enableComma),
	)
}

//
func binaryTailExpr(enableComma bool) ParserFn {
	return If(enableComma,
		First(
			FlatGroup(
				First(
					binaryOp(), binaryCommaOp(),
					binaryDotOp(),
				), sp0(),
				refExpressionImpl(true),
			),
			FlatGroup(
				binaryTypeCastAsOp(), sp0(),
				refTypeExpressionImpl(true),
			),
		),
		First(
			FlatGroup(
				First(
					binaryOp(),
					binaryDotOp(),
				), sp0(),
				refExpressionImpl(false),
			),
			FlatGroup(
				binaryTypeCastAsOp(), sp0(),
				refTypeExpressionImpl(false),
			),
		),
	)
}

//
func expressionImpl(enableComma bool) ParserFn {
	return FlatGroup(
		LookAheadN(First(End(), CharClass(")", "]", "}", "?", ":", ",", ";"))),
		First(
			FlatGroup(
				First(
					ifExpression(),
					loopExpression(),
					whileExpression(),
					doWhileExpression(),
					forInExpression(),
					tradForExpression(),
					lambdaExpression(),

					groupingParenthesesExpr(),
					FlatGroup(literalOrSymbol(), sp0()),

					mustBracketMultipleStatement(),
				),
				ZeroOrMoreTimes(First(
					actualArgumentsExpr(),
					acceessExpr(),
					sliceExpr(),
					FlatGroup(unaryPostfixOp(), sp0()),
				)),
				ZeroOrOnce(First(
					ternaryTailExpr(enableComma),
					binaryTailExpr(enableComma),
				)),
			),
			FlatGroup(
				unaryPrefixOp(), sp0(),
				refExpressionImpl(enableComma),
			),
			FlatGroup(LookAhead(
				First(
					FlatGroup(
						CharClassN("{"),
						Error(emsg.ParserErr00002),
					),
					Unmatched(),
				),
			)),
		),
	)
}

//
func expressionImplEnableComma() ParserFn {
	return expressionImpl(true)
}

//
func expressionImplDisableComma() ParserFn {
	return expressionImpl(false)
}

//
func refExpressionImpl(enableComma bool) ParserFn {
	return If(enableComma,
		ref(expressionImplEnableComma),
		ref(expressionImplDisableComma),
	)
}

//
func expression(enableComma bool) ParserFn {
	return Trans(
		expressionImpl(enableComma),
		ProdRule,
	)
}

//
func expressionEnableComma() ParserFn {
	return expression(true)
}

//
func expressionDisableComma() ParserFn {
	return expression(false)
}

//
func refExpression(enableComma bool) ParserFn {
	return If(enableComma,
		ref(expressionEnableComma),
		ref(expressionDisableComma),
	)
}

// TODO:
func typeExpressionImpl(enableComma bool) ParserFn {
	return FlatGroup(
		LookAheadN(First(End(), CharClass(")", "]", "}", "?", ":", ",", ";"))),
		First(
			FlatGroup(
				First(
					typeGroupingParenthesesExpr(),

					FlatGroup(typeExpr(), sp0()),
				),
				ZeroOrMoreTimes(First(
					actualArgumentsExpr(),
					acceessExpr(),
					sliceExpr(),
					FlatGroup(unaryPostfixOp(), sp0()),
				)),
				ZeroOrOnce(First(
					ternaryTailExpr(enableComma),
					binaryTailExpr(enableComma),
				)),
			),
		),
	)
}

//
func typeExpressionImplEnableComma() ParserFn {
	return typeExpressionImpl(true)
}

//
func typeExpressionImplDisableComma() ParserFn {
	return typeExpressionImpl(false)
}

//
func refTypeExpressionImpl(enableComma bool) ParserFn {
	return If(enableComma,
		ref(typeExpressionImplEnableComma),
		ref(typeExpressionImplDisableComma),
	)
}

//
func typeExpressionEnableComma() ParserFn {
	return expression(true)
}

//
func typeExpressionDisableComma() ParserFn {
	return expression(false)
}

//
func refTypeExpression(enableComma bool) ParserFn {
	return If(enableComma,
		ref(typeExpressionEnableComma),
		ref(typeExpressionDisableComma),
	)
}

//
func defExpressionHead() ParserFn {
	return First(
		FlatGroup(
			erase(Seq("let")), sp1(), erase(Seq("mut")),
			Zero(Ast{
				ClassName: clsz.VarOrConst,
				Type:      AstType_String,
				Value:     "let mut",
			}),
		),
		FlatGroup(
			erase(Seq("let")),
			Zero(Ast{
				ClassName: clsz.VarOrConst,
				Type:      AstType_String,
				Value:     "let",
			}),
		),
		FlatGroup(
			erase(Seq("const")),
			Zero(Ast{
				ClassName: clsz.VarOrConst,
				Type:      AstType_String,
				Value:     "const",
			}),
		),
	)
}

//
func defExpressionBody() ParserFn {
	return FlatGroup(
		First(
			FlatGroup(
				erase(Seq("Some")), sp0(), // TODO:
				erase(Seq("(")), sp0(),
				symbol(mnem.Symbol), sp0(),
				First(
					FlatGroup(
						erase(Seq(":")), sp0(),
						typeExpr(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
				erase(Seq(")")), sp0(),
			),
			FlatGroup(
				symbol(mnem.Symbol), sp0(),
				First(
					FlatGroup(
						erase(Seq(":")), sp0(),
						typeExpr(),
					),
					Zero(Ast{ClassName: clsz.Placeholder}),
				),
			),
		),
		First(
			Seq("="),
		),
		sp0(),
		refExpression(false),
	)
}

//
func defExpression() ParserFn {
	return Trans(
		FlatGroup(
			defExpressionHead(),
			sp1(),
			Group(
				defExpressionBody(),
				ZeroOrMoreTimes(
					erase(Seq(",")), sp0(),
					defExpressionBody(),
				),
			),
		),
		TransDefStatement,
	)
}
