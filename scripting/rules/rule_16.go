package rules

import (
	"errors"
	"math"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Exponentiation `op1 ** op2`
var precedence16 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"**"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				ops, ok := promoteBinaryOperands(ctx, true, true, AstType_Float, asts[0], asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}
				op1 := ops[0]
				op2 := ops[1]

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat {

					return AstSlice{{
						OpCode:         mnem.Imm_f64,
						ClassName:      clsz.C_imm_f64,
						Type:           AstType_Float,
						Value:          math.Pow(op1.Value.(float64), op2.Value.(float64)),
						SourcePosition: op1.SourcePosition,
					}}, nil

				} else {
					return AstSlice{{
						OpCode:         mnem.Pow_f,
						ClassName:      "$pow_f",
						Type:           AstType_AstCons,
						Value:          AstCons{Car: op1, Cdr: op2},
						SourcePosition: op1.SourcePosition,
					}}, nil
				}
			},
		),
	},
	Rtol: true,
}
