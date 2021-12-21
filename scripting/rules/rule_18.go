package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Postfix Increment `op1 ++`
// Postfix Decrement `op1 --`
var precedence18 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.UnaryPostfixOp, []string{"++", "--"}),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperand(ctx, false, true, AstType_Float, asts[0])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}

				if op1.OpCode&mnem.Lvalue == 0 {
					return nil, errors.New(emsg.RuleErr00002 + asts[1].Value.(string))
				}

				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}

				var opcode AstOpCodeType
				var err error

				if asts[1].Value.(string) == "++" {
					opcode = mnem.PostIncr | mnem.Lvalue
				} else {
					opcode = mnem.PostDecr | mnem.Lvalue
				}

				ret, err = generateUnaryOpResult(ctx, opcode, op1)
				if err != nil {
					return nil, err
				}

				if ret.OpCode&mnem.ReturnTypeMask != mnem.ReturnAny {
					return AstSlice{ret}, nil
				} else {
					return nil, errors.New(emsg.RuleErr00006 + asts[1].Value.(string))
				}
			},
		),
	},
	Rtol: false,
}
