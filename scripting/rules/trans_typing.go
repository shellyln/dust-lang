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
func TransTypeObjectExpr(ctx ParserContext, asts AstSlice) (AstSlice, error) { // TODO:

	typeInfo := tsys.TypeInfo{
		Name:      "",
		ElName:    "",
		Flags:     0, // masked by MetaInfoMask (include type id)
		Id:        0,
		TypeBits:  0,
		Primitive: 0,                 // AstOpCodeType compatible rune/int/uint/float/bool/string and bits. masked by (ReturnTypeMask|BitLenMask)
		LenIdx:    0,                 // array length/arg|tuple|enum item index
		Of:        []tsys.TypeInfo{}, // ret[0]+arguments[1:]/members/enum/generic type[0]+params[1:]
	}

	// TODO: BUG: search existing type and number the type id

	// TODO:
	return AstSlice{Ast{
		OpCode:    tsys.TypeInfo_Primitive_Any.Flags, // TODO: BUG:
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     typeInfo, // TODO: BUG:
	}}, nil
}

//
func TransTypeFuncExpr(ctx ParserContext, asts AstSlice) (AstSlice, error) { // TODO:

	// [0]: placeholder
	// [1]: arguments (AstSlice TypeInfo)
	// [2]: return type (TypeInfo)

	argsAsts := asts[1].Value.(AstSlice)
	of := make([]tsys.TypeInfo, len(argsAsts)+1, len(argsAsts)+1)
	of[0] = asts[2].Value.(tsys.TypeInfo)
	for i := 0; i < len(argsAsts); i++ {
		of[1+i] = argsAsts[i].Value.(tsys.TypeInfo)
	}

	flags := of[0].Flags
	flags &^= mnem.TypeIdMask
	if flags&mnem.Callable != 0 {
		flags &^= (mnem.TypeMarkerMask | mnem.ReturnTypeMask | mnem.BitLenMask)
	}
	flags |= mnem.Callable

	// TODO: BUG: search existing type and number the type id

	typeInfo := tsys.TypeInfo{
		Name:      "",
		ElName:    "",
		Flags:     flags, // masked by MetaInfoMask (include type id) // TODO: BUG: type id
		Id:        0,     // TODO: BUG: type id
		TypeBits:  tsys.TypeInfo_TypeBits_Function,
		Primitive: 0,  // AstOpCodeType compatible rune/int/uint/float/bool/string and bits. masked by (ReturnTypeMask|BitLenMask)
		LenIdx:    0,  // array length/arg|tuple|enum item index
		Of:        of, // ret[0]+arguments[1:]/members/enum/generic type[0]+params[1:]
	}

	// TODO:
	return AstSlice{Ast{
		OpCode:    tsys.TypeInfo_Primitive_Any.Flags, // TODO: BUG:
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     typeInfo,
	}}, nil
}

//
func TransTypeSliceExpr(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	// [0]: TypeInfo

	elTypeInfo := asts[0].Value.(tsys.TypeInfo)
	sliceTypeInfo, err := xctx.GetRegisteredTypeInfo(elTypeInfo.Slice(xctx.DummyTypeId))
	if err != nil {
		return nil, err
	}

	return AstSlice{Ast{
		OpCode:    sliceTypeInfo.Flags,
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     sliceTypeInfo,
	}}, nil
}

//
func TransTypeGenericsExpr(ctx ParserContext, asts AstSlice) (AstSlice, error) { // TODO:
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	// [0]    : symbol
	// [1]... : TypeInfo

	// TODO: find the custom type by symbol name ([0])
	//       (from scope namespace)

	_ /*typeInfo*/, err := xctx.GetTypeInfoByName(asts[0].Value.(string))
	if err != nil {
		// return nil, err // TODO:
	}

	// TODO: BUG: search existing type and number the type id

	// TODO:
	return AstSlice{Ast{
		OpCode:    tsys.TypeInfo_Primitive_Any.Flags,
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     tsys.TypeInfo_Primitive_Any,
	}}, nil
}

//
func TransTypeExpr(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	// [0]: symbol

	typeInfo, err := xctx.GetTypeInfoByName(asts[0].Value.(string))
	if err != nil {
		return nil, err
	}

	return AstSlice{Ast{
		OpCode:    typeInfo.Flags,
		ClassName: clsz.Type,
		Type:      AstType_Any,
		Value:     typeInfo,
	}}, nil
}
