package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Comma / Sequence `op1, op2`
var precedence01 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{","}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				if mnem.Imm_begin < asts[0].OpCode && asts[0].OpCode < mnem.Imm_end {
					return asts[2:3], nil
				} else {
					retType, _, err := getReturnTypeAndFlags(xctx, false, false, false, asts[2])
					if err != nil {
						return nil, err
					}

					return AstSlice{{
						OpCode:         mnem.Seq | retType,
						ClassName:      clsz.C_op_seq,
						Type:           AstType_AstCons,
						Value:          AstCons{Car: asts[0], Cdr: asts[2]},
						SourcePosition: asts[0].SourcePosition,
					}}, nil
				}
			},
		),
	},
	Rtol: false,
}
