package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

// Equality `op1 == op2`
// Inequality `op1 != op2`
// Strict Equality `op1 === op2`
// Strict Inequality `op1 !== op2`
var precedence11 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				anyOperand(),
				isOperator(clsz.BinaryOp, []string{"===", "!==", "==", "!="}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				var ops [2]Ast
				var op1, op2 Ast
				var promoted bool

				switch asts[1].Value.(string) {
				case "==", "!=":
					{
						ops, promoted = promoteBinaryOperands(ctx, false, true, AstType_String, asts[0], asts[2])
						op1 = ops[0]
						op2 = ops[1]
					}
				case "===", "!==":
					{
						var ok1, ok2 bool
						op1, ok1 = promoteUnaryOperand(ctx, false, true, AstType_String, asts[0])
						op2, ok2 = promoteUnaryOperand(ctx, false, true, AstType_String, asts[2])
						promoted = ok1 && ok2
					}
				}

				// NOTE: ptr, any, nil, unit value are not promoted.
				// (However, if `forcePromoting` is false, then `ok` will be true.)

				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}
				var opcode AstOpCodeType

				switch asts[1].Value.(string) {
				case "==":
					if promoted {
						switch op1.OpCode & mnem.ReturnTypeMask {
						case mnem.ReturnInt:
							opcode = mnem.CmpEqi_bool
						case mnem.ReturnUint:
							opcode = mnem.CmpEqu_bool
						case mnem.ReturnFloat:
							opcode = mnem.CmpEqf_bool
						case mnem.ReturnBool:
							opcode = mnem.CmpEqbool_bool
						case mnem.ReturnString:
							opcode = mnem.CmpEqstr_bool
						default:
							opcode = mnem.CmpEqdyn_bool
						}
					} else {
						opcode = mnem.CmpEqdyn_bool
					}

				case "===":
					if op1.OpCode&mnem.ReturnTypeMask != 0 {
						if op1.OpCode&mnem.ReturnTypeMask != op2.OpCode&mnem.ReturnTypeMask {
							return AstSlice{FalseAst}, nil
						} else {
							switch op1.OpCode & mnem.ReturnTypeMask {
							case mnem.ReturnInt:
								opcode = mnem.CmpEqi_bool
							case mnem.ReturnUint:
								opcode = mnem.CmpEqu_bool
							case mnem.ReturnFloat:
								opcode = mnem.CmpEqf_bool
							case mnem.ReturnBool:
								opcode = mnem.CmpEqbool_bool
							case mnem.ReturnString:
								opcode = mnem.CmpEqstr_bool
							default:
								opcode = mnem.CmpStrictEqdyn_bool
							}
						}
					} else if op2.OpCode&mnem.ReturnTypeMask != 0 {
						return AstSlice{FalseAst}, nil
					} else {
						opcode = mnem.CmpStrictEqdyn_bool
					}

				case "!=":
					if promoted {
						switch op1.OpCode & mnem.ReturnTypeMask {
						case mnem.ReturnInt:
							opcode = mnem.CmpNotEqi_bool
						case mnem.ReturnUint:
							opcode = mnem.CmpNotEqu_bool
						case mnem.ReturnFloat:
							opcode = mnem.CmpNotEqf_bool
						case mnem.ReturnBool:
							opcode = mnem.CmpNotEqbool_bool
						case mnem.ReturnString:
							opcode = mnem.CmpNotEqstr_bool
						default:
							opcode = mnem.CmpNotEqdyn_bool
						}
					} else {
						opcode = mnem.CmpNotEqdyn_bool
					}

				case "!==":
					if op1.OpCode&mnem.ReturnTypeMask != 0 {
						if op1.OpCode&mnem.ReturnTypeMask != op2.OpCode&mnem.ReturnTypeMask {
							return AstSlice{TrueAst}, nil
						} else {
							switch op1.OpCode & mnem.ReturnTypeMask {
							case mnem.ReturnInt:
								opcode = mnem.CmpNotEqi_bool
							case mnem.ReturnUint:
								opcode = mnem.CmpNotEqu_bool
							case mnem.ReturnFloat:
								opcode = mnem.CmpNotEqf_bool
							case mnem.ReturnBool:
								opcode = mnem.CmpNotEqbool_bool
							case mnem.ReturnString:
								opcode = mnem.CmpNotEqstr_bool
							default:
								opcode = mnem.CmpStrictNotEqdyn_bool
							}
						}
					} else if op2.OpCode&mnem.ReturnTypeMask != 0 {
						return AstSlice{TrueAst}, nil
					} else {
						opcode = mnem.CmpStrictNotEqdyn_bool
					}
				}

				if opcode == 0 {
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
					case mnem.CmpEqi_bool:
						v = op1.Value.(int64) == op2.Value.(int64)

					case mnem.CmpNotEqi_bool:
						v = op1.Value.(int64) != op2.Value.(int64)
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
					case mnem.CmpEqu_bool:
						v = op1.Value.(uint64) == op2.Value.(uint64)

					case mnem.CmpNotEqu_bool:
						v = op1.Value.(uint64) != op2.Value.(uint64)
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
					case mnem.CmpEqf_bool:
						v = op1.Value.(float64) == op2.Value.(float64)

					case mnem.CmpNotEqf_bool:
						v = op1.Value.(float64) != op2.Value.(float64)
					}
					ret.Value = v

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnBool &&
					mnem.Imm_begin < op2.OpCode && op2.OpCode < mnem.Imm_end &&
					op2.OpCode&mnem.ReturnTypeMask == mnem.ReturnBool {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					var v bool
					switch opcode & mnem.OpCodeMask {
					case mnem.CmpEqbool_bool:
						v = op1.Value.(bool) == op2.Value.(bool)

					case mnem.CmpNotEqbool_bool:
						v = op1.Value.(bool) != op2.Value.(bool)
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
					case mnem.CmpEqstr_bool:
						v = op1.Value.(string) == op2.Value.(string)

					case mnem.CmpNotEqstr_bool:
						v = op1.Value.(string) != op2.Value.(string)
					}
					ret.Value = v

				} else if op1.OpCode&mnem.OpCodeMask == mnem.Imm_nil &&
					op2.OpCode&mnem.OpCodeMask == mnem.Imm_nil {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					switch opcode & mnem.OpCodeMask {
					case mnem.CmpEqdyn_bool, mnem.CmpStrictEqdyn_bool:
						ret.Value = true

					case mnem.CmpNotEqdyn_bool, mnem.CmpStrictNotEqdyn_bool:
						ret.Value = false
					}

				} else if op1.OpCode&mnem.OpCodeMask == mnem.Imm_unitval &&
					op2.OpCode&mnem.OpCodeMask == mnem.Imm_unitval {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Bool

					switch opcode & mnem.OpCodeMask {
					case mnem.CmpEqdyn_bool, mnem.CmpStrictEqdyn_bool:
						ret.Value = true

					case mnem.CmpNotEqdyn_bool, mnem.CmpStrictNotEqdyn_bool:
						ret.Value = false
					}

				} else {
					// NOTE: Imm_ptr and Imm_data will be compared at runtime

					ret = Ast{
						OpCode:         opcode,
						Type:           AstType_AstCons,
						Value:          AstCons{Car: op1, Cdr: op2},
						SourcePosition: op1.SourcePosition,
					}
				}

				return AstSlice{ret}, nil
			},
		),
	},
	Rtol: false,
}
