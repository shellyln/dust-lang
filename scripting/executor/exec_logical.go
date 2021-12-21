package executor

import (
	"reflect"
	"strconv"
	"unsafe"

	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

//
func execBitwiseAndLogicalOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	opcode := ast.OpCode &^ mnem.FlagsMask

	switch opcode {
	case mnem.LogicalNotbool_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value:  !ast.Value.(AstCons).Car.Value.(bool),
			Value: !*(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil

	case mnem.LogicalAndbool_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(bool) &&
			// 	ast.Value.(AstCons).Cdr.Value.(bool),
			Value: *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) &&
				*(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.LogicalOrbool_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(bool) ||
			// 	ast.Value.(AstCons).Cdr.Value.(bool),
			Value: *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ||
				*(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseNot_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value:  ^ast.Value.(AstCons).Car.Value.(int64),
			Value: ^*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseNot_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value:  ^ast.Value.(AstCons).Car.Value.(uint64),
			Value: ^*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseLShift_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) <<
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <<
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseLShift_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) <<
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <<
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseSRShift_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) >>
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >>
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseSRShift_u, mnem.BitwiseURShift_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) >>
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >>
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseURShift_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: int64(uint64(ast.Value.(AstCons).Car.Value.(int64)) >>
			// 	ast.Value.(AstCons).Cdr.Value.(uint64)),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >>
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseAnd_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) &
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) &
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseAnd_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) &
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) &
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseXor_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) ^
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ^
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseXor_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) ^
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ^
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.BitwiseOr_i:
		return Ast{
			OpCode: mnem.Imm_i64,
			Type:   AstType_Int,
			// Value: ast.Value.(AstCons).Car.Value.(int64) |
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) |
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.BitwiseOr_u:
		return Ast{
			OpCode: mnem.Imm_u64,
			Type:   AstType_Uint,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) |
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) |
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpEqi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) ==
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ==
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpEqu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) ==
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ==
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpEqf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) ==
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ==
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpEqbool_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(bool) ==
			// 	ast.Value.(AstCons).Cdr.Value.(bool),
			Value: *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ==
				*(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpEqstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) ==
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) ==
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpNotEqi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) !=
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) !=
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpNotEqu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) !=
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) !=
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpNotEqf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) !=
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) !=
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpNotEqbool_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(bool) !=
			// 	ast.Value.(AstCons).Cdr.Value.(bool),
			Value: *(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) !=
				*(*bool)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpNotEqstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) !=
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) !=
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpEqdyn_bool, mnem.CmpNotEqdyn_bool, mnem.CmpStrictEqdyn_bool, mnem.CmpStrictNotEqdyn_bool:
		{
			var result bool
			// op1any := ast.Value.(AstCons).Car.Value
			// op2any := ast.Value.(AstCons).Cdr.Value
			op1any := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value
			op2any := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value

			if op1any == nil && op2any == nil {
				if opcode == mnem.CmpEqdyn_bool || opcode == mnem.CmpStrictEqdyn_bool {
					result = true
				}
				return Ast{
					OpCode: mnem.Imm_bool,
					Type:   AstType_Bool,
					Value:  result,
				}, true, nil, nil
			}

			var ok1i, ok1u, ok1f, ok1bool, ok1str, ok1unitval, ok1nil bool
			var ok2i, ok2u, ok2f, ok2bool, ok2str, ok2unitval, ok2nil bool
			var op1i, op2i int64
			var op1u, op2u uint64
			var op1f, op2f float64
			var op1bool, op2bool bool
			var op1str, op2str string

			switch val := op1any.(type) {
			case int64:
				op1i = val
				ok1i = true
			case uint64:
				op1u = val
				ok1u = true
			case float64:
				op1f = val
				ok1f = true
			case bool:
				op1bool = val
				ok1bool = true
			case string:
				op1str = val
				ok1str = true
			case *Unit:
				ok1unitval = true
			case nil:
				ok1nil = true
			}

			switch val := op2any.(type) {
			case int64:
				op2i = val
				ok2i = true
			case uint64:
				op2u = val
				ok2u = true
			case float64:
				op2f = val
				ok2f = true
			case bool:
				op2bool = val
				ok2bool = true
			case string:
				op2str = val
				ok2str = true
			case *Unit:
				ok2unitval = true
			case nil:
				ok2nil = true
			}

			if opcode == mnem.CmpEqdyn_bool || opcode == mnem.CmpNotEqdyn_bool {
				if ok1i {
					op1u = uint64(op1i)
					ok1u = true
					if !ok1str {
						op1str = strconv.FormatInt(op1i, 10)
					}
				}
				if ok2i {
					op2u = uint64(op2i)
					ok2u = true
					if !ok2str {
						op2str = strconv.FormatInt(op2i, 10)
					}
				}

				if ok1u {
					op1f = float64(op1u)
					ok1f = true
					if !ok1str {
						op1str = strconv.FormatUint(op1u, 10)
					}
				}
				if ok2u {
					op2f = float64(op2u)
					ok2f = true
					if !ok2str {
						op2str = strconv.FormatUint(op2u, 10)
					}
				}

				if ok1f {
					op1bool = op1f != 0.0
					ok1bool = true
					if !ok1str {
						op1str = strconv.FormatFloat(op1f, 'g', -1, 64)
					}
				}
				if ok2f {
					op2bool = op2f != 0.0
					ok2bool = true
					if !ok2str {
						op2str = strconv.FormatFloat(op2f, 'g', -1, 64)
					}
				}

				if ok1bool {
					if !ok1str {
						if ok1bool {
							op1str = "true"
						} else {
							op1str = "false"
						}
					}
				}
				if ok2bool {
					if !ok2str {
						if ok2bool {
							op2str = "true"
						} else {
							op2str = "false"
						}
					}
				}
			}

			if ok1i && ok2i {
				result = op1i == op2i
			} else if ok1u && ok2u {
				result = op1u == op2u
			} else if ok1f && ok2f {
				result = op1f == op2f
			} else if ok1bool && ok2bool {
				result = op1bool == op2bool
			} else if ok1str && ok2str {
				result = op1str == op2str
			} else if ok1unitval && ok2unitval {
				result = true
			} else if ok1nil && ok2nil {
				result = true
			} else {
				// NOTE: DO NOT let it by value!
				// op1raw := (*rawInterface)(unsafe.Pointer(&op1any))
				// op2raw := (*rawInterface)(unsafe.Pointer(&op2any))
				op1raw := (*rawInterface2)(unsafe.Pointer(&op1any))
				op2raw := (*rawInterface2)(unsafe.Pointer(&op2any))

				if reflect.ValueOf(op1any).Kind() == reflect.Slice &&
					reflect.ValueOf(op2any).Kind() == reflect.Slice {
					// NOTE: Two slices that point same memory comparison is always false
					//       due to slice header structure.

					// NOTE: DO NOT let it by value!
					// NOTE: `.Ptr` is referenced and tracked via `interface{}`. This is just a view.
					// NOTE: Possible invalid pointer due to memory relocation (owing to future GC algorithm changes).
					// op1SliceRaw := (*rawSliceHeader)(unsafe.Pointer(op1raw.Ptr)) // NOTE: go vet warning;
					// op2SliceRaw := (*rawSliceHeader)(unsafe.Pointer(op2raw.Ptr)) // NOTE: go vet warning
					op1SliceRaw := (*rawSliceHeader)(op1raw.Ptr)
					op2SliceRaw := (*rawSliceHeader)(op2raw.Ptr)

					result = op1SliceRaw.Data == op2SliceRaw.Data &&
						op1SliceRaw.Len == op2SliceRaw.Len &&
						op1SliceRaw.Cap == op2SliceRaw.Cap
				} else {
					result = op1raw.Ptr == op2raw.Ptr
				}
			}

			if opcode == mnem.CmpNotEqdyn_bool || opcode == mnem.CmpStrictNotEqdyn_bool {
				result = !result
			}

			return Ast{
				OpCode: mnem.Imm_bool,
				Type:   AstType_Bool,
				Value:  result,
			}, true, nil, nil
		}

	case mnem.CmpLTi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) <
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLTu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) <
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLTf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) <
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLTstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) <
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpLEi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) <=
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <=
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLEu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) <=
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <=
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLEf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) <=
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <=
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpLEstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) <=
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) <=
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpGTi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) >
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGTu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) >
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGTf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) >
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGTstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) >
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil

	case mnem.CmpGEi_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(int64) >=
			// 	ast.Value.(AstCons).Cdr.Value.(int64),
			Value: *(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >=
				*(*int64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGEu_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(uint64) >=
			// 	ast.Value.(AstCons).Cdr.Value.(uint64),
			Value: *(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >=
				*(*uint64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGEf_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(float64) >=
			// 	ast.Value.(AstCons).Cdr.Value.(float64),
			Value: *(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >=
				*(*float64)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	case mnem.CmpGEstr_bool:
		return Ast{
			OpCode: mnem.Imm_bool,
			Type:   AstType_Bool,
			// Value: ast.Value.(AstCons).Car.Value.(string) >=
			// 	ast.Value.(AstCons).Cdr.Value.(string),
			Value: *(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car.Value)).Ptr) >=
				*(*string)((*rawInterface2)(unsafe.Pointer(&(*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr.Value)).Ptr),
		}, true, nil, nil
	}

	return *ast, false, nil, nil
}
