package zeros

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

//
type Unit struct{}

//
var UnitSingleton = Unit{}

//
var NilAst = Ast{
	OpCode:    mnem.Imm_nil,
	ClassName: clsz.C_imm_nil,
	Type:      AstType_Nil,
	Value:     nil,
}

//
var UnitAst = Ast{
	OpCode:    mnem.Imm_unitval | AstOpCodeType(tsys.TypeInfo_Primitive_Unit.Id<<mnem.TypeIdOffset),
	ClassName: clsz.C_imm_ndef,
	Type:      AstType_Any,
	Value:     &UnitSingleton,
}

//
var TrueAst = Ast{
	OpCode:    mnem.Imm_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id<<mnem.TypeIdOffset),
	ClassName: clsz.C_imm_bool,
	Type:      AstType_Bool,
	Value:     true,
}

//
var FalseAst = Ast{
	OpCode:    mnem.Imm_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id<<mnem.TypeIdOffset),
	ClassName: clsz.C_imm_bool,
	Type:      AstType_Bool,
	Value:     false,
}

//
var ZeroStrAst = Ast{
	OpCode:    mnem.Imm_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id<<mnem.TypeIdOffset),
	ClassName: clsz.C_imm_str,
	Type:      AstType_String,
	Value:     "",
}
