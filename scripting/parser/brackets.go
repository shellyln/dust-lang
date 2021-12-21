package parser

import (
	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func groupingParenthesesExpr() ParserFn {
	return FlatGroup(
		erase(CharClass("(")), sp0(),
		refExpression(true),
		erase(CharClass(")")), sp0(),
	)
}

//
func typeGroupingParenthesesExpr() ParserFn {
	return FlatGroup(
		erase(CharClass("(")), sp0(),
		refTypeExpression(true),
		erase(CharClass(")")), sp0(),
	)
}

//
func acceessExpr() ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.Access,
			Type:      AstType_String,
			Value:     "[]",
		}),
		erase(CharClass("[")), sp0(),
		LookAheadN(CharClass("..")),
		refExpression(true),
		erase(CharClass("]")), sp0(),
	)
}

//
func sliceExpr() ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.Slice,
			Type:      AstType_String,
			Value:     "[]",
		}),
		erase(CharClass("[")), sp0(),
		First(
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}), sp0(),
				erase(CharClass(":", "..")), sp0(),
				refExpression(true),
			),
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}), sp0(),
				erase(CharClass(":", "..")), sp0(),
				Zero(Ast{ClassName: clsz.Placeholder}), sp0(),
			),
			FlatGroup(
				refExpression(true),
				// NOTE: Go style and Rust style
				erase(CharClass(":", "..")), sp0(),
				refExpression(true),
			),
			FlatGroup(
				refExpression(true),
				erase(CharClass(":", "..")), sp0(),
				Zero(Ast{ClassName: clsz.Placeholder}),
			),
		),
		erase(CharClass("]")), sp0(),
	)
}

//
func actualArgumentsExpr() ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.Call,
			Type:      AstType_String,
			Value:     "()",
		}),
		erase(CharClass("(")), sp0(),
		Group(First(
			FlatGroup(
				erase(CharClass(")")), sp0(),
			),
			FlatGroup(
				refExpression(false),
				ZeroOrMoreTimes(
					erase(CharClass(",")), sp0(),
					refExpression(false),
				),
				ZeroOrOnce(
					erase(CharClass(",")), sp0(),
				),
				erase(CharClass(")")), sp0(),
			),
			Error(emsg.ParserErr00003),
		)),
	)
}

//
func listLiteral() ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.List,
			Type:      AstType_String,
			Value:     "[]",
		}),
		First(
			FlatGroup(
				WordBoundary(),
				// NOTE: Rust standard or other macros
				CharClass(
					"vec!",
					// NOTE: TODO: typed array
					"i8!", "u8!", "i16!", "u16!", "i32!", "u32!", "i64!", "u64!",
					"f32!", "f64!",
					"bool!", "str!",
				), sp0(),
				erase(CharClass("[")),
			),
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}), // macro name
				erase(CharClass("[")),
			),
		),
		sp0(),
		Group(First(
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}), // size
				erase(CharClass("]")), sp0(),
			),
			Trans(
				FlatGroup(
					refExpression(false),
					erase(CharClass(";")), sp0(),
					refExpression(false),
					erase(CharClass("]")), sp0(),
				),
				func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
					return AstSlice{asts[1], asts[0]}, nil
				},
			),
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}), // size
				// TODO: spread
				refExpression(false),
				ZeroOrMoreTimes(
					erase(CharClass(",")), sp0(),
					// TODO: spread
					refExpression(false),
				),
				ZeroOrOnce(
					erase(CharClass(",")), sp0(),
				),
				erase(CharClass("]")), sp0(),
			),
			Error(emsg.ParserErr00004),
		)),
	)
}

//
func objectKey() ParserFn {
	return Trans(
		First(
			stringLiteral(),
			symbol(mnem.Imm_str),
		),
	)
}

//
func objectKeyValuePair() ParserFn {
	return FlatGroup(
		First(
			FlatGroup(
				objectKey(), sp0(),
			),
			FlatGroup(
				erase(CharClass("[")), sp0(),
				refExpression(false),
				erase(CharClass("]")), sp0(),
			),
		),
		// NOTE: JS/Go style and Rust style
		erase(CharClass(":", "=>")), sp0(),
		First(
			refExpression(false),
			Error(emsg.ParserErr00008),
		),
	)
}

//
func objectLiteral() ParserFn {
	return FlatGroup(
		Zero(Ast{
			ClassName: clsz.Object,
			Type:      AstType_String,
			Value:     "{}",
		}),
		First(
			FlatGroup(
				WordBoundary(),
				// TODO: any user defined struct names
				// NOTE: Rust maplit crate or some other macros
				CharClass("hashmap!", "map!", "collection!"), sp0(),
				erase(CharClass("{")),
			),
			FlatGroup(
				Zero(Ast{ClassName: clsz.Placeholder}),
				erase(CharClass("{")),
			),
		),
		sp0(),
		Group(First(
			FlatGroup(
				erase(CharClass("}")), sp0(),
			),
			FlatGroup(
				// TODO: spread
				objectKeyValuePair(),
				ZeroOrMoreTimes(
					erase(CharClass(",")), sp0(),
					// TODO: spread
					objectKeyValuePair(),
				),
				ZeroOrOnce(
					erase(CharClass(",")), sp0(),
				),
				First(
					FlatGroup(erase(CharClass("}")), sp0()),
					Error(emsg.ParserErr00008),
				),
			),
		)),
	)
}
