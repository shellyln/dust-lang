package executor

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

// TODO: BUG: `xtor.ToAstType` returns w/o checking type marker flags
func ToAstType(ty AstOpCodeType) AstType {
	// NOTE: AstOpCodeType(0): any
	//       AstType(0)      : nil
	//       AstType_Any     : any

	if ty&mnem.TypeMarkerMask != 0 {
		return AstType_Any
	}

	ty &= mnem.ReturnTypeMask
	if ty == mnem.ReturnAny {
		return AstType_Any
	} else {
		return AstType(ty)
	}
}

//
func ToOpCodeType(ty AstType) AstOpCodeType {
	// NOTE: AstOpCodeType(0): any
	//       AstType(0)      : nil
	//       AstType_Any     : any
	if ty == AstType_Any {
		return mnem.ReturnAny
	} else {
		return AstOpCodeType(ty)
	}
}
