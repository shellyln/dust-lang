package parser

import (
	"math"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/extra"
	. "github.com/shellyln/takenoco/string"
)

//
func integerSizeSuffix() ParserFn {
	return First(
		Seq("usize"),
		Seq("u64"),
		Seq("u32"),
		Seq("u16"),
		Seq("u8"),
		Seq("isize"),
		Seq("i64"),
		Seq("i32"),
		Seq("i16"),
		Seq("i8"),
	)
}

//
func floatSizeSuffix() ParserFn {
	return First(
		Seq("f64"),
		Seq("f32"),
	)
}

//
func numberSizeDetector(signed TransformerFn, unsigned TransformerFn) ParserFn {
	return Match(2)(
		Case{If: TopIsStr("isize"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_i64|AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_i64,
				),
			},
		},
		Case{If: TopIsStr("usize"),
			Let: []TransformerFn{Pop, unsigned,
				SetOpCodeAndClassName(
					mnem.Imm_u64|AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_u64,
				),
			},
		},
		Case{If: TopIsStr("i64"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_i64|AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_i64,
				),
			},
		},
		Case{If: TopIsStr("u64"),
			Let: []TransformerFn{Pop, unsigned,
				SetOpCodeAndClassName(
					mnem.Imm_u64|AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_u64,
				),
			},
		},
		Case{If: TopIsStr("i32"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_i32|AstOpCodeType(tsys.TypeInfo_Primitive_I32.Id<<mnem.TypeIdOffset),
					clsz.C_imm_i32,
				),
			},
		},
		Case{If: TopIsStr("u32"),
			Let: []TransformerFn{Pop, unsigned,
				SetOpCodeAndClassName(
					mnem.Imm_u32|AstOpCodeType(tsys.TypeInfo_Primitive_U32.Id<<mnem.TypeIdOffset),
					clsz.C_imm_u32,
				),
			},
		},
		Case{If: TopIsStr("i16"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_i16|AstOpCodeType(tsys.TypeInfo_Primitive_I16.Id<<mnem.TypeIdOffset),
					clsz.C_imm_i16,
				),
			},
		},
		Case{If: TopIsStr("u16"),
			Let: []TransformerFn{Pop, unsigned,
				SetOpCodeAndClassName(
					mnem.Imm_u16|AstOpCodeType(tsys.TypeInfo_Primitive_U16.Id<<mnem.TypeIdOffset),
					clsz.C_imm_u16,
				),
			},
		},
		Case{If: TopIsStr("i8"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_i8|AstOpCodeType(tsys.TypeInfo_Primitive_I8.Id<<mnem.TypeIdOffset),
					clsz.C_imm_i8,
				),
			},
		},
		Case{If: TopIsStr("u8"),
			Let: []TransformerFn{Pop, unsigned,
				SetOpCodeAndClassName(
					mnem.Imm_u8|AstOpCodeType(tsys.TypeInfo_Primitive_U8.Id<<mnem.TypeIdOffset),
					clsz.C_imm_u8,
				),
			},
		},
		Case{If: TopIsStr("f64"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_f64|AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_f64,
				),
			},
		},
		Case{If: TopIsStr("f32"),
			Let: []TransformerFn{Pop, signed,
				SetOpCodeAndClassName(
					mnem.Imm_f32|AstOpCodeType(tsys.TypeInfo_Primitive_F32.Id<<mnem.TypeIdOffset),
					clsz.C_imm_f32,
				),
			},
		},
	)
}

//
func numberConstSizeDetector() ParserFn {
	return Match(2)(
		Case{If: TopIsStr("f64"),
			Let: []TransformerFn{Pop,
				SetOpCodeAndClassName(
					mnem.Imm_f64|AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id<<mnem.TypeIdOffset),
					clsz.C_imm_f64,
				),
			},
		},
		Case{If: TopIsStr("f32"),
			Let: []TransformerFn{Pop,
				SetOpCodeAndClassName(
					mnem.Imm_f32|AstOpCodeType(tsys.TypeInfo_Primitive_F32.Id<<mnem.TypeIdOffset),
					clsz.C_imm_f32,
				),
			},
		},
	)
}

