package rules

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Logical OR `op1 || op2`
var precedence06 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"||"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok1 := promoteUnaryOperand(ctx, true, false, AstType_Bool, asts[0])
				op2, ok2 := promoteUnaryOperand(ctx, true, false, AstType_Bool, asts[2])

				// NOTE: ptr, any, nil, unit value are not promoted.

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnBool {

					if op1.Value.(bool) {
						return AstSlice{op1}, nil
					} else {
						return AstSlice{op2}, nil
					}

				} else {
					// NOTE: Imm_ptr and Imm_data will be compared at runtime

					if !ok1 {
						op1 = Ast{
							OpCode:    mnem.Convdyn_bool,
							ClassName: "$convdyn_bool",
							Type:      AstType_AstCons,
							Value:     AstCons{Car: op1},
						}
					}
					if !ok2 {
						op2 = Ast{
							OpCode:    mnem.Convdyn_bool,
							ClassName: "$convdyn_bool",
							Type:      AstType_AstCons,
							Value:     AstCons{Car: op2},
						}
					}

					return AstSlice{{
						OpCode:    mnem.LogicalOrbool_bool,
						ClassName: "$logicalor",
						Type:      AstType_AstCons,
						Value: AstCons{
							Car: op1,
							Cdr: Ast{
								OpCode:    mnem.Quote_bool,
								ClassName: clsz.C_op_quote,
								Type:      AstType_AstCons,
								Value:     AstCons{Car: op2},
							},
						},
						SourcePosition: op1.SourcePosition,
					}}, nil
				}
			},
		),
	},
	Rtol: false,
}
