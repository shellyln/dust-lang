package rules

import (
	"errors"
	"strings"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

// Member Access `op1 . op2`
// Computed Member Access `op1[op2]`
// Slicing `op1[start..end]`
// Type cast `op1 as T` / `op1.(T)`
// Function Call `op1(args)`
var precedence20 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				LookBehindN(2, 2,
					anyOperand(),
					isOperator(clsz.BinaryOp, []string{"|>"}),
				),
				anyOperand(),
				isOperator(clsz.Call, []string{"()"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				op1 := asts[0]
				op2 := asts[2]
				args := op2.Value.(AstSlice)

				opcode := mnem.DynCall
				className := "$dyncall"

				// TODO: Func opcode
				if op1.OpCode&mnem.OpCodeMask == mnem.Func || op1.OpCode&mnem.OpCodeMask == mnem.Lambda {
					//
				}

				if op1.OpCode&mnem.OpCodeMask == mnem.Symbol {
					name := op1.Value.(string)
					if name == "Some" {
						if len(args) == 1 {
							return AstSlice{args[0]}, nil
						} else {
							return nil, errors.New(emsg.RuleErr00009)
						}
					}

					// TODO: BUG: resolve symbol and get address.
				}

				// NOTE: For cases of MapIndex and Index, call type and argument types are not determined.
				// TODO: Do not check for the presence of the symbol until back reference have been implemented.
				// TODO: Do not check `IsCallable` is symbol is any
				retType, _, _ := getReturnTypeAndFlags(xctx, true, false, false, op1)
				// TODO: BUG: ignore error

				// TODO: check op1 is callable

				opcode |= retType

				slice := make(AstSlice, 0, len(args)+1)
				slice = append(slice, op1)
				slice = append(slice, args...)

				return AstSlice{{
					OpCode:         opcode,
					ClassName:      className,
					Type:           AstType_ListOfAst,
					Value:          slice,
					SourcePosition: op1.SourcePosition,
				}}, nil
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.List, []string{"[]"}),
				anyOperand(),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				// [0]: operator
				// [1]: macro name or placeholder
				// [2]: [0]: size or placeholder, [1:]: body

				body := asts[2].Value.(AstSlice)
				args := body[1:]

				var itemTypeInfo tsys.TypeInfo // `any` type
				if asts[1].ClassName != clsz.Placeholder {
					s := asts[1].Value.(string)
					if strings.HasSuffix(s, "!") {
						s = s[:len(s)-1]
						itemTypeInfo, _ = xctx.GetTypeInfoByName(s)
					}
				}

				itemType := xtor.ToAstType(itemTypeInfo.Flags)
				sliceTypeInfo, err := xctx.GetRegisteredTypeInfo(itemTypeInfo.Slice(xctx.DummyTypeId))
				if err != nil {
					return nil, err
				}

				if body[0].ClassName == clsz.Placeholder {
					slice := make(AstSlice, 0, len(args))
					slice = append(slice, args...)

					if itemTypeInfo.Flags != mnem.ReturnAny {
						for i := 0; i < len(slice); i++ {
							op, ok := promoteUnaryOperand(ctx, true, false, itemType, slice[i])
							if !ok {
								return nil, errors.New(emsg.RuleErr00021 + asts[1].Value.(string))
							}
							slice[i] = op
						}
					}

					return AstSlice{{
						OpCode:         mnem.List | mnem.Indexable | sliceTypeInfo.Flags,
						ClassName:      "$list",
						Type:           AstType_ListOfAst,
						Value:          slice,
						SourcePosition: asts[0].SourcePosition,
					}}, nil
				} else {
					opCar, ok := promoteUnaryOperand(ctx, true, true, AstType_Int, body[0])
					if !ok {
						return nil, errors.New(emsg.RuleErr00013)
					}

					opCdr := body[1]

					if itemTypeInfo.Flags != mnem.ReturnAny {
						opCdr, ok = promoteUnaryOperand(ctx, true, false, itemType, opCdr)
						if !ok {
							return nil, errors.New(emsg.RuleErr00021 + asts[1].Value.(string))
						}
					}

					return AstSlice{{
						OpCode:    mnem.FilledList | mnem.Indexable | sliceTypeInfo.Flags,
						ClassName: "$filledlist",
						Type:      AstType_AstCons,
						Value: AstCons{
							Car: opCar,
							Cdr: opCdr,
						},
						SourcePosition: asts[0].SourcePosition,
					}}, nil
				}
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.Object, []string{"{}"}),
				anyOperand(),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				// [0]: operator
				// [1]: macro name or placeholder
				// [2]: body

				args := asts[2].Value.(AstSlice)
				slice := make(AstSlice, 0, len(args))
				slice = append(slice, args...)

				for i := 0; i < len(slice); i += 2 {
					if slice[i].OpCode != mnem.Imm_str {
						key, ok := promoteUnaryOperand(ctx, true, true, AstType_String, slice[i])
						if !ok {
							return nil, errors.New(emsg.RuleErr00012)
						}
						slice[i] = key
					}
				}

				return AstSlice{{
					OpCode:         mnem.Object | mnem.Indexable,
					ClassName:      "$object",
					Type:           AstType_ListOfAst,
					Value:          slice,
					SourcePosition: asts[0].SourcePosition,
				}}, nil
			},
		),
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.Access, []string{"[]", "."}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				var opcode mnem.AstOpCodeType
				var className string
				op1 := asts[0]

				op2, ok := promoteUnaryOperand(ctx, false, true, AstType_String, asts[2])

				if asts[1].Value.(string) == "." && op2.OpCode&mnem.OpCodeMask == mnem.Symbol {
					opcode = mnem.Mapindex
					className = "$mapindex"
					op2.OpCode = mnem.Imm_str
					op2.ClassName = clsz.C_imm_str

				} else if ok && asts[1].Value.(string) == "[]" && op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnString {
					opcode = mnem.Mapindex
					className = "$mapindex"

				} else {
					op2, ok = promoteUnaryOperand(ctx, true, true, AstType_Int, asts[2])
					if !ok {
						return nil, errors.New(emsg.RuleErr00013)
					}

					opcode = mnem.Index
					className = "$index"

					// NOTE: The result of the index operation of the string is u8.
					//       The result of the index operation of the T[] is T.
					retType, _, err := getReturnTypeAndFlags(xctx, false, false, true, op1)
					if err != nil {
						return nil, err
					}

					opcode |= retType
				}

				return AstSlice{{
					OpCode:         opcode | mnem.Lvalue,
					ClassName:      className,
					Type:           AstType_AstCons,
					Value:          AstCons{Car: op1, Cdr: op2},
					SourcePosition: op1.SourcePosition,
				}}, nil
			},
		),
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.Slice, []string{"[]"}),
				anyOperand(),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				var op1, op2, op3 Ast
				op1 = asts[0]

				if asts[2].ClassName == clsz.Placeholder {
					op2 = NilAst
				} else {
					op2, ok = promoteUnaryOperand(ctx, true, true, AstType_Int, asts[2])
					if !ok {
						return nil, errors.New(emsg.RuleErr00014)
					}
				}

				if asts[3].ClassName == clsz.Placeholder {
					op3 = NilAst
				} else {
					op3, ok = promoteUnaryOperand(ctx, true, true, AstType_Int, asts[3])
					if !ok {
						return nil, errors.New(emsg.RuleErr00014)
					}
				}

				retType, _, err := getReturnTypeAndFlags(xctx, false, false, false, op1)
				if err != nil {
					return nil, err
				}

				return AstSlice{{
					OpCode:         mnem.Slice | retType | mnem.Lvalue,
					ClassName:      "$slice",
					Type:           AstType_ListOfAst,
					Value:          AstSlice{op1, op2, op3},
					SourcePosition: op1.SourcePosition,
				}}, nil
			},
		),
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.TypeCast, []string{"as"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				var opcode mnem.AstOpCodeType
				var className string
				op1 := asts[0]
				op2 := asts[2] // TypeInfo

				retType, _, err := getReturnTypeAndFlags(xctx, false, false, false, op1)
				if err != nil {
					return nil, err
				}

				var opcodeTmp AstOpCodeType
				opcodeTmp, className = getConversionOpCode(retType)
				opcode |= opcodeTmp

				typeInfo := op2.Value.(tsys.TypeInfo)
				opcodeTmp = typeInfo.Flags

				if retType&mnem.ReturnTypeMask == opcodeTmp&mnem.ReturnTypeMask ||
					mnem.ReturnAny == opcodeTmp&mnem.ReturnTypeMask {

					return AstSlice{op1}, nil
				} else {
					opcode |= opcodeTmp

					return AstSlice{{
						OpCode:    opcode,
						ClassName: className,
						Type:      AstType_AstCons,
						Value:     AstCons{Car: op1},
					}}, nil
				}
			},
		),
	},
	Rtol: false,
}
