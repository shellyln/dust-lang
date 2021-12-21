package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Assignment `op1 = op2`
var precedence03 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"="}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				op1 := asts[0]
				op1Type, _, err := getReturnTypeAndFlags(xctx, false, false, false, op1)
				if err != nil {
					return nil, errors.New(emsg.RuleErr00002 + asts[1].Value.(string))
				}
				if op1Type&mnem.Lvalue == 0 {
					return nil, errors.New(emsg.RuleErr00002 + asts[1].Value.(string))
				}

				// TODO: BUG: Complex type equality checking
				op2, ok := promoteUnaryOperand(ctx, true, false, xtor.ToAstType(op1Type), asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00003)
				}

				return AstSlice{{
					OpCode:         mnem.Assign | op1Type, // op1 is already lvalue
					ClassName:      "$assign",
					Type:           AstType_AstCons,
					Value:          AstCons{Car: op1, Cdr: op2},
					SourcePosition: asts[0].SourcePosition,
				}}, nil
			},
		),
	},
	Rtol: true,
}