//
func binaryNumberLiteral() ParserFn {
	return FlatGroup(
		FlatGroup(
			erase(SeqI("0b")),
			BinaryNumberStr(),
		),
		First(
			integerSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "i64"}),
		),
		WordBoundary(),
		numberSizeDetector(ParseIntRadix(2), ParseUintRadix(2)),
	)
}

//
func octalNumberLiteral() ParserFn {
	return FlatGroup(
		FlatGroup(
			erase(Seq("0o")),
			OctalNumberStr(),
		),
		First(
			integerSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "i64"}),
		),
		WordBoundary(),
		numberSizeDetector(ParseIntRadix(8), ParseUintRadix(8)),
	)
}

//
func hexNumberLiteral() ParserFn {
	return FlatGroup(
		FlatGroup(
			erase(SeqI("0x")),
			HexNumberStr(),
		),
		First(
			integerSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "i64"}),
		),
		WordBoundary(),
		numberSizeDetector(ParseIntRadix(16), ParseUintRadix(16)),
	)
}

//
func integerNumberLiteral() ParserFn {
	return FlatGroup(
		IntegerNumberStr(),
		First(
			First(
				integerSizeSuffix(),
				floatSizeSuffix(),
			),
			Zero(Ast{Type: AstType_String, Value: "i64"}),
		),
		WordBoundary(),
		numberSizeDetector(ParseInt, ParseUint),
	)
}

//
func floatNumberLiteral() ParserFn {
	return FlatGroup(
		FloatNumberStr(),
		First(
			floatSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "f64"}),
		),
		WordBoundary(),
		numberSizeDetector(ParseFloat, ParseFloat),
	)
}

//
func negativeInfinityLiteral() ParserFn {
	return FlatGroup(
		erase(Seq("-Infinity")),
		Zero(Ast{
			Type:  AstType_Float,
			Value: math.Inf(-1),
		}),
		erase(ZeroOrMoreTimes(Seq("_"))),
		First(
			floatSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "f64"}),
		),
		WordBoundary(),
		numberConstSizeDetector(),
	)
}

//
func positiveInfinityLiteral() ParserFn {
	return FlatGroup(
		erase(First(
			Seq("+Infinity"),
			Seq("Infinity"),
		)),
		Zero(Ast{
			Type:  AstType_Float,
			Value: math.Inf(0),
		}),
		erase(ZeroOrMoreTimes(Seq("_"))),
		First(
			floatSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "f64"}),
		),
		WordBoundary(),
		numberConstSizeDetector(),
	)
}

//
func nanLiteral() ParserFn {
	return FlatGroup(
		erase(Seq("NaN")),
		Zero(Ast{
			Type:  AstType_Float,
			Value: math.NaN(),
		}),
		erase(ZeroOrMoreTimes(Seq("_"))),
		First(
			floatSizeSuffix(),
			Zero(Ast{Type: AstType_String, Value: "f64"}),
		),
		WordBoundary(),
		numberConstSizeDetector(),
	)
}

//
func numberLiteral() ParserFn {
	return First(
		negativeInfinityLiteral(),
		positiveInfinityLiteral(),
		nanLiteral(),
		binaryNumberLiteral(),
		octalNumberLiteral(),
		hexNumberLiteral(),
		floatNumberLiteral(),
		integerNumberLiteral(),
	)
}

//
func trueLiteral() ParserFn {
	return FlatGroup(
		erase(First(
			Seq("true"),
		)),
		WordBoundary(),
		Zero(TrueAst),
	)
}

//
func falseLiteral() ParserFn {
	return FlatGroup(
		erase(First(
			Seq("false"),
		)),
		WordBoundary(),
		Zero(FalseAst),
	)
}

//
func boolLiteral() ParserFn {
	return First(trueLiteral(), falseLiteral())
}

//
func nullLiteral() ParserFn {
	return FlatGroup(
		erase(First(
			CharClass("null", "None"), // NOTE: None is not global symbol in Rust
		)),
		WordBoundary(),
		Zero(NilAst),
	)
}

