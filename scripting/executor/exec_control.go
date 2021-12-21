package executor

import (
	"errors"
	"reflect"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

//
func execControlOp(ctx *ExecutionContext, ast *Ast) (Ast, bool, interface{}, error) {

	switch ast.OpCode & mnem.OpCodeMask {
	case mnem.Scope:
		ctx.PopScope()
		// return ast.Value.(AstCons).Car, true, nil, nil
		return (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, true, nil, nil

	case mnem.Quote:
		// return ast.Value.(AstCons).Car, true, nil, nil
		return (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Car, true, nil, nil

	case mnem.If:
		{
			// slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			condLen := len(slice) - 1
			var w *Ast

			hit := false
			for i := 0; i < condLen; i += 2 {
				v, _, _, thrown, err := slice[i].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if thrown != nil || err != nil {
					return Ast{}, false, thrown, err
				}

				// if v.Value.(bool) {
				if *(*bool)((*rawInterface2)(unsafe.Pointer(&v.Value)).Ptr) {
					w = &slice[i+1]
					hit = true
				}
			}
			if !hit {
				w = &slice[condLen]
			}

			out, _, _, thrown, err := w.Traverse(
				traverseWayThereCallback,
				traverseWayBackCallback,
				traverseChildrenErrCallback,
				ctx,
			)
			if thrown != nil || err != nil {
				return out, false, thrown, err
			}
			return out, true, nil, nil

			// TODO: unregister labels
		}

	case mnem.While:
		{
			// slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			// [0]: variables definition, condition
			// [1]: body
			// [2]: label

			var v Ast
			var thrown interface{}
			var err error
			out := NilAst

		LOOP_WHILE:
			for {
				v, _, _, thrown, err = slice[0].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if thrown != nil || err != nil {
					return Ast{}, false, thrown, err
				}

				// if !v.Value.(bool) {
				if !*(*bool)((*rawInterface2)(unsafe.Pointer(&v.Value)).Ptr) {
					return out, true, nil, nil
				}

				out, _, _, thrown, err = slice[1].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if err != nil {
					return Ast{}, false, nil, err
				}
				if thrown != nil {
					switch payload := thrown.(type) {
					case breakPayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[2].Type == AstType_String && slice[2].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								// TODO: cast value (see also: execDynCallOp)
								out = Ast{
									OpCode: mnem.Imm_data,
									Type:   payload.Type,
									Value:  payload.Value,
								}
								break LOOP_WHILE
							}
						}
					case continuePayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[2].Type == AstType_String && slice[2].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								out = NilAst
								continue LOOP_WHILE
							}
						}
					}

					return Ast{}, true, thrown, err
				}
			}
			return out, true, nil, nil

			// TODO: unregister labels
		}

	case mnem.DoWhile:
		{
			// slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			// [0]: variables definition, condition
			// [1]: body
			// [2]: label

			var v Ast
			var thrown interface{}
			var err error
			out := NilAst

		LOOP_DOWHILE:
			for {
				out, _, _, thrown, err = slice[1].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if err != nil {
					return Ast{}, false, thrown, err
				}
				if thrown != nil {
					hit := false
					switch payload := thrown.(type) {
					case breakPayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[2].Type == AstType_String && slice[2].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								// TODO: cast value (see also: execDynCallOp)
								out = Ast{
									OpCode: mnem.Imm_data,
									Type:   payload.Type,
									Value:  payload.Value,
								}
								break LOOP_DOWHILE
							}
						}
					case continuePayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[2].Type == AstType_String && slice[2].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								hit = true
								out = NilAst
							}
						}
					}
					if !hit {
						return Ast{}, true, thrown, err
					}
				}

				v, _, _, thrown, err = slice[0].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if thrown != nil || err != nil {
					return Ast{}, false, thrown, err
				}

				//if !v.Value.(bool) {
				if !*(*bool)((*rawInterface2)(unsafe.Pointer(&v.Value)).Ptr) {
					return out, true, nil, nil
				}
			}
			return out, true, nil, nil

			// TODO: unregister labels
		}

	case mnem.ForIn, mnem.ForIni, mnem.ForInu, mnem.ForInf, mnem.ForInbool, mnem.ForInstr:
		{
			// slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			// [0]: variables definition
			// [1]: symbol
			// [2]: body
			// [3]: iterator/array
			// [4]: label

			itemType := AstType_Any
			switch ast.OpCode & mnem.OpCodeMask {
			case mnem.ForIni:
				itemType = AstType_Int
			case mnem.ForInu:
				itemType = AstType_Uint
			case mnem.ForInf:
				itemType = AstType_Float
			case mnem.ForInbool:
				itemType = AstType_Bool
			case mnem.ForInstr:
				itemType = AstType_String
			}

			// vi, ok := ctx.DefineVariable(slice[1].Value.(string), VariableInfo{
			vi, ok := ctx.DefineVariable(*(*string)((*rawInterface2)(unsafe.Pointer(&slice[1].Value)).Ptr), VariableInfo{
				Flags: ToOpCodeType(itemType),
			})
			if !ok {
				return *ast, true, nil, errors.New(emsg.ExecErr00002 + slice[1].Value.(string))
			}

			rv := reflect.ValueOf(slice[3].Value)
			rvLen := rv.Len()

			var thrown interface{}
			var err error
			out := NilAst

		LOOP_FORIN:
			for i := 0; i < rvLen; i++ {
				rvi := rv.Index(i)
				vi.Value = rvi.Interface()

				out, _, _, thrown, err = slice[2].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if err != nil {
					return Ast{}, false, thrown, err
				}
				if thrown != nil {
					hit := false
					switch payload := thrown.(type) {
					case breakPayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[4].Type == AstType_String && slice[4].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								// TODO: cast value (see also: execDynCallOp)
								out = Ast{
									OpCode: mnem.Imm_data,
									Type:   payload.Type,
									Value:  payload.Value,
								}
								break LOOP_FORIN
							}
						}
					case continuePayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[4].Type == AstType_String && slice[4].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								hit = true
								out = NilAst
							}
						}
					}
					if !hit {
						return Ast{}, true, thrown, err
					}
				}
			}
			return out, true, nil, nil

			// TODO: unregister labels
		}

	case mnem.For:
		{
			// slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			// [0]: variables definition, condition
			// [1]: body
			// [2]: inclement/declement/etc
			// [3]: label

			var v Ast
			var thrown interface{}
			var err error
			out := NilAst

		LOOP_FOR:
			for {
				v, _, _, thrown, err = slice[0].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if thrown != nil || err != nil {
					return Ast{}, false, thrown, err
				}

				//if !v.Value.(bool) {
				if !*(*bool)((*rawInterface2)(unsafe.Pointer(&v.Value)).Ptr) {
					return out, true, nil, nil
				}

				out, _, _, thrown, err = slice[1].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if err != nil {
					return Ast{}, false, thrown, err
				}
				if thrown != nil {
					hit := false
					switch payload := thrown.(type) {
					case breakPayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[3].Type == AstType_String && slice[3].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								// TODO: cast value (see also: execDynCallOp)
								out = Ast{
									OpCode: mnem.Imm_data,
									Type:   payload.Type,
									Value:  payload.Value,
								}
								break LOOP_FOR
							}
						}
					case continuePayload:
						{
							ok := true
							if payload.Label != "" {
								if slice[3].Type == AstType_String && slice[3].Value.(string) == payload.Label {
									ok = true
								} else {
									ok = false
								}
							}
							if ok {
								hit = true
								out = NilAst
							}
						}
					}
					if !hit {
						return Ast{}, true, thrown, err
					}
				}

				_, _, _, thrown, err = slice[2].Traverse(
					traverseWayThereCallback,
					traverseWayBackCallback,
					traverseChildrenErrCallback,
					ctx,
				)
				if thrown != nil || err != nil {
					return Ast{}, false, thrown, err
				}
			}
			return out, true, nil, nil

			// TODO: unregister labels
		}

	case mnem.Break:
		{
			// cons := ast.Value.(AstCons)
			cons := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			return *ast, true, breakPayload{
				// Label: cons.Car.Value.(string),
				Label: *(*string)((*rawInterface2)(unsafe.Pointer(&cons.Car.Value)).Ptr),
				Type:  ToAstType(ast.OpCode),
				Value: cons.Cdr.Value,
			}, nil
		}

	case mnem.Continue:
		{
			// cons := ast.Value.(AstCons)
			cons := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			return *ast, true, continuePayload{
				// Label: cons.Car.Value.(string),
				Label: *(*string)((*rawInterface2)(unsafe.Pointer(&cons.Car.Value)).Ptr),
			}, nil
		}

	case mnem.Ret:
		{
			// cons := ast.Value.(AstCons)
			cons := (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			return *ast, true, returnPayload{
				Type:  ToAstType(ast.OpCode),
				Value: cons.Car.Value,
			}, nil
		}

	case mnem.Seq:
		// return ast.Value.(AstCons).Cdr, true, nil, nil
		return (*AstCons)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr).Cdr, true, nil, nil

	case mnem.Last:
		{
			//slice := ast.Value.(AstSlice)
			slice := *(*AstSlice)((*rawInterface2)(unsafe.Pointer(&ast.Value)).Ptr)
			if 0 < len(slice) {
				return slice[len(slice)-1], true, nil, nil
			} else {
				return Ast{
					OpCode: mnem.Imm_nil,
				}, true, nil, nil
			}
		}
	}

	return *ast, false, nil, nil
}
