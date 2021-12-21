package executor

import (
	"errors"
	"reflect"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

//
func execDynCallOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {
	var err error

	if ast.OpCode&mnem.OpCodeMask == mnem.DynCall {
		retType := ast.OpCode & mnem.ReturnTypeMask
		opcode, astType := getImmediateOpCode(retType)
		// args := ast.Value.(AstSlice)
		args := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)

		if (args[0].OpCode &^ mnem.FlagsMask) == mnem.Imm_data {
			if payload, ok := args[0].Value.(Ast); ok &&
				(payload.OpCode&mnem.OpCodeMask == mnem.Func ||
					payload.OpCode&mnem.OpCodeMask == mnem.Lambda) {

				var ret, thrown interface{}
				ret, thrown, err = callStandardFunc(ctx, &payload, args[1:])

				if thrown != nil {
					payload, ok := thrown.(returnPayload)
					if ok {
						retvalOpcode, retvalAstType := getImmediateOpCode(AstOpCodeType(payload.Type))
						w, ok := CastNumberLiteral(Ast{
							OpCode: retvalOpcode,
							Type:   retvalAstType,
							Value:  payload.Value,
						}, astType)

						if ok {
							return w, true, nil, nil
						} else {
							return Ast{}, true, nil, errors.New(emsg.ExecErr00003)
						}
					}
				}

				return Ast{
					OpCode: opcode,
					Type:   astType,
					Value:  ret,
				}, true, thrown, err
			}
		}
		{
			// TODO: BUG: ast.Type is wrong!
			var ret interface{}
			ret, err = callNativeFunc(args[0].Value, args[1:])

			return Ast{
				OpCode: opcode,
				Type:   astType,
				Value:  ret,
			}, true, nil, err
		}
	}

	return *ast, false, nil, nil
}

//
func callStandardFunc(ctx *ExecutionContext, fn *Ast, args AstSlice) (interface{}, interface{}, error) {
	// formalArgsAndBody := fn.Value.(AstSlice)
	formalArgsAndBody := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&fn.Value)).Ptr)
	formalArgs := formalArgsAndBody[:len(formalArgsAndBody)-1]
	body := &formalArgsAndBody[len(formalArgsAndBody)-1]

	if fn.OpCode&mnem.OpCodeMask == mnem.Func {
		ctx.PushFuncScope(fn.OpCode & mnem.ReturnTypeMask)
	} else {
		ctx.PushLambdaScope(fn.OpCode & mnem.ReturnTypeMask)
	}
	defer ctx.PopScope()

	_, ok := ctx.DefineVariable("recurse", VariableInfo{
		Flags: fn.OpCode&mnem.ReturnTypeMask | mnem.Callable,
		Value: *fn,
	})
	if !ok {
		return nil, nil, errors.New(emsg.ExecErr00004)
	}

	formalArgsLen := len(formalArgs)
	if formalArgsLen != len(args) {
		return nil, nil, errors.New(emsg.ExecErr00005)
	}

	for i := 0; i < formalArgsLen; i++ {
		// typeInfo, err := ctx.GetTypeInfoById(uint32((formalArgs[i].OpCode & mnem.TypeIdMask) >> mnem.TypeIdOffset))
		// if err != nil {
		// 	return nil,nil, err
		// }
		argType, ok, _ := tsys.ExtractReturnTypeAndFlags(ctx, false, false, false, formalArgs[i].OpCode)

		var v Ast
		if ok {
			v, ok = CastNumberLiteral(args[i], ToAstType(argType))
			if !ok {
				return nil, nil, errors.New(emsg.ExecErr00003)
			}
		} else {
			// TODO: BUG: check same type id
			v = args[i]
		}

		// _, ok = ctx.DefineVariable(formalArgs[i].Value.(string), VariableInfo{
		_, ok = ctx.DefineVariable(*(*string)((*rawInterface2)(unsafe.Pointer(&formalArgs[i].Value)).Ptr), VariableInfo{
			Flags: argType | mnem.Lvalue,
			Value: v.Value,
		})
		if !ok {
			return nil, nil, errors.New(emsg.ExecErr00004)
		}
	}

	ret, _, _, thrown, err := body.Traverse(
		traverseWayThereCallback,
		traverseWayBackCallback,
		traverseChildrenErrCallback,
		ctx,
	)
	if thrown != nil || err != nil {
		return Ast{}, thrown, err
	}

	return ret.Value, nil, nil
}

//
func callNativeFunc(fn interface{}, args AstSlice) (interface{}, error) {
	var err error

	ft := reflect.TypeOf(fn)
	fv := reflect.ValueOf(fn)

	ftNumIn := ft.NumIn()
	var va []reflect.Value
	if ft.IsVariadic() {
		va = make([]reflect.Value, len(args), len(args))
	} else {
		va = make([]reflect.Value, ftNumIn, ftNumIn)
	}

	for i := 0; i < ftNumIn; i++ {
		argtyp := ft.In(i)
		switch argtyp.Kind() {
		case reflect.Slice:
			// TODO: recursive
			if i == ftNumIn-1 && ft.IsVariadic() {
				elmtyp := argtyp.Elem()
				srcLen := len(args) - i
				for j := 0; j < srcLen; j++ {

					// TODO: recursive
					switch elmtyp.Kind() {
					case reflect.Slice:
					default:
					}

					va[i+j] = toReflectValue(args[i+j].Value).Convert(elmtyp)
				}
			} else {
				elmtyp := argtyp.Elem()
				srcSlice := reflect.ValueOf(args[i].Value)
				srcSliceLen := srcSlice.Len()

				destSlice := reflect.MakeSlice(argtyp, 0, srcSliceLen)
				for j := 0; j < srcSliceLen; j++ {

					// TODO: recursive
					switch elmtyp.Kind() {
					case reflect.Slice:
					default:
					}

					destSlice = reflect.Append(
						destSlice,
						toReflectValue(srcSlice.Index(j).Interface()).Convert(elmtyp),
					)
				}
				va[i] = destSlice
			}
		default:
			va[i] = reflect.ValueOf(args[i].Value).Convert(argtyp)
		}
	}

	var ret interface{}
	result := fv.Call(va)
	resultLen := len(result)
	err = nil

	if 0 < resultLen {
		ret = result[0].Interface()
		if 2 <= resultLen && result[resultLen-1].CanInterface() {
			err, _ = result[resultLen-1].Interface().(error)
		}
	}

	return ret, err
}

//
func toReflectValue(v interface{}) reflect.Value {
	switch w := v.(type) {
	case int64:
		return reflect.ValueOf(w)
	case int32:
		return reflect.ValueOf(w)
	case int:
		return reflect.ValueOf(w)
	case int16:
		return reflect.ValueOf(w)
	case int8:
		return reflect.ValueOf(w)
	case uint64:
		return reflect.ValueOf(w)
	case uint32:
		return reflect.ValueOf(w)
	case uint:
		return reflect.ValueOf(w)
	case uint16:
		return reflect.ValueOf(w)
	case uint8:
		return reflect.ValueOf(w)
	case float64:
		return reflect.ValueOf(w)
	case float32:
		return reflect.ValueOf(w)
	case bool:
		return reflect.ValueOf(w)
	case string:
		return reflect.ValueOf(w)
	default:
		return reflect.ValueOf(w)
	}
}

//
func getImmediateOpCode(retType AstOpCodeType) (opcode AstOpCodeType, astType AstType) {
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
	default:
		opcode = mnem.Imm_data
		astType = AstType_Any
	}
	return
}
