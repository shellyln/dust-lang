package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Conditional (ternary) operator `op1 ? op2 : op3`
var precedence04b = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.TernaryOp, []string{"?:"}),
				anyOperand(),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				op1, ok := promoteUnaryOperand(ctx, true, false, AstType_Bool, asts[0])
				if !ok {
					return nil, errors.New(emsg.RuleErr00004 + asts[1].Value.(string))
				}

				if op1.OpCode == mnem.Imm_bool {
					if op1.Value.(bool) {
						return asts[2:3], nil // TODO: type conversion
					} else {
						return asts[3:4], nil // TODO: type conversion
					}
				} else {
					op2Type, _, err := getReturnTypeAndFlags(xctx, false, false, false, asts[2])
					if err != nil {
						return nil, err
					}
					op3Type, _, err := getReturnTypeAndFlags(xctx, false, false, false, asts[3])
					if err != nil {
						return nil, err
					}

					if op2Type != op3Type { // TODO: BUG: compare type id
						// TODO: Cast is needed
					}

					return AstSlice{{
						OpCode:    mnem.If, // TODO: type
						ClassName: clsz.C_op_if,
						Type:      AstType_ListOfAst,
						Value: AstSlice{
							Ast{
								OpCode:         mnem.Quote,
								ClassName:      clsz.C_op_quote,
								Type:           AstType_AstCons,
								Value:          AstCons{Car: op1},
								SourcePosition: op1.SourcePosition,
							},
							Ast{
								OpCode:         mnem.Quote, // TODO: type
								ClassName:      clsz.C_op_quote,
								Type:           AstType_AstCons,
								Value:          AstCons{Car: asts[2]}, // TODO: type conversion
								SourcePosition: op1.SourcePosition,
							},
							Ast{
								OpCode:         mnem.Quote, // TODO: type
								ClassName:      clsz.C_op_quote,
								Type:           AstType_AstCons,
								Value:          AstCons{Car: asts[3]}, // TODO: type conversion
								SourcePosition: op1.SourcePosition,
							},
						},
						SourcePosition: op1.SourcePosition,
					}}, nil
				}
			},
		),
	},
	Rtol: true,
}
