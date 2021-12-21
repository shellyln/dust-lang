package executor

import (
	"errors"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func execValueOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	switch ast.OpCode &^ mnem.FlagsMask {
	case mnem.Imm_i64, mnem.Imm_u64, mnem.Imm_f64,
		mnem.Imm_bool, mnem.Imm_str,
		mnem.Imm_ptr, mnem.Imm_data,
		mnem.Imm_nil, mnem.Imm_unitval:
		return *ast, true, nil, nil

	case mnem.Imm_i32, mnem.Imm_i16, mnem.Imm_i8:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   ast.Type,
			Value:  ast.Value,
		}, true, nil, nil
	case mnem.Imm_u32, mnem.Imm_u16, mnem.Imm_u8:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   ast.Type,
			Value:  ast.Value,
		}, true, nil, nil
	case mnem.Imm_f32:
		return Ast{
			OpCode: mnem.Imm_f64,
			Type:   ast.Type,
			Value:  ast.Value,
		}, true, nil, nil

	case mnem.Symbol, mnem.Symbol_i, mnem.Symbol_u, mnem.Symbol_f,
		mnem.Symbol_bool, mnem.Symbol_str:
		{
			// vi, ok := ctx.GetVariableInfo(ast.Value.(string))
			vi, ok := ctx.GetVariableInfo(*(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
			if ok {
				lvalueFlag := mnem.Lvalue
				if vi.Flags&mnem.Lvalue == 0 {
					lvalueFlag = 0
				}
				if vi.Flags&mnem.Callable != 0 || vi.Flags&mnem.Maybe != 0 || vi.Flags&mnem.Indexable != 0 {
					return Ast{
						OpCode:  mnem.Imm_data | lvalueFlag,
						Type:    AstType_Any,
						Value:   vi.Value,
						Address: vi,
					}, true, nil, nil
				}

				v := vi.Value

				if v == nil {
					return Ast{
						OpCode:  mnem.Imm_nil,
						Type:    AstType_Nil,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				}

				switch vi.Flags & mnem.ReturnTypeMask {
				case mnem.ReturnInt:
					// switch vi.Flags & mnem.BitLenMask {
					// case mnem.Bits32:
					// 	// v = int32(v.(int64))
					// 	v = int32(*(*int64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// case mnem.Bits16:
					// 	// v = int16(v.(int64))
					// 	v = int16(*(*int64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// case mnem.Bits8:
					// 	// v = int8(v.(int64))
					// 	v = int8(*(*int64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// }

					// NOTE: promote to 64bit
					return Ast{
						OpCode:  mnem.Imm_i64 | lvalueFlag | vi.Flags&mnem.BitLenMask,
						Type:    AstType_Int,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				case mnem.ReturnUint:
					// switch vi.Flags & mnem.BitLenMask {
					// case mnem.Bits32:
					// 	// v = uint32(v.(uint64))
					// 	v = uint32(*(*uint64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// case mnem.Bits16:
					// 	// v = uint16(v.(uint64))
					// 	v = uint16(*(*uint64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// case mnem.Bits8:
					// 	// v = uint8(v.(uint64))
					// 	v = uint8(*(*uint64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// }

					// NOTE: promote to 64bit
					return Ast{
						OpCode:  mnem.Imm_u64 | lvalueFlag | vi.Flags&mnem.BitLenMask,
						Type:    AstType_Uint,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				case mnem.ReturnFloat:
					// switch vi.Flags & mnem.BitLenMask {
					// case mnem.Bits32:
					// 	// v = float32(v.(float64))
					// 	v = float32(*(*float64)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					// }

					// NOTE: promote to 64bit
					return Ast{
						OpCode:  mnem.Imm_f64 | lvalueFlag | vi.Flags&mnem.BitLenMask,
						Type:    AstType_Float,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				case mnem.ReturnBool:
					return Ast{
						OpCode:  mnem.Imm_bool | lvalueFlag,
						Type:    AstType_Bool,
						Value:   vi.Value,
						Address: vi,
					}, true, nil, nil
				case mnem.ReturnString:
					return Ast{
						OpCode:  mnem.Imm_str | lvalueFlag,
						Type:    AstType_String,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				case mnem.ReturnAny:
					return Ast{
						OpCode:  mnem.Imm_data | lvalueFlag | vi.Flags&mnem.BitLenMask,
						Type:    AstType_Any,
						Value:   v,
						Address: vi,
					}, true, nil, nil
				default:
					return *ast, true, nil, errors.New(emsg.ExecErr00019)
				}
			} else {
				return *ast, true, nil, errors.New(emsg.ExecErr00020 + ast.Value.(string))
			}
		}
	}

	return *ast, false, nil, nil
}
