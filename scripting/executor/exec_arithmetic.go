package executor

import (
	"math"
	"unsafe"

	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func execArithmeticOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	switch ast.OpCode &^ mnem.FlagsMask {
	case mnem.PreIncr_i:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetInt(op1.Address.GetInt() + 1)
			return Ast{
				OpCode: mnem.Imm_i64 | mnem.Lvalue,
				Type:   AstType_Int,
				// Value:   op1.Value.(int64) + 1,
				Value:   *(*int64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) + 1,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PreIncr_u:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetUint(op1.Address.GetUint() + 1)
			return Ast{
				OpCode: mnem.Imm_u64 | mnem.Lvalue,
				Type:   AstType_Uint,
				// Value:   op1.Value.(uint64) + 1,
				Value:   *(*uint64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) + 1,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PreIncr_f:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetFloat(op1.Address.GetFloat() + 1.0)
			return Ast{
				OpCode: mnem.Imm_f64 | mnem.Lvalue,
				Type:   AstType_Float,
				// Value:   op1.Value.(float64) + 1.0,
				Value:   *(*float64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) + 1.0,
				Address: op1.Address,
			}, true, nil, nil
		}

	case mnem.PreDecr_i:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetInt(op1.Address.GetInt() - 1)
			return Ast{
				OpCode: mnem.Imm_i64 | mnem.Lvalue,
				Type:   AstType_Int,
				// Value:   op1.Value.(int64) - 1,
				Value:   *(*int64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) - 1,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PreDecr_u:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetUint(op1.Address.GetUint() - 1)
			return Ast{
				OpCode: mnem.Imm_u64 | mnem.Lvalue,
				Type:   AstType_Uint,
				// Value:   op1.Value.(uint64) - 1,
				Value:   *(*uint64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) - 1,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PreDecr_f:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetFloat(op1.Address.GetFloat() - 1.0)
			return Ast{
				OpCode: mnem.Imm_f64 | mnem.Lvalue,
				Type:   AstType_Float,
				// Value:   op1.Value.(float64) - 1.0,
				Value:   *(*float64)((*rawInterface2)(unsafe.Pointer(&op1.Value)).Ptr) - 1.0,
				Address: op1.Address,
			}, true, nil, nil
		}

	case mnem.PostIncr_i:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetInt(op1.Address.GetInt() + 1)
			return Ast{
				OpCode:  mnem.Imm_i64 | mnem.Lvalue,
				Type:    AstType_Int,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PostIncr_u:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetUint(op1.Address.GetUint() + 1)
			return Ast{
				OpCode:  mnem.Imm_u64 | mnem.Lvalue,
				Type:    AstType_Uint,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PostIncr_f:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetFloat(op1.Address.GetFloat() + 1.0)
			return Ast{
				OpCode:  mnem.Imm_f64 | mnem.Lvalue,
				Type:    AstType_Float,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}

	case mnem.PostDecr_i:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetInt(op1.Address.GetInt() - 1)
			return Ast{
				OpCode:  mnem.Imm_i64 | mnem.Lvalue,
				Type:    AstType_Int,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PostDecr_u:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetUint(op1.Address.GetUint() - 1)
			return Ast{
				OpCode:  mnem.Imm_u64 | mnem.Lvalue,
				Type:    AstType_Uint,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}
	case mnem.PostDecr_f:
		{
			// op1 := ast.Value.(AstCons).Car
			op1 := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car
			op1.Address.SetFloat(op1.Address.GetFloat() - 1.0)
			return Ast{
				OpCode:  mnem.Imm_f64 | mnem.Lvalue,
				Type:    AstType_Float,
				Value:   op1.Value,
				Address: op1.Address,
			}, true, nil, nil
		}

	case mnem.Neg_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value:  -ast.Value.(AstCons).Car.Value.(int64),
			Value: -*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil
	case mnem.Neg_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value:  -ast.Value.(AstCons).Car.Value.(uint64),
			Value: -*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil
	case mnem.Neg_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value:  -ast.Value.(AstCons).Car.Value.(float64),
			Value: -*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil

	case mnem.Pow_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			Value: math.Pow(
				// ast.Value.(AstCons).Car.Value.(float64),
				// ast.Value.(AstCons).Cdr.Value.(float64),
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
			),
		}, true, nil, nil

	case mnem.Mul_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) *
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) *
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Mul_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) *
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) *
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Mul_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value: ast.Value.(AstCons).Car.Value.(float64) *
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) *
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.Div_i:
		// {
		// 	op2 := ast.Value.(AstCons).Cdr.Value.(int64)
		// 	if op2 != 0 {
		// 		return Ast{
		// 			OpCode: mnem.Imm_i64,
		// 			Type:   AstType_Int,
		// 			Value:  ast.Value.(AstCons).Car.Value.(int64) / op2,
		// 		}, true, nil, nil
		// 	} else {
		// 		return *ast, true, nil, errors.New(emsg.ExecErr00001)
		// 	}
		// }
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) /
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Div_u:
		// {
		// 	op2 := ast.Value.(AstCons).Cdr.Value.(uint64)
		// 	if op2 != 0 {
		// 		return Ast{
		// 			OpCode: mnem.Imm_u64,
		// 			Type:   AstType_Uint,
		// 			Value:  ast.Value.(AstCons).Car.Value.(uint64) / op2,
		// 		}, true, nil, nil
		// 	} else {
		// 		return *ast, true, nil, errors.New(emsg.ExecErr00001)
		// 	}
		// }
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) /
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Div_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value: ast.Value.(AstCons).Car.Value.(float64) /
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) /
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.Mod_i:
		// {
		// 	op2 := ast.Value.(AstCons).Cdr.Value.(int64)
		// 	if op2 != 0 {
		// 		return Ast{
		// 			OpCode: mnem.Imm_i64,
		// 			Type:   AstType_Int,
		// 			Value:  ast.Value.(AstCons).Car.Value.(int64) % op2,
		// 		}, true, nil, nil
		// 	} else {
		// 		return *ast, true, nil, errors.New(emsg.ExecErr00001)
		// 	}
		// }
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) %
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Mod_u:
		// {
		// 	op2 := ast.Value.(AstCons).Cdr.Value.(uint64)
		// 	if op2 != 0 {
		// 		return Ast{
		// 			OpCode: mnem.Imm_u64,
		// 			Type:   AstType_Uint,
		// 			Value:  ast.Value.(AstCons).Car.Value.(uint64) % op2,
		// 		}, true, nil, nil
		// 	} else {
		// 		return *ast, true, nil, errors.New(emsg.ExecErr00001)
		// 	}
		// }
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) %
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Mod_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			Value: math.Mod(
				// ast.Value.(AstCons).Car.Value.(float64),
				// ast.Value.(AstCons).Cdr.Value.(float64),
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
			),
		}, true, nil, nil

	case mnem.Add_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) +
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) +
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Add_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) +
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) +
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Add_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value: ast.Value.(AstCons).Car.Value.(float64) +
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) +
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.Sub_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) -
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) -
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Sub_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) -
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) -
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.Sub_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value: ast.Value.(AstCons).Car.Value.(float64) -
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) -
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.Concat_str:
		return Ast{
			OpCode: mnem.Imm_str,
			Type:   AstType_String,
			// Value: ast.Value.(AstCons).Car.Value.(string) +
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) +
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	}

	return *ast, false, nil, nil
}
