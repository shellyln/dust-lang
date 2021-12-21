package executor

import (
	"errors"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func execAssignOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	switch ast.OpCode & mnem.OpCodeMask {
	case mnem.Assign:
		{
			opcode := mnem.Imm_data
			astType := AstType_Any
			retType := ast.OpCode & mnem.ReturnTypeMask

			// op1 := ast.Value.(AstCons).Car
			// op2 := ast.Value.(AstCons).Cdr
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op2 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr

			if ast.OpCode&mnem.Callable != 0 {
				op1.Address.SetAny(op2.Value)
			} else if ast.OpCode&mnem.Maybe != 0 {
				op1.Address.SetAny(op2.Value)
			} else if ast.OpCode&mnem.Indexable != 0 {
				op1.Address.SetAny(op2.Value)
			} else {
				switch retType {
				case mnem.ReturnInt:
					opcode = mnem.Imm_i64
					astType = AstType_Int
					// op1.Address.SetInt(op2.Value.(int64))
					op1.Address.SetInt(*(*int64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnUint:
					opcode = mnem.Imm_u64
					astType = AstType_Uint
					// op1.Address.SetUint(op2.Value.(uint64))
					op1.Address.SetUint(*(*uint64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnFloat:
					opcode = mnem.Imm_f64
					astType = AstType_Float
					// op1.Address.SetFloat(op2.Value.(float64))
					op1.Address.SetFloat(*(*float64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnBool:
					opcode = mnem.Imm_bool
					astType = AstType_Bool
					// op1.Address.SetBool(op2.Value.(bool))
					op1.Address.SetBool(*(*bool)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnString:
					opcode = mnem.Imm_str
					astType = AstType_String
					// op1.Address.SetString(op2.Value.(string))
					op1.Address.SetString(*(*string)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				default:
					op1.Address.SetAny(op2.Value)
				}
			}

			return Ast{
				OpCode:  opcode | mnem.Lvalue,
				Type:    astType,
				Value:   op2.Value,
				Address: op1.Address,
			}, true, nil, nil
		}

	case mnem.DefVar, mnem.DefConst:
		{
			opcode := mnem.Imm_data
			astType := AstType_Any
			retType := ast.OpCode & mnem.ReturnTypeMask

			// op1 := ast.Value.(AstCons).Car
			// op2 := ast.Value.(AstCons).Cdr
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op2 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr

			flags := (ast.OpCode & mnem.FlagsMask) &^ mnem.Lvalue
			if ast.OpCode&mnem.OpCodeMask == mnem.DefVar {
				flags |= mnem.Lvalue
			}

			// vi, ok := ctx.DefineVariable(op1.Value.(string), VariableInfo{
			vi, ok := ctx.DefineVariable(*(*string)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr), VariableInfo{
				Flags: retType | flags,
			})
			if !ok {
				return *ast, true, nil, errors.New(emsg.ExecErr00002 + op1.Value.(string))
			}

			if ast.OpCode&mnem.Callable != 0 {
				vi.SetAny(op2.Value)
			} else if ast.OpCode&mnem.Maybe != 0 {
				vi.SetAny(op2.Value)
			} else if ast.OpCode&mnem.Indexable != 0 {
				vi.SetAny(op2.Value)
			} else {
				switch retType {
				case mnem.ReturnInt:
					opcode = mnem.Imm_i64
					astType = AstType_Int
					// vi.SetInt(op2.Value.(int64))
					vi.SetInt(*(*int64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnUint:
					opcode = mnem.Imm_u64
					astType = AstType_Uint
					// vi.SetUint(op2.Value.(uint64))
					vi.SetUint(*(*uint64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnFloat:
					opcode = mnem.Imm_f64
					astType = AstType_Float
					// vi.SetFloat(op2.Value.(float64))
					vi.SetFloat(*(*float64)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnBool:
					opcode = mnem.Imm_bool
					astType = AstType_Bool
					// vi.SetBool(op2.Value.(bool))
					vi.SetBool(*(*bool)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				case mnem.ReturnString:
					opcode = mnem.Imm_str
					astType = AstType_String
					// vi.SetString(op2.Value.(string))
					vi.SetString(*(*string)((*rawInterface2)(unsafe.Pointer(&op2.Value)).Ptr))
				default:
					vi.SetAny(op2.Value)
				}
			}

			return Ast{
				OpCode:  opcode | mnem.Lvalue,
				Type:    astType,
				Value:   op2.Value,
				Address: vi,
			}, true, nil, nil
		}
	}

	return *ast, false, nil, nil
}