//
func undefinedLiteral() ParserFn {
	return FlatGroup(
		erase(Seq("undefined")),
		WordBoundary(),
		Zero(UnitAst),
	)
}

//
func unitValueLiteral() ParserFn {
	return FlatGroup(
		erase(Seq("()")),
		Zero(UnitAst),
	)
}

// TODO:
func stringLiteralInner(cc string, multiline bool) ParserFn {
	return FlatGroup(
		erase(Seq(cc)),
		ZeroOrMoreTimes(
			First(
				FlatGroup(
					erase(Seq("\\")),
					First(
						CharClass("\\", "'", "\"", "`"),
						replaceStr(CharClass("n"), "\n"),
						replaceStr(CharClass("r"), "\r"),
						replaceStr(CharClass("v"), "\v"),
						replaceStr(CharClass("t"), "\t"),
						replaceStr(CharClass("b"), "\b"),
						replaceStr(CharClass("f"), "\f"),
						Trans(
							FlatGroup(
								erase(CharClass("u")),
								Repeat(Times{Min: 4, Max: 4}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								erase(CharClass("u{")),
								Repeat(Times{Min: 1, Max: 6}, HexNumber()),
								erase(CharClass("}")),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								erase(CharClass("x")),
								Repeat(Times{Min: 2, Max: 2}, HexNumber()),
							),
							ParseIntRadix(16),
							StringFromInt,
						),
						Trans(
							FlatGroup(
								Repeat(Times{Min: 3, Max: 3}, OctNumber()),
							),
							ParseIntRadix(8),
							StringFromInt,
						),
					),
				),
				If(multiline,
					OneOrMoreTimes(CharClassN(cc, "\\")),
					OneOrMoreTimes(
						First(
							FlatGroup(
								CharClass("\r", "\n"),
								Error(emsg.ParserErr00007),
							),
							CharClassN(cc, "\\"),
						),
					),
				),
			),
		),
		First(
			FlatGroup(End(), Error(emsg.ParserErr00006)),
			erase(Seq(cc)),
		),
	)
}

//
func singleQuoteStringliteral() ParserFn {
	return Trans(
		stringLiteralInner("'", true),
		Concat,
		SetOpCodeAndClassName(
			mnem.Imm_str|AstOpCodeType(tsys.TypeInfo_Primitive_String.Id<<mnem.TypeIdOffset),
			clsz.C_imm_str,
		),
	)
}

//
func doubleQuoteStringliteral() ParserFn {
	return Trans(
		stringLiteralInner("\"", true),
		Concat,
		SetOpCodeAndClassName(
			mnem.Imm_str|AstOpCodeType(tsys.TypeInfo_Primitive_String.Id<<mnem.TypeIdOffset),
			clsz.C_imm_str,
		),
	)
}

//
func backQuoteStringliteral() ParserFn {
	return Trans(
		stringLiteralInner("`", true), // TODO: heredocs
		Concat,
		SetOpCodeAndClassName(
			mnem.Imm_str|AstOpCodeType(tsys.TypeInfo_Primitive_String.Id<<mnem.TypeIdOffset),
			clsz.C_imm_str,
		),
	)
}

//
func stringLiteral() ParserFn {
	return First(
		singleQuoteStringliteral(),
		doubleQuoteStringliteral(),
		backQuoteStringliteral(),
	)
}

//
func literal() ParserFn {
	return First(
		numberLiteral(),
		boolLiteral(),
		nullLiteral(),
		undefinedLiteral(),
		unitValueLiteral(),
		stringLiteral(),
		listLiteral(),
		objectLiteral(),
	)
}

//
func symbol(opcode AstOpCodeType) ParserFn {
	return Trans(
		FlatGroup(
			erase(ZeroOrOnce(Seq("r#"))),
			Once(First(
				Alpha(),
				CharClass("_", "$"),
			)),
			ZeroOrMoreTimes(First(
				Alnum(),
				CharClass("_", "$"),
			)),
		),
		Concat,
		SetOpCodeAndClassName(opcode, clsz.C_symbol),
	)
}

//
func literalOrSymbol() ParserFn {
	return First(literal(), symbol(mnem.Symbol))
}
