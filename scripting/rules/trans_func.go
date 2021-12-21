package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

//
func TransFormalArgument(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: argument name
	// [1]: type (TypeInfo)

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	typeInfo := asts[1].Value.(tsys.TypeInfo)

	_, ok = xctx.DefineVariable(asts[0].Value.(string), xtor.VariableInfo{
		Flags: typeInfo.Flags & mnem.ReturnTypeMask,
	})
	if !ok {
		return nil, errors.New(emsg.RuleErr00018 + asts[0].Value.(string))
	}

	return AstSlice{Ast{
		OpCode:    mnem.Arg | typeInfo.Flags,
		ClassName: clsz.C_op_arg,
		Type:      AstType_String,
		Value:     asts[0].Value.(string),
	}}, nil
}

// For two-pass parsing of function definitions
func TransFormalArgumentForPreScan(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: argument name
	// [1]: type (TypeInfo)

	typeInfo := asts[1].Value.(tsys.TypeInfo)

	return AstSlice{Ast{
		OpCode:    mnem.Arg | typeInfo.Flags,
		ClassName: clsz.C_op_arg,
		Type:      AstType_String,
		Value:     asts[0].Value.(string), // TODO:
	}}, nil
}

//
func TransFormalArgumentForTyping(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: argument name
	// [1]: type (TypeInfo)

	typeInfo := asts[1].Value.(tsys.TypeInfo)

	return AstSlice{Ast{
		OpCode:    typeInfo.Flags,
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     typeInfo,
	}}, nil
}

func makeFuncAst(xctx *xtor.ExecutionContext, asts AstSlice, fnOpcode AstOpCodeType) (Ast, AstSlice, AstOpCodeType, error) {
	// [0]: func name
	// [1]: Slice of Arg
	// [2]: type (TypeInfo) or placeholder -> returnType

	retTypeAndFlags := mnem.ReturnAny

	if asts[2].ClassName == clsz.Placeholder {
		// nothing to do
	} else {
		typeInfo := asts[2].Value.(tsys.TypeInfo)
		retTypeAndFlags = typeInfo.Flags
	}

	retTypeAndFlags |= mnem.Callable

	args := asts[1].Value.(AstSlice)
	argsLen := len(args)
	slice := make(AstSlice, 0, argsLen+1)
	slice = append(slice, args...)

	// TODO: BUG: build new TypeInfo for this function and number a type id
	/*retTypeInfo, err :=*/
	xctx.GetTypeInfoById(uint32((retTypeAndFlags & mnem.TypeIdMask) >> mnem.TypeIdOffset))
	typeInfo := tsys.TypeInfo{
		Name:     "",
		Flags:    0,
		Id:       0,
		TypeBits: tsys.TypeInfo_TypeBits_Function,
		Of:       []tsys.TypeInfo{
			//
		},
	}

	fn := Ast{
		OpCode:    fnOpcode | retTypeAndFlags&^mnem.TypeIdMask | AstOpCodeType(typeInfo.Id<<mnem.TypeIdOffset),
		ClassName: clsz.C_op_func,
		Type:      AstType_ListOfAst,
		Value:     slice,
	}

	return fn, slice, retTypeAndFlags, nil
}

// TODO: Make the implementation of the three functions common. (TransPreScanFuncDefinitions, TransDefineLambdaSymbols, TransDefineFuncSymbols)
// For two-pass parsing of function definitions
func TransPreScanFuncDefinitions(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: func name
	// [1]: Slice of Arg
	// [2]: symbol or placeholder -> returnType

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	fn, argsAndBody, fnTypeAndFlags, err := makeFuncAst(xctx, asts, mnem.Lambda)
	if err != nil {
		return nil, err
	}

	argsAndBody = append(argsAndBody, Ast{})
	fn.Value = argsAndBody

	symbolName := asts[0].Value.(string)

	_, ok = xctx.DefineVariable(symbolName, xtor.VariableInfo{
		Flags: fnTypeAndFlags,
		// NOTE: Set argument types information for compiling the call operator.
		Value: asts[0],
	})
	if !ok {
		// TODO: Set this if the function is declared but not defined.
		return nil, errors.New(emsg.RuleErr00018 + symbolName)
	}

	return asts, nil
}

// TODO: Make the implementation of the three functions common. (TransPreScanFuncDefinitions, TransDefineLambdaSymbols, TransDefineFuncSymbols)
func TransDefineLambdaSymbols(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: placeholder
	// [1]: Slice of Arg
	// [2]: symbol or placeholder -> returnType

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	fn, argsAndBody, fnTypeAndFlags, err := makeFuncAst(xctx, asts, mnem.Lambda)
	if err != nil {
		return nil, err
	}

	argsAndBody = append(argsAndBody, Ast{})
	fn.Value = argsAndBody

	_, ok = xctx.DefineVariable("recurse", xtor.VariableInfo{
		Flags: fnTypeAndFlags,
		// NOTE: Set argument types information for compiling the call operator.
		Value: fn,
	})
	if !ok {
		return nil, errors.New(emsg.RuleErr00018 + asts[0].Value.(string))
	}

	return asts, nil
}

