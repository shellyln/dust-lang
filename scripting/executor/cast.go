package executor

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func CastNumberLiteral(ast Ast, to AstType) (Ast, bool) {
	if ast.Type == to {
		return ast, true
	}

	ast.Address = nil
	ok := true

	// TODO: BUG: Set type id.

	switch to {
	case AstType_Rune:
		switch ast.Type {
		case AstType_Int:
			// ast.Value = rune(ast.Value.(rune))
			ast.Value = rune(*(*int64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Uint:
			// ast.Value = rune(ast.Value.(uint64))
			ast.Value = rune(*(*uint64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Float:
			// ast.Value = rune(ast.Value.(float64))
			ast.Value = rune(*(*float64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Bool:
			// if ast.Value.(bool) {
			if *(*bool)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) {
				ast.Value = rune(1)
			} else {
				ast.Value = rune(0)
			}
		case AstType_String:
			{
				// v, err := strconv.ParseInt(ast.Value.(string), 10, 64)
				v, err := strconv.ParseInt(*(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr), 10, 64)
				if err != nil {
					ast.Value = rune(0)
				} else {
					ast.Value = rune(v)
				}
			}
		default:
			{
				rv := reflect.ValueOf(ast.Value)
				switch rv.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ast.Value = rune(rv.Int())
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					ast.Value = rune(rv.Uint())
				case reflect.Float64, reflect.Float32:
					ast.Value = rune(rv.Float())
				case reflect.Bool:
					if rv.Bool() {
						ast.Value = rune(1)
					} else {
						ast.Value = rune(0)
					}
				case reflect.String:
					v, err := strconv.ParseInt(rv.String(), 10, 64)
					if err != nil {
						ast.Value = rune(0)
					} else {
						ast.Value = rune(v)
					}
				default:
					ast.Value = rune(0)
					ok = false
				}
			}
		}

	case AstType_Int:
		switch ast.Type {
		case AstType_Rune:
			// ast.Value = int64(ast.Value.(rune))
			ast.Value = int64(*(*rune)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Uint:
			// ast.Value = int64(ast.Value.(uint64))
			ast.Value = int64(*(*uint64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Float:
			// ast.Value = int64(ast.Value.(float64))
			ast.Value = int64(*(*float64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Bool:
			// if ast.Value.(bool) {
			if *(*bool)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) {
				ast.Value = int64(1)
			} else {
				ast.Value = int64(0)
			}
		case AstType_String:
			{
				// v, err := strconv.ParseInt(ast.Value.(string), 10, 64)
				v, err := strconv.ParseInt(*(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr), 10, 64)
				if err != nil {
					ast.Value = int64(0)
				} else {
					ast.Value = v
				}
			}
		default:
			{
				rv := reflect.ValueOf(ast.Value)
				switch rv.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ast.Value = rv.Int()
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					ast.Value = int64(rv.Uint())
				case reflect.Float64, reflect.Float32:
					ast.Value = int64(rv.Float())
				case reflect.Bool:
					if rv.Bool() {
						ast.Value = int64(1)
					} else {
						ast.Value = int64(0)
					}
				case reflect.String:
					v, err := strconv.ParseInt(rv.String(), 10, 64)
					if err != nil {
						ast.Value = int64(0)
					} else {
						ast.Value = int64(v)
					}
				default:
					ast.Value = int64(0)
					ok = false
				}
			}
		}

	case AstType_Uint:
		switch ast.Type {
		case AstType_Rune:
			// ast.Value = uint64(ast.Value.(rune))
			ast.Value = uint64(*(*rune)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Int:
			// ast.Value = uint64(ast.Value.(int64))
			ast.Value = uint64(*(*int64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Float:
			// ast.Value = uint64(ast.Value.(float64))
			ast.Value = uint64(*(*float64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Bool:
			// if ast.Value.(bool) {
			if *(*bool)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) {
				ast.Value = uint64(1)
			} else {
				ast.Value = uint64(0)
			}
		case AstType_String:
			{
				// v, err := strconv.ParseUint(ast.Value.(string), 10, 64)
				v, err := strconv.ParseUint(*(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr), 10, 64)
				if err != nil {
					ast.Value = uint64(0)
				} else {
					ast.Value = v
				}
			}
		default:
			{
				rv := reflect.ValueOf(ast.Value)
				switch rv.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ast.Value = uint64(rv.Int())
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					ast.Value = rv.Uint()
				case reflect.Float64, reflect.Float32:
					ast.Value = uint64(rv.Float())
				case reflect.Bool:
					if rv.Bool() {
						ast.Value = uint64(1)
					} else {
						ast.Value = uint64(0)
					}
				case reflect.String:
					v, err := strconv.ParseUint(rv.String(), 10, 64)
					if err != nil {
						ast.Value = uint64(0)
					} else {
						ast.Value = uint64(v)
					}
				default:
					ast.Value = uint64(0)
					ok = false
				}
			}
		}

	case AstType_Float:
		switch ast.Type {
		case AstType_Rune:
			// ast.Value = float64(ast.Value.(rune))
			ast.Value = float64(*(*rune)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Int:
			// ast.Value = float64(ast.Value.(int64))
			ast.Value = float64(*(*int64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Uint:
			// ast.Value = float64(ast.Value.(uint64))
			ast.Value = float64(*(*uint64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr))
		case AstType_Bool:
			// if ast.Value.(bool) {
			if *(*bool)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) {
				ast.Value = 1.0
			} else {
				ast.Value = 0.0
			}
		case AstType_String:
			{
				// v, err := strconv.ParseFloat(ast.Value.(string), 64)
				v, err := strconv.ParseFloat(*(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr), 64)
				if err != nil {
					ast.Value = float64(0)
				} else {
					ast.Value = v
				}
			}
		default:
			{
				rv := reflect.ValueOf(ast.Value)
				switch rv.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					ast.Value = float64(rv.Int())
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					ast.Value = float64(rv.Uint())
				case reflect.Float64, reflect.Float32:
					ast.Value = rv.Float()
				case reflect.Bool:
					if rv.Bool() {
						ast.Value = float64(1)
					} else {
						ast.Value = float64(0)
					}
				case reflect.String:
					v, err := strconv.ParseFloat(rv.String(), 64)
					if err != nil {
						ast.Value = float64(0)
					} else {
						ast.Value = float64(v)
					}
				default:
					ast.Value = float64(0)
					ok = false
				}
			}
		}

	case AstType_Bool:
		switch ast.Type {
		case AstType_Rune:
			// if ast.Value.(rune) != 0 {
			if *(*rune)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) != 0 {
				ast.Value = true
			} else {
				ast.Value = false
			}
		case AstType_Int:
			// if ast.Value.(int64) != 0 {
			if *(*int64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) != 0 {
				ast.Value = true
			} else {
				ast.Value = false
			}
		case AstType_Uint:
			// if ast.Value.(uint64) != 0 {
			if *(*uint64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) != 0 {
				ast.Value = true
			} else {
				ast.Value = false
			}
		case AstType_Float:
			// if ast.Value.(float64) != 0.0 {
			if *(*float64)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr) != 0.0 {
				ast.Value = true
			} else {
				ast.Value = false
			}
		case AstType_String:
			{
				// v := ast.Value.(string)
				v := *(*string)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
				if v != "" {
					ast.Value = true
				} else {
					ast.Value = false
				}
			}
		default:
			{
				rv := reflect.ValueOf(ast.Value)
				switch rv.Kind() {
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					if rv.Int() != 0 {
						ast.Value = true
					} else {
						ast.Value = false
					}
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					if rv.Uint() != 0 {
						ast.Value = true
					} else {
						ast.Value = false
					}
				case reflect.Float64, reflect.Float32:
					if rv.Float() != 0.0 {
						ast.Value = true
					} else {
						ast.Value = false
					}
				case reflect.Bool:
					ast.Value = rv.Bool()
				case reflect.String:
					if rv.String() != "" {
						ast.Value = true
					} else {
						ast.Value = false
					}
				default:
					ast.Value = false
					ok = false
				}
			}
		}

	case AstType_String:
		ast.Value = fmt.Sprintf("%v", ast.Value)

	case AstType_Any:
		ast.OpCode = mnem.Imm_data
		ast.Type = AstType_Any
		return ast, true

	default:
		return ast, false
	}

	// NOTE: clear `lvalue` flag
	// TODO: BUG: case of imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef
	ast.OpCode = (ast.OpCode & mnem.OpCodeMask) | AstOpCodeType(to)
	ast.Type = to
	return ast, ok
}
