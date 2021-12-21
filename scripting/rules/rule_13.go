package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Bitwise Left Shift `op1 << op2`
// Bitwise Right Shift `op1 >> op2`
// Bitwise Unsigned Right Shift `op1 >>> op2`
var precedence13 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"<<", ">>"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperandForBitwise(ctx, asts[0])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}

				op2, ok := promoteUnaryOperand(ctx, true, true, AstType_Uint, asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}

				var opcode AstOpCodeType
				switch asts[1].Value.(string) {
				case "<<":
					opcode = mnem.BitwiseLShift
				case ">>":
					opcode = mnem.BitwiseSRShift
				default:
					return nil, errors.New(emsg.RuleErr00007 + asts[1].Value.(string))
				}

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					var v int64
					switch opcode {
					case mnem.BitwiseLShift:
						v = op1.Value.(int64) << op2.Value.(uint64)
					case mnem.BitwiseSRShift:
						v = op1.Value.(int64) >> op2.Value.(uint64)
					}
					return AstSlice{{
						OpCode:         mnem.Imm_i64,
						ClassName:      clsz.C_imm_i64,
						Type:           AstType_Int,
						Value:          v,
						SourcePosition: op1.SourcePosition,
					}}, nil

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					var v uint64
					switch opcode {
					case mnem.BitwiseLShift:
						v = op1.Value.(uint64) << op2.Value.(uint64)
					case mnem.BitwiseSRShift:
						v = op1.Value.(uint64) >> op2.Value.(uint64)
					}
					return AstSlice{{
						OpCode:         mnem.Imm_u64,
						ClassName:      clsz.C_imm_u64,
						Type:           AstType_Uint,
						Value:          v,
						SourcePosition: op1.SourcePosition,
					}}, nil
				}

				ret, err := generateUnaryOpResult(ctx, opcode, op1)
				if err != nil {
					return nil, err
				}
				cons := ret.Value.(AstCons)
				cons.Cdr = op2
				ret.Value = cons

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
				isOperator(clsz.BinaryOp, []string{">>>"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				opx, ok := promoteUnaryOperandForBitwise(ctx, asts[0])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}
				optyp := opx.OpCode & mnem.ReturnTypeMask

				op1, ok := promoteUnaryOperand(ctx, true, true, AstType_Uint, asts[0])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}

				op2, ok := promoteUnaryOperand(ctx, true, true, AstType_Uint, asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
				}

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					v := op1.Value.(uint64) >> op2.Value.(uint64)
					if optyp == mnem.ReturnInt {
						return AstSlice{{
							OpCode:         mnem.Imm_i64,
							ClassName:      clsz.C_imm_i64,
							Type:           AstType_Int,
							Value:          int64(v),
							SourcePosition: op1.SourcePosition,
						}}, nil
					} else {
						return AstSlice{{
							OpCode:         mnem.Imm_u64,
							ClassName:      clsz.C_imm_u64,
							Type:           AstType_Uint,
							Value:          v,
							SourcePosition: op1.SourcePosition,
						}}, nil
					}
				}

				if op1.OpCode&mnem.ReturnTypeMask != optyp {
					op1, ok = promoteUnaryOperand(ctx, true, false, xtor.ToAstType(optyp), asts[0])
					if !ok {
						return nil, errors.New(emsg.RuleErr00005 + asts[1].Value.(string))
					}
				}

				ret, err := generateUnaryOpResult(ctx, mnem.BitwiseURShift, op1)
				if err != nil {
					return nil, err
				}
				cons := ret.Value.(AstCons)
				cons.Cdr = op2
				ret.Value = cons

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
