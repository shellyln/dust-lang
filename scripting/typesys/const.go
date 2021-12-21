package typesys

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
)

const (
	TypeInfo_TypeBits_Unit               uint16 = 1 << 0
	TypeInfo_TypeBits_Pointer            uint16 = 1 << 1
	TypeInfo_TypeBits_Ref                uint16 = 1 << 2
	TypeInfo_TypeBits_Function           uint16 = 1 << 3
	TypeInfo_TypeBits_NativeFn           uint16 = 1 << 4
	TypeInfo_TypeBits_Maybe              uint16 = 1 << 5
	TypeInfo_TypeBits_Array              uint16 = 1 << 6
	TypeInfo_TypeBits_Slice              uint16 = 1 << 7
	TypeInfo_TypeBits_Object             uint16 = 1 << 8
	TypeInfo_TypeBits_Struct             uint16 = 1 << 9
	TypeInfo_TypeBits_Enum               uint16 = 1 << 10 // std::option::Option use with Maybe
	TypeInfo_TypeBits_NamedItem          uint16 = 1 << 11 // use with Object,Struct,Enum
	TypeInfo_TypeBits_Tuple              uint16 = 1 << 12
	TypeInfo_TypeBits_Trait              uint16 = 1 << 13
	TypeInfo_TypeBits_GenericType        uint16 = 1 << 14
	TypeInfo_TypeBits_GenericPlaceholder uint16 = 1 << 15

	// GenericType/ -+- "std::option::Option"/ -+- "Some"/ --- GenericPlaceholder(T)
	//               |                          +- "None"/ --- unit
	//               +- GenericPlaceholder(T)

	// "std::option::Option"/ -+- "Some"/ --- s
	//                         +- "None"/ --- unit
)

var TypeInfo_Primitive_Any = TypeInfo{
	Name: "any",
	Id:   0,
}

var TypeInfo_Primitive_Unit = TypeInfo{
	Name:     "()",
	Flags:    1 << mnem.TypeIdOffset,
	Id:       1,
	TypeBits: TypeInfo_TypeBits_Unit,
}

var TypeInfo_Primitive_I64 = TypeInfo{
	Name:      "i64",
	Flags:     2<<mnem.TypeIdOffset | mnem.ReturnInt,
	Id:        2,
	Primitive: uint16(mnem.ReturnInt),
}

var TypeInfo_Primitive_I32 = TypeInfo{
	Name:      "i32",
	Flags:     3<<mnem.TypeIdOffset | mnem.ReturnInt | mnem.Bits32,
	Id:        3,
	Primitive: uint16(mnem.ReturnInt | mnem.Bits32),
}

var TypeInfo_Primitive_I16 = TypeInfo{
	Name:      "i16",
	Flags:     4<<mnem.TypeIdOffset | mnem.ReturnInt | mnem.Bits16,
	Id:        4,
	Primitive: uint16(mnem.ReturnInt | mnem.Bits16),
}

var TypeInfo_Primitive_I8 = TypeInfo{
	Name:      "i8",
	Flags:     5<<mnem.TypeIdOffset | mnem.ReturnInt | mnem.Bits8,
	Id:        5,
	Primitive: uint16(mnem.ReturnInt | mnem.Bits8),
}

var TypeInfo_Primitive_U64 = TypeInfo{
	Name:      "u64",
	Flags:     6<<mnem.TypeIdOffset | mnem.ReturnUint,
	Id:        6,
	Primitive: uint16(mnem.ReturnUint),
}

var TypeInfo_Primitive_U32 = TypeInfo{
	Name:      "u32",
	Flags:     7<<mnem.TypeIdOffset | mnem.ReturnUint | mnem.Bits32,
	Id:        7,
	Primitive: uint16(mnem.ReturnUint | mnem.Bits32),
}

var TypeInfo_Primitive_U16 = TypeInfo{
	Name:      "u16",
	Flags:     8<<mnem.TypeIdOffset | mnem.ReturnUint | mnem.Bits16,
	Id:        8,
	Primitive: uint16(mnem.ReturnUint | mnem.Bits16),
}

var TypeInfo_Primitive_U8 = TypeInfo{
	Name:      "u8",
	Flags:     9<<mnem.TypeIdOffset | mnem.ReturnUint | mnem.Bits8,
	Id:        9,
	Primitive: uint16(mnem.ReturnUint | mnem.Bits8),
}

var TypeInfo_Primitive_F64 = TypeInfo{
	Name:      "f64",
	Flags:     10<<mnem.TypeIdOffset | mnem.ReturnFloat,
	Id:        10,
	Primitive: uint16(mnem.ReturnFloat),
}

var TypeInfo_Primitive_F32 = TypeInfo{
	Name:      "f32",
	Flags:     11<<mnem.TypeIdOffset | mnem.ReturnFloat | mnem.Bits32,
	Id:        11,
	Primitive: uint16(mnem.ReturnFloat | mnem.Bits32),
}

var TypeInfo_Primitive_Bool = TypeInfo{
	Name:      "bool",
	Flags:     12<<mnem.TypeIdOffset | mnem.ReturnBool,
	Id:        12,
	Primitive: uint16(mnem.ReturnBool),
}

var TypeInfo_Primitive_String = TypeInfo{
	Name:      "String",
	Flags:     13<<mnem.TypeIdOffset | mnem.ReturnString,
	Id:        13,
	Primitive: uint16(mnem.ReturnString),
}

// It is initialized in the init() function.
var TypeInfo_Generic_Enum_Option TypeInfo
