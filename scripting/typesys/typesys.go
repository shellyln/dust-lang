package typesys

import (
	"encoding/json"
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

// Function, NativeFn

var TypeInfoInitialMap map[string]TypeInfo
var TypeIdInitialMax uint32

func init() {
	TypeInfoInitialMap = make(map[string]TypeInfo)

	TypeInfoPrimitives := []TypeInfo{
		TypeInfo_Primitive_Any,
		TypeInfo_Primitive_Unit,
		TypeInfo_Primitive_I64,
		TypeInfo_Primitive_I32,
		TypeInfo_Primitive_I16,
		TypeInfo_Primitive_I8,
		TypeInfo_Primitive_U64,
		TypeInfo_Primitive_U32,
		TypeInfo_Primitive_U16,
		TypeInfo_Primitive_U8,
		TypeInfo_Primitive_F64,
		TypeInfo_Primitive_F32,
		TypeInfo_Primitive_Bool,
		TypeInfo_Primitive_String,
	}

	TypeIdInitialMax = uint32(len(TypeInfoPrimitives))

	genId := func() uint32 {
		TypeIdInitialMax += 1
		return TypeIdInitialMax
	}

	// Initialize primitive types
	for _, item := range TypeInfoPrimitives {
		if bytes, err := json.Marshal(item); err == nil {
			TypeInfoInitialMap[string(bytes)] = item
		} else {
			panic(err)
		}
	}

	// Initialize slice types
	for _, item := range TypeInfoPrimitives {
		z := item.Slice(genId)
		if bytes, err := json.Marshal(z); err == nil {
			z.Flags = item.Flags&^mnem.TypeIdMask | mnem.Indexable | AstOpCodeType(z.Id<<uint32(mnem.TypeIdOffset))
			TypeInfoInitialMap[string(bytes)] = z
		} else {
			panic(err)
		}
	}

	// Initialize `Option<T>` type
	{
		idGeneric := genId()
		idEnum := genId()
		idSome := genId()
		idT := genId()
		idNone := genId()

		t := TypeInfo{
			Name:     "T",
			Flags:    AstOpCodeType(idT << mnem.TypeIdOffset),
			Id:       idT,
			TypeBits: TypeInfo_TypeBits_GenericPlaceholder,
		}
		name := "std::option::Option<" + t.Name + ">"

		TypeInfo_Generic_Enum_Option = TypeInfo{
			Name:     name,
			Flags:    mnem.Maybe | AstOpCodeType(idGeneric<<mnem.TypeIdOffset),
			Id:       idGeneric,
			TypeBits: TypeInfo_TypeBits_GenericType,
			Of: []TypeInfo{{
				Name:     name,
				Flags:    mnem.Maybe | AstOpCodeType(idEnum<<mnem.TypeIdOffset),
				Id:       idEnum,
				TypeBits: TypeInfo_TypeBits_Enum | TypeInfo_TypeBits_Maybe,
				Of: []TypeInfo{{
					ElName:   "Some",
					Flags:    t.Flags&^mnem.TypeIdMask | AstOpCodeType(idSome<<mnem.TypeIdOffset),
					Id:       idSome,
					TypeBits: TypeInfo_TypeBits_Enum | TypeInfo_TypeBits_NamedItem,
					Of:       []TypeInfo{t},
				}, {
					ElName:   "None",
					Flags:    TypeInfo_Primitive_Unit.Flags&^mnem.TypeIdMask | AstOpCodeType(idNone<<mnem.TypeIdOffset),
					Id:       idNone,
					TypeBits: TypeInfo_TypeBits_Enum | TypeInfo_TypeBits_NamedItem,
					Of:       []TypeInfo{TypeInfo_Primitive_Unit},
				}},
			}, t},
		}
		if bytes, err := json.Marshal(TypeInfo_Generic_Enum_Option); err == nil {
			TypeInfoInitialMap[string(bytes)] = TypeInfo_Generic_Enum_Option
		} else {
			panic(err)
		}
	}

	// Initialize option types
	for _, item := range TypeInfoPrimitives {
		z := item.Option(genId)
		if bytes, err := json.Marshal(z); err == nil {
			z.Flags = item.Flags&^mnem.TypeIdMask | mnem.Maybe | AstOpCodeType(z.Id<<uint32(mnem.TypeIdOffset))
			TypeInfoInitialMap[string(bytes)] = z
		} else {
			panic(err)
		}
	}
}

//
func GetTypeInfo(name string) (TypeInfo, error) {
	switch name {
	case "any":
		return TypeInfo_Primitive_Any, nil
	case "int", "isize", "i64":
		return TypeInfo_Primitive_I64, nil
	case "i32", "rune":
		return TypeInfo_Primitive_I32, nil
	case "i16":
		return TypeInfo_Primitive_I16, nil
	case "i8":
		return TypeInfo_Primitive_I8, nil
	case "uint", "usize", "u64":
		return TypeInfo_Primitive_U64, nil
	case "u32", "char":
		return TypeInfo_Primitive_U32, nil
	case "u16":
		return TypeInfo_Primitive_U16, nil
	case "u8":
		return TypeInfo_Primitive_U8, nil
	case "float", "f64":
		return TypeInfo_Primitive_F64, nil
	case "f32":
		return TypeInfo_Primitive_F32, nil
	case "bool":
		return TypeInfo_Primitive_Bool, nil
	case "string", "String":
		return TypeInfo_Primitive_String, nil
	default:
		return TypeInfo_Primitive_Any, errors.New(emsg.ExecErr00016)
	}
}
