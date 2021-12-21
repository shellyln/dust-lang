package typesys

import (
	"strings"

	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
)

//
func traverseTypeInfoForSpecialize(
	x TypeInfo, genId func() uint32, formalParams []TypeInfo, actualParams []TypeInfo) (TypeInfo, bool) {

	replaced := false
	childrenReplaced := false
	for i := 0; i < len(formalParams); i++ {
		if x.Id == formalParams[i].Id {
			x = actualParams[i]
			replaced = true
			break
		}
	}

	if x.Of != nil {
		slice := make([]TypeInfo, len(x.Of), len(x.Of))
		for i := 0; i < len(x.Of); i++ {
			var rep bool
			slice[i], rep = traverseTypeInfoForSpecialize(x.Of[i], genId, formalParams, actualParams)
			if rep {
				childrenReplaced = true
				c := slice[i]

				for j := 0; j < len(formalParams); j++ {
					x.Name = strings.ReplaceAll(x.Name, "<"+formalParams[j].Name+">", "<"+actualParams[j].Name+">")
				}

				// TODO: BUG: replace Id (x.Id and x.Flags)

				x.Flags &^= (mnem.ReturnTypeMask | mnem.BitLenMask)

				if x.TypeBits&TypeInfo_TypeBits_Function != 0 ||
					x.TypeBits&TypeInfo_TypeBits_NativeFn != 0 {

					if c.Flags&mnem.Callable == 0 {
						x.Flags |= c.Flags & (mnem.ReturnTypeMask | mnem.BitLenMask)
					}
					x.Flags |= mnem.Callable

				} else if x.TypeBits&TypeInfo_TypeBits_Maybe != 0 {
					if c.Flags&(mnem.Callable|c.Flags&mnem.Maybe) == 0 {
						x.Flags |= c.Flags & (mnem.ReturnTypeMask | mnem.BitLenMask)
					}
					x.Flags |= mnem.Maybe

				} else if x.TypeBits&TypeInfo_TypeBits_Array != 0 ||
					x.TypeBits&TypeInfo_TypeBits_Slice != 0 {

					if c.Flags&(mnem.Callable|c.Flags&mnem.Maybe|c.Flags&mnem.Indexable) == 0 {
						x.Flags |= c.Flags & (mnem.ReturnTypeMask | mnem.BitLenMask)
					}
					x.Flags |= mnem.Indexable

				} else if x.TypeBits&TypeInfo_TypeBits_Enum != 0 {
					// enum and enum item
					// nothing to do

				} else if x.TypeBits&TypeInfo_TypeBits_Object != 0 ||
					x.TypeBits&TypeInfo_TypeBits_Struct != 0 {

					x.Flags |= mnem.Indexable
				}
			}
		}
		x.Of = slice
	}

	return x, replaced || childrenReplaced
}
