package executor

import (
	"errors"
	"reflect"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func execObjectOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	switch ast.OpCode & mnem.OpCodeMask {
	case mnem.List:
		{
			orig := ast.Value.(AstSlice)
			origLen := len(orig)
			var payload interface{}

			switch ast.OpCode & mnem.ReturnTypeMask {
			case mnem.ReturnInt:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]int32, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = int32(orig[i].Value.(int64))
							slice[i] = int32(*(*int64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				case mnem.Bits16:
					{
						slice := make([]int16, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = int16(orig[i].Value.(int64))
							slice[i] = int16(*(*int64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				case mnem.Bits8:
					{
						slice := make([]int8, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = int8(orig[i].Value.(int64))
							slice[i] = int8(*(*int64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				default:
					{
						slice := make([]int64, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = orig[i].Value.(int64)
							slice[i] = *(*int64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)
						}
						payload = slice
					}
				}
			case mnem.ReturnUint:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]uint32, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = uint32(orig[i].Value.(uint64))
							slice[i] = uint32(*(*uint64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				case mnem.Bits16:
					{
						slice := make([]uint16, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = uint16(orig[i].Value.(uint64))
							slice[i] = uint16(*(*uint64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				case mnem.Bits8:
					{
						slice := make([]uint8, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = uint8(orig[i].Value.(uint64))
							slice[i] = uint8(*(*uint64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				default:
					{
						slice := make([]uint64, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = orig[i].Value.(uint64)
							slice[i] = *(*uint64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)
						}
						payload = slice
					}
				}
			case mnem.ReturnFloat:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]float32, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = float32(orig[i].Value.(float64))
							slice[i] = float32(*(*float64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr))
						}
						payload = slice
					}
				default:
					{
						slice := make([]float64, origLen, origLen)
						for i := 0; i < origLen; i++ {
							// slice[i] = orig[i].Value.(float64)
							slice[i] = *(*float64)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)
						}
						payload = slice
					}
				}
			case mnem.ReturnBool:
				{
					slice := make([]bool, origLen, origLen)
					for i := 0; i < origLen; i++ {
						// slice[i] = orig[i].Value.(bool)
						slice[i] = *(*bool)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)
					}
					payload = slice
				}
			case mnem.ReturnString:
				{
					slice := make([]string, origLen, origLen)
					for i := 0; i < origLen; i++ {
						// slice[i] = orig[i].Value.(string)
						slice[i] = *(*string)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)
					}
					payload = slice
				}
			default:
				{
					slice := make([]interface{}, origLen, origLen)
					for i := 0; i < origLen; i++ {
						slice[i] = orig[i].Value
					}
					payload = slice
				}
			}

			return Ast{
				OpCode: mnem.Imm_data | mnem.Indexable,
				Type:   AstType_ListOfAny,
				Value:  payload,
			}, true, nil, nil
		}

	case mnem.FilledList:
		{
			cons := ast.Value.(AstCons)
			size := cons.Car.Value.(int64)
			var payload interface{}

			switch ast.OpCode & mnem.ReturnTypeMask {
			case mnem.ReturnInt:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]int32, size, size)
						// v := int32(cons.Cdr.Value.(int64))
						v := int32(*(*int64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				case mnem.Bits16:
					{
						slice := make([]int16, size, size)
						// v := int16(cons.Cdr.Value.(int64))
						v := int16(*(*int64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				case mnem.Bits8:
					{
						slice := make([]int8, size, size)
						// v := int8(cons.Cdr.Value.(int64))
						v := int8(*(*int64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				default:
					{
						slice := make([]int64, size, size)
						// v := cons.Cdr.Value.(int64)
						v := *(*int64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr)
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				}
			case mnem.ReturnUint:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]uint32, size, size)
						// v := uint32(cons.Cdr.Value.(uint64))
						v := uint32(*(*uint64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				case mnem.Bits16:
					{
						slice := make([]uint16, size, size)
						// v := uint16(cons.Cdr.Value.(uint64))
						v := uint16(*(*uint64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				case mnem.Bits8:
					{
						slice := make([]uint8, size, size)
						// v := uint8(cons.Cdr.Value.(uint64))
						v := uint8(*(*uint64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				default:
					{
						slice := make([]uint64, size, size)
						// v := cons.Cdr.Value.(uint64)
						v := *(*uint64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr)
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				}
			case mnem.ReturnFloat:
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					{
						slice := make([]float32, size, size)
						// v := float32(cons.Cdr.Value.(float64))
						v := float32(*(*float64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				default:
					{
						slice := make([]float64, size, size)
						// v := cons.Cdr.Value.(float64)
						v := *(*float64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr)
						for i := int64(0); i < size; i++ {
							slice[i] = v
						}
						payload = slice
					}
				}
			case mnem.ReturnBool:
				{
					slice := make([]bool, size, size)
					// v := cons.Cdr.Value.(bool)
					v := *(*bool)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr)
					for i := int64(0); i < size; i++ {
						slice[i] = v
					}
					payload = slice
				}
			case mnem.ReturnString:
				{
					slice := make([]string, size, size)
					// v := cons.Cdr.Value.(string)
					v := *(*string)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr)
					for i := int64(0); i < size; i++ {
						slice[i] = v
					}
					payload = slice
				}
			default:
				{
					slice := make([]interface{}, size, size)
					for i := int64(0); i < size; i++ {
						slice[i] = cons.Cdr.Value
					}
					payload = slice
				}
			}

			return Ast{
				OpCode: mnem.Imm_data | mnem.Indexable,
				Type:   AstType_ListOfAny,
				Value:  payload,
			}, true, nil, nil
		}

	case mnem.Object:
		{
			orig := ast.Value.(AstSlice)
			origLen := len(orig)
			dict := make(map[string]interface{}, origLen/2)

			for i := 0; i < origLen; i += 2 {
				// dict[orig[i].Value.(string)] = orig[i+1].Value
				dict[*(*string)((*rawInterface2)(unsafe.Pointer(&orig[i].Value)).Ptr)] = orig[i+1].Value
			}
			return Ast{
				OpCode: mnem.Imm_data | mnem.Indexable,
				Type:   AstType_ListOfAny,
				Value:  dict,
			}, true, nil, nil
		}

	case mnem.Index:
		{
			opcode := mnem.Imm_data
			astType := AstType_Any
			retType := ast.OpCode & mnem.ReturnTypeMask

			switch retType {
			case mnem.ReturnInt:
				opcode = mnem.Imm_i64
				astType = AstType_Int
			case mnem.ReturnUint:
				opcode = mnem.Imm_u64
				astType = AstType_Uint
			case mnem.ReturnFloat:
				opcode = mnem.Imm_f64
				astType = AstType_Float
			case mnem.ReturnBool:
				opcode = mnem.Imm_bool
				astType = AstType_Bool
			case mnem.ReturnString:
				opcode = mnem.Imm_str
				astType = AstType_String
			}

			cons := ast.Value.(AstCons)
			rv := reflect.ValueOf(cons.Car.Value)
			// i := int(cons.Cdr.Value.(int64))
			i := int(*(*int64)((*rawInterface2)(unsafe.Pointer(&cons.Cdr.Value)).Ptr))

			if 0 <= i && i < rv.Len() {
				rvi := rv.Index(i)
				var v interface{}
				v = rvi.Interface()

				// NOTE: promote to 64bit
				switch ast.OpCode & mnem.BitLenMask {
				case mnem.Bits32:
					switch retType {
					case mnem.ReturnInt:
						// v = int64(v.(int32))
						v = int64(*(*int32)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					case mnem.ReturnUint:
						// v = uint64(v.(uint32))
						v = uint64(*(*uint32)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					case mnem.ReturnFloat:
						// v = float64(v.(float32))
						v = float64(*(*float32)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					}
				case mnem.Bits16:
					switch retType {
					case mnem.ReturnInt:
						// v = int64(v.(int16))
						v = int64(*(*int16)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					case mnem.ReturnUint:
						// v = uint64(v.(uint16))
						v = uint64(*(*uint16)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					}
				case mnem.Bits8:
					switch retType {
					case mnem.ReturnInt:
						// v = int64(v.(int8))
						v = int64(*(*int8)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					case mnem.ReturnUint:
						// v = uint64(v.(uint8))
						v = uint64(*(*uint8)((*rawInterface2)(unsafe.Pointer(&v)).Ptr))
					}
				}
				return Ast{
					OpCode:  opcode,
					Type:    astType,
					Value:   v,
					Address: ReflectionBox{Val: rvi},
				}, true, nil, nil
			} else {
				return *ast, false, nil, errors.New(emsg.ExecErr00017)
			}
		}

	case mnem.Slice:
		{
			opcode := mnem.Imm_data
			astType := AstType_Any
			retType := ast.OpCode & mnem.ReturnTypeMask

			slice := ast.Value.(AstSlice)
			rv := reflect.ValueOf(slice[0].Value)
			rvLen := rv.Len()

			var startTmp, endTmp int64
			var ok bool

			startTmp, ok = slice[1].Value.(int64)
			if !ok {
				startTmp = 0
			}
			endTmp, ok = slice[2].Value.(int64)
			if !ok {
				endTmp = int64(rvLen)
			}

			start := int(startTmp)
			end := int(endTmp)

			if 0 <= start && start <= rvLen && 0 <= end && end <= rvLen && start <= end {
				rvi := rv.Slice(start, end)
				return Ast{
					OpCode:  opcode | retType,
					Type:    astType,
					Value:   rvi.Interface(),
					Address: ReflectionBox{Val: rvi},
				}, true, nil, nil
			} else {
				return *ast, false, nil, errors.New(emsg.ExecErr00017)
			}
		}

	case mnem.Mapindex:
		{
			opcode := mnem.Imm_data
			astType := AstType_Any
			retType := ast.OpCode & mnem.ReturnTypeMask

			switch retType {
			case mnem.ReturnInt:
				opcode = mnem.Imm_i64
				astType = AstType_Int
			case mnem.ReturnUint:
				opcode = mnem.Imm_u64
				astType = AstType_Uint
			case mnem.ReturnFloat:
				opcode = mnem.Imm_f64
				astType = AstType_Float
			case mnem.ReturnBool:
				opcode = mnem.Imm_bool
				astType = AstType_Bool
			case mnem.ReturnString:
				opcode = mnem.Imm_str
				astType = AstType_String
			}

			cons := ast.Value.(AstCons)
			rv := reflect.ValueOf(cons.Car.Value)
			k := reflect.ValueOf(cons.Cdr.Value.(string))
			v := rv.MapIndex(k)

			if v.IsValid() {
				return Ast{
					OpCode:  opcode,
					Type:    astType,
					Value:   v.Interface(),
					Address: MapContainerReflectionBox{Container: rv, Key: k},
				}, true, nil, nil
			} else {
				return Ast{
					OpCode:  opcode,
					Type:    astType,
					Value:   nil, // TODO: Zero value per type
					Address: &NotInitializedMapContainerReflectionBox{Container: rv, Key: k},
				}, true, nil, nil
			}
		}
	}

	return *ast, false, nil, nil
}
