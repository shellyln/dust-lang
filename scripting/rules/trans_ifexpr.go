package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

//
func TransIfExpression(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: multiple conditions and bodies

	args := asts[0].Value.(AstSlice)
	// [i+0]: condition (if / else if) or placeholder (else / nop)
	// [i+1]: body

	origLen := len(args)/2 - 1
	slice := make(AstSlice, origLen, origLen)

	isConstCondExpr := true

	for i := 0; i < len(args); i += 4 {
		switch args[i].ClassName {
		case clsz.IfCondition:
			{
				cond, ok := promoteUnaryOperand(ctx, true, false, AstType_Bool, args[i+2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00004 + args[i+2].Value.(string))
				}

				if isConstCondExpr && cond.OpCode == mnem.Imm_bool {
					if cond.Value.(bool) {
						return AstSlice{args[i+3]}, nil
					}
				} else {
					isConstCondExpr = false
				}

				// variables definition and condition
				if args[i+1].ClassName == clsz.Placeholder {
					slice[i/2] = Ast{
						OpCode:    mnem.Quote,
						ClassName: clsz.C_op_quote,
						Type:      AstType_AstCons,
						Value:     AstCons{Car: cond},
					}
				} else {
					slice[i/2] = Ast{
						OpCode:    mnem.Quote,
						ClassName: clsz.C_op_quote,
						Type:      AstType_AstCons,
						Value: AstCons{
							Car: Ast{
								OpCode:    mnem.Last | (cond.OpCode &^ mnem.OpCodeMask),
								ClassName: clsz.C_op_last,
								Type:      AstType_ListOfAst,
								Value: AstSlice{
									args[i+1],
									cond,
								},
							},
						},
					}
				}

				// body
				slice[i/2+1] = Ast{
					OpCode:    mnem.Quote, // TODO: type
					ClassName: clsz.C_op_quote,
					Type:      AstType_AstCons,
					Value:     AstCons{Car: args[i+3]}, // TODO: type conversion
				}
			}

		case clsz.IfConditionElse:
			{
				if isConstCondExpr {
					return AstSlice{args[i+3]}, nil
				}

				slice[i/2] = Ast{
					OpCode:    mnem.Quote, // TODO: type
					ClassName: clsz.C_op_quote,
					Type:      AstType_AstCons,
					Value:     AstCons{Car: args[i+3]}, // TODO: type conversion
				}
			}

		case clsz.IfConditionNop:
			{
				v := NilAst

				if isConstCondExpr {
					return AstSlice{v}, nil
				}

				slice[i/2] = v
			}

		default:
			return nil, errors.New(emsg.RuleErr00019)
		}
	}

	return AstSlice{Ast{
		OpCode:    mnem.Scope, // TODO: type
		ClassName: clsz.C_op_scope,
		Type:      AstType_AstCons,
		Value: AstCons{Car: Ast{
			OpCode:    mnem.If, // TODO: type
			ClassName: clsz.C_op_if,
			Type:      AstType_ListOfAst,
			Value:     slice,
		}},
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}
