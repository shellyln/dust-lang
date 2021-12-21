package typesys

import (
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

// TODO:
type TypeInfo struct {
	Name      string        `json:"-"`
	ElName    string        `json:"n,omitempty"`
	Flags     AstOpCodeType `json:"-"` // masked by MetaInfoMask (include type id)
	Id        uint32        `json:"-"` // type id
	TypeBits  uint16        `json:"t,omitempty"`
	Primitive uint16        `json:"p,omitempty"` // AstOpCodeType compatible rune/int/uint/float/bool/string and bits. masked by (ReturnTypeMask|BitLenMask)
	LenIdx    uint          `json:"l,omitempty"` // array length/arg|tuple|enum item index
	Of        []TypeInfo    `json:"o,omitempty"` // ret[0]+arguments[1:]/members/enum/generic type[0]+params[1:]
}

//
func (s TypeInfo) Unwrap() TypeInfo {
	// TODO: BUG: case of `Maybe`: current -> [0] Some -> [0] T

	return s.Of[0]
}

//
func (s TypeInfo) Slice(genId func() uint32) TypeInfo {
	flags := s.Flags

	// Meaning of `Lvalue|Callable|Maybe|Indexable|Bits8|ReturnUint` is:
	// The address of the variable that can be assigned to and
	// It is function that returns
	//   Option<?> of
	//   Array<?> of
	//   Uint8

	if flags&mnem.Callable != 0 || flags&mnem.Maybe != 0 || flags&mnem.Indexable != 0 {
		flags &^= (mnem.ReturnTypeMask | mnem.BitLenMask)
	}
	flags |= mnem.Indexable

	id := genId()

	return TypeInfo{
		Name:     "[" + s.Name + "]",
		Flags:    (flags &^ mnem.TypeIdMask) | AstOpCodeType(id),
		Id:       id,
		TypeBits: TypeInfo_TypeBits_Slice,
		Of:       []TypeInfo{s},
	}
}

//
func (s TypeInfo) Option(genId func() uint32) TypeInfo {
	return TypeInfo_Generic_Enum_Option.Specialize(genId, []TypeInfo{s})
}

//
func (s TypeInfo) Specialize(genId func() uint32, actualParams []TypeInfo) TypeInfo {
	z, _ := traverseTypeInfoForSpecialize(s.Of[0], genId, s.Of[1:], actualParams)
	return z
}
