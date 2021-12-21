package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Pipeline Function Call `op1 |> op2(args)`
var precedence04a = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"|>"}),
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
				op3 := asts[4]
				args := op3.Value.(AstSlice)

				opcode := mnem.DynCall
				className := "$dyncall"

				// TODO: Func opcode
				if op2.OpCode&mnem.OpCodeMask == mnem.Func {
					//
				}

				// TODO: Func, Symbol->Func opcode
				// NOTE: For cases of MapIndex and Index, call type and argument types are not determined.
				retType, _, err := getReturnTypeAndFlags(xctx, true, false, false, op2)
				if err != nil {
					return nil, err
				}
				// TODO: Do not check for the presence of the symbol until back reference have been implemented.

				// TODO: check op1 is callable

				opcode |= retType

				slice := make(AstSlice, 0, len(args)+2)
				slice = append(slice, op2)
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
	},
	Rtol: false,
}
