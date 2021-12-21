package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Addition `op1 + op2`
// Subtraction `op1 - op2`
var precedence14 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"+"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
				if !ok {
					return nil, errors.New(emsg.InternalErr00001)
				}

				opcode := mnem.Add

				ops, ok := promoteBinaryOperands(ctx, false, true, AstType_Float, asts[0], asts[2])
				if !ok {
					op1Type, _, _ := getReturnTypeAndFlags(xctx, false, false, false, asts[0])
					op2Type, _, _ := getReturnTypeAndFlags(xctx, false, false, false, asts[2])

					if op1Type&mnem.ReturnTypeMask == mnem.ReturnString || op2Type&mnem.ReturnTypeMask == mnem.ReturnString {
						opcode = mnem.Concat_str

						ops, ok = promoteBinaryOperands(ctx, true, true, AstType_String, asts[0], asts[2])
						if !ok {
							return nil, errors.New(emsg.RuleErr00010 + asts[1].Value.(string))
						}
					} else {
						return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
					}
				}

				op1 := ops[0]
				op2 := ops[1]
				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Int
					ret.Value = op1.Value.(int64) + op2.Value.(int64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Uint
					ret.Value = op1.Value.(uint64) + op2.Value.(uint64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Float
					ret.Value = op1.Value.(float64) + op2.Value.(float64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnString &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnString {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_String
					ret.Value = op1.Value.(string) + op2.Value.(string)

				} else {
					var err error
					ret, err = generateBinaryOpResult(ctx, opcode, op1, op2)
					if err != nil {
						return nil, err
					}
				}

				if ret.OpCode&mnem.ReturnTypeMask != mnem.ReturnAny {
					return AstSlice{ret}, nil
				} else {
					return nil, errors.New(emsg.RuleErr00006 + asts[1].Value.(string))
				}
			},
		),
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"-"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				ops, ok := promoteBinaryOperands(ctx, false, true, AstType_Float, asts[0], asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}
				op1 := ops[0]
				op2 := ops[1]
				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}

				opcode := mnem.Sub

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Int
					ret.Value = op1.Value.(int64) - op2.Value.(int64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Uint
					ret.Value = op1.Value.(uint64) - op2.Value.(uint64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Float
					ret.Value = op1.Value.(float64) - op2.Value.(float64)

				} else {
					var err error
					ret, err = generateBinaryOpResult(ctx, opcode, op1, op2)
					if err != nil {
						return nil, err
					}
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