// TODO: Make the implementation of the three functions common. (TransPreScanFuncDefinitions, TransDefineLambdaSymbols, TransDefineFuncSymbols)
func TransDefineFuncSymbols(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: func name
	// [1]: Slice of Arg
	// [2]: symbol or placeholder -> returnType

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	fn, argsAndBody, fnTypeAndFlags, err := makeFuncAst(xctx, asts, mnem.Func)
	if err != nil {
		return nil, err
	}

	argsAndBody = append(argsAndBody, Ast{})
	fn.Value = argsAndBody

	_, ok = xctx.DefineVariable("recurse", xtor.VariableInfo{
		Flags: fnTypeAndFlags,
		// NOTE: Set argument types information for compiling the call operator.
		Value: fn,
	})
	if !ok {
		return nil, errors.New(emsg.RuleErr00018 + asts[0].Value.(string))
	}

	_, ok = xctx.DefineVariable(asts[0].Value.(string), xtor.VariableInfo{
		Flags: fnTypeAndFlags,
		// NOTE: Set argument types information for compiling the call operator.
		Value: fn,
	})
	if !ok {
		return nil, errors.New(emsg.RuleErr00018 + asts[0].Value.(string))
	}

	return asts, nil
}

//
func TransLambdaExpression(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: placeholder
	// [1]: Slice of Arg
	// [2]: symbol or placeholder -> returnType
	// [3]: body

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	fn, argsAndBody, fnTypeAndFlags, err := makeFuncAst(xctx, asts, mnem.Lambda)
	if err != nil {
		return nil, err
	}

	retTypeAndFlags, _, err := tsys.ExtractReturnTypeAndFlags(xctx, true, false, false, fnTypeAndFlags)
	if err != nil {
		return nil, err
	}
	retTypeAndFlags &= mnem.ReturnTypeMask | mnem.TypeMarkerMask

	body, ok := promoteUnaryOperand(ctx, true, false, xtor.ToAstType(retTypeAndFlags), asts[3])
	if !ok {
		return nil, errors.New(emsg.RuleErr00003)
	}

	// TODO: traverse children return statements and transform them to add type conversion operations.
	// TODO: transform tail recursion

	argsAndBody = append(argsAndBody, body)
	fn.Value = argsAndBody

	return AstSlice{Ast{
		OpCode:    mnem.Quote | fnTypeAndFlags,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value: AstCons{Car: Ast{
			// NOTE: DON'T add return type bits. This symbol should be retrieved as `imm_data`.
			OpCode:    mnem.Imm_data | fnTypeAndFlags&mnem.FlagsMask,
			ClassName: clsz.C_imm_data,
			Type:      AstType_Any,
			Value:     fn,
		}},
	}}, nil
}

//
func TransFuncStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: func name
	// [1]: Slice of Arg
	// [2]: symbol or placeholder -> returnType
	// [3]: body

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	fn, argsAndBody, fnTypeAndFlags, err := makeFuncAst(xctx, asts, mnem.Func)
	if err != nil {
		return nil, err
	}

	retTypeAndFlags, _, err := tsys.ExtractReturnTypeAndFlags(xctx, true, false, false, fnTypeAndFlags)
	if err != nil {
		return nil, err
	}
	retTypeAndFlags &= mnem.ReturnTypeMask | mnem.TypeMarkerMask

	body, ok := promoteUnaryOperand(ctx, true, false, xtor.ToAstType(retTypeAndFlags), asts[3])
	if !ok {
		return nil, errors.New(emsg.RuleErr00003)
	}

	// TODO: traverse children return statements and transform them to add type conversion operations.
	// TODO: transform tail recursion

	argsAndBody = append(argsAndBody, body)
	fn.Value = argsAndBody

	return AstSlice{Ast{
		OpCode:    mnem.DefConst,
		ClassName: clsz.C_op_def,
		Type:      AstType_AstCons,
		Value: AstCons{
			Car: Ast{
				OpCode:    mnem.Quote,
				ClassName: clsz.C_op_quote,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: asts[0]},
			},
			Cdr: Ast{
				OpCode:    mnem.Quote | fnTypeAndFlags,
				ClassName: clsz.C_op_quote,
				Type:      AstType_AstCons,
				Value: AstCons{
					Car: Ast{
						// NOTE: DON'T add return type bits. This symbol should be retrieved as `imm_data`.
						OpCode:    mnem.Imm_data | fnTypeAndFlags&mnem.FlagsMask,
						ClassName: clsz.C_imm_data,
						Type:      AstType_Any,
						Value:     fn,
					},
				},
			},
		},
	}}, nil
}
