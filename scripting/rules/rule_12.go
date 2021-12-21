package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Less Than `op1 < op2`
// Less Than Or Equal `op1 <= op2`
// Greater Than `op1 > op2`
// Greater Than Or Equal `op1 >= op2`
var precedence12 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"<", "<=", ">", ">="}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				ops, ok := promoteBinaryOperands(ctx, false, true, AstType_String, asts[0], asts[2])
				if !ok {
					return nil, errors.New(emsg.RuleErr00008 + asts[1].Value.(string))
				}
				op1 := ops[0]
				op2 := ops[1]
				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}
				var opcode AstOpCodeType

				if op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnBool {
					return nil, errors.New(emsg.RuleErr00008 + asts[1].Value.(string))
				}

				switch asts[1].Value.(string) {
				case "<":
					switch op1.OpCode & mnem.ReturnTypeMask {
					case mnem.ReturnInt:
						opcode = mnem.CmpLTi_bool
					case mnem.ReturnUint:
						opcode = mnem.CmpLTu_bool
					case mnem.ReturnFloat:
						opcode = mnem.CmpLTf_bool
					case mnem.ReturnString:
						opcode = mnem.CmpLTstr_bool
					}

				case "<=":
					switch op1.OpCode & mnem.ReturnTypeMask {
					case mnem.ReturnInt:
						opcode = mnem.CmpLEi_bool
					case mnem.ReturnUint:
						opcode = mnem.CmpLEu_bool
					case mnem.ReturnFloat:
						opcode = mnem.CmpLEf_bool
					case mnem.ReturnString:
						opcode = mnem.CmpLEstr_bool
					}

				case ">":
					switch op1.OpCode & mnem.ReturnTypeMask {
					case mnem.ReturnInt:
						opcode = mnem.CmpGTi_bool
					case mnem.ReturnUint:
						opcode = mnem.CmpGTu_bool
					case mnem.ReturnFloat:
						opcode = mnem.CmpGTf_bool
					case mnem.ReturnString:
						opcode = mnem.CmpGTstr_bool
					}

				case ">=":
					switch op1.OpCode & mnem.ReturnTypeMask {
					case mnem.ReturnInt:
						opcode = mnem.CmpGEi_bool
					case mnem.ReturnUint:
						opcode = mnem.CmpGEu_bool
					case mnem.ReturnFloat:
						opcode = mnem.CmpGEf_bool
					case mnem.ReturnString:
						opcode = mnem.CmpGEstr_bool
					}

				default:
					return nil, errors.New(emsg.RuleErr00007 + asts[1].Value.(string))
				}

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					var v bool
					switch opcode & mnem.OpCodeMask {
					case mnem.CmpLTi_bool:
						v = op1.Value.(int64) < op2.Value.(int64)
					case mnem.CmpLEi_bool:
						v = op1.Value.(int64) <= op2.Value.(int64)
					case mnem.CmpGTi_bool:
						v = op1.Value.(int64) > op2.Value.(int64)
					case mnem.CmpGEi_bool:
						v = op1.Value.(int64) >= op2.Value.(int64)
					}
					ret.Value = v

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					var v bool
					switch opcode & mnem.OpCodeMask {
					case mnem.CmpLTu_bool:
						v = op1.Value.(uint64) < op2.Value.(uint64)
					case mnem.CmpLEu_bool:
						v = op1.Value.(uint64) <= op2.Value.(uint64)
					case mnem.CmpGTu_bool:
						v = op1.Value.(uint64) > op2.Value.(uint64)
					case mnem.CmpGEu_bool:
						v = op1.Value.(uint64) >= op2.Value.(uint64)
					}
					ret.Value = v

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					var v bool
					switch opcode & mnem.OpCodeMask {
					case mnem.CmpLTf_bool:
						v = op1.Value.(float64) < op2.Value.(float64)
					case mnem.CmpLEf_bool:
						v = op1.Value.(float64) <= op2.Value.(float64)
					case mnem.CmpGTf_bool:
						v = op1.Value.(float64) > op2.Value.(float64)
					case mnem.CmpGEf_bool:
						v = op1.Value.(float64) >= op2.Value.(float64)
					}
					ret.Value = v

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnString &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnString {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					var v bool
					switch opcode & mnem.OpCodeMask {
					case mnem.CmpLTstr_bool:
						v = op1.Value.(string) < op2.Value.(string)
					case mnem.CmpLEstr_bool:
						v = op1.Value.(string) <= op2.Value.(string)
					case mnem.CmpGTstr_bool:
						v = op1.Value.(string) > op2.Value.(string)
					case mnem.CmpGEstr_bool:
						v = op1.Value.(string) >= op2.Value.(string)
					}
					ret.Value = v

				} else {
					ret = Ast{
						OpCode:         opcode,
						Type:           AstType_AstCons,
						Value:          AstCons{Car: op1, Cdr: op2},
						SourcePosition: op1.SourcePosition,
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
