package executor

import (
	"errors"
	"fmt"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func execConvertOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	switch ast.OpCode &^ mnem.FlagsMask {
	case mnem.Convi_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value:  uint64(ast.Value.(AstCons).Car.Value.(int64)),
			Value: uint64(*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convi_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value:  float64(ast.Value.(AstCons).Car.Value.(int64)),
			Value: float64(*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convi_bool:
		// if ast.Value.(AstCons).Car.Value.(int64) != 0 {
		if *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) != 0 {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  true,
			}, true, nil, nil
		} else {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  false,
			}, true, nil, nil
		}

	case mnem.Convu_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value:  int64(ast.Value.(AstCons).Car.Value.(uint64)),
			Value: int64(*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convu_f:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   AstType_Float,
			// Value:  float64(ast.Value.(AstCons).Car.Value.(uint64)),
			Value: float64(*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convu_bool:
		// if ast.Value.(AstCons).Car.Value.(uint64) != 0 {
		if *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) != 0 {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  true,
			}, true, nil, nil
		} else {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  false,
			}, true, nil, nil
		}

	case mnem.Convf_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value:  int64(ast.Value.(AstCons).Car.Value.(float64)),
			Value: int64(*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convf_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value:  uint64(ast.Value.(AstCons).Car.Value.(float64)),
			Value: uint64(*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)),
		}, true, nil, nil
	case mnem.Convf_bool:
		// if ast.Value.(AstCons).Car.Value.(float64) != 0.0 {
		if *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) != 0.0 {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  true,
			}, true, nil, nil
		} else {
			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  false,
			}, true, nil, nil
		}

	case mnem.Convbool_i:
		{
			// w := ast.Value.(AstCons).Car.Value.(bool)
			w := *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)
			var v int64
			if w {
				v = 1
			}
			return Ast{
				OpCode: mnem.Imm_i64,
				Type:   AstType_Int,
				Value:  v,
			}, true, nil, nil
		}
	case mnem.Convbool_u:
		{
			// w := ast.Value.(AstCons).Car.Value.(bool)
			w := *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)
			var v uint64
			if w {
				v = 1
			}
			return Ast{
				OpCode: mnem.Imm_u64,
				Type:   AstType_Uint,
				Value:  v,
			}, true, nil, nil
		}
	case mnem.Convbool_f:
		{
			// w := ast.Value.(AstCons).Car.Value.(bool)
			w := *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr)
			var v float64
			if w {
				v = 1
			}
			return Ast{
				OpCode: mnem.Imm_f64,
				Type:   AstType_Float,
				Value:  v,
			}, true, nil, nil
		}

	case mnem.Convstr_i, mnem.Convdyn_i:
		// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
		{
			// out, ok := CastNumberLiteral(ast.Value.(AstCons).Car, AstType_Int)
			out, ok := CastNumberLiteral((*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, AstType_Int)
			if ok {
				return out, true, nil, nil
			} else {
				return *ast, true, nil, errors.New(
					emsg.ExecErr00006 + fmt.Sprintf("%v", ast.Value),
				)
			}
		}
	case mnem.Convstr_u, mnem.Convdyn_u:
		// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
		{
			// out, ok := CastNumberLiteral(ast.Value.(AstCons).Car, AstType_Uint)
			out, ok := CastNumberLiteral((*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, AstType_Uint)
			if ok {
				return out, true, nil, nil
			} else {
				return *ast, true, nil, errors.New(
					emsg.ExecErr00007 + fmt.Sprintf("%v", ast.Value),
				)
			}
		}
	case mnem.Convstr_f, mnem.Convdyn_f:
		// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
		{
			// out, ok := CastNumberLiteral(ast.Value.(AstCons).Car, AstType_Float)
			out, ok := CastNumberLiteral((*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, AstType_Float)
			if ok {
				return out, true, nil, nil
			} else {
				return *ast, true, nil, errors.New(
					emsg.ExecErr00008 + fmt.Sprintf("%v", ast.Value),
				)
			}
		}
	case mnem.Convstr_bool, mnem.Convdyn_bool:
		// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
		{
			// out, ok := CastNumberLiteral(ast.Value.(AstCons).Car, AstType_Bool)
			out, ok := CastNumberLiteral((*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, AstType_Bool)
			if ok {
				return out, true, nil, nil
			} else {
				return *ast, true, nil, errors.New(
					emsg.ExecErr00009 + fmt.Sprintf("%v", ast.Value),
				)
			}
		}
	case mnem.Convi_str, mnem.Convu_str, mnem.Convf_str, mnem.Convbool_str,
		mnem.Convdyn_str:
		// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
		{
			// out, ok := CastNumberLiteral(ast.Value.(AstCons).Car, AstType_String)
			out, ok := CastNumberLiteral((*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, AstType_String)
			if ok {
				return out, true, nil, nil
			} else {
				return *ast, true, nil, errors.New(
					emsg.ExecErr00010 + fmt.Sprintf("%v", ast.Value),
				)
			}
		}

		// TODO: convdynforce
		// TODO: [u8] <-> String conversion
		// TODO: complex type casting
	}

	return *ast, false, nil, nil
}
