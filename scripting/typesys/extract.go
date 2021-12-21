package typesys

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func ExtractReturnTypeAndFlags(
	ctx TypeManagerContext,
	getCallResultType bool, getMaybeResultType bool, getIndexResultType bool, flags AstOpCodeType) (AstOpCodeType, bool, error) {

	// Meaning of `Lvalue|Callable|Maybe|Indexable|Bits8|ReturnUint` is:
	// The address of the variable that can be assigned to and
	// It is function that returns
	//   Option<?> of
	//   Array<?> of
	//   Uint8

	typeInfo, err := ctx.GetTypeInfoById(uint32((flags & mnem.TypeIdMask) >> mnem.TypeIdOffset))
	if err != nil {
		return flags &^ mnem.OpCodeMask, false, err
	}

	if getCallResultType {
		if flags&mnem.Callable == 0 {
			// It is NOT callable
			return mnem.ReturnAny, false, errors.New(emsg.ExecErr00014)
		}
		flags &^= mnem.Callable

		// TODO: BUG: Find and set new type id
		// typeInfo.Unwrap() // BUG: Symbol `f` (type of `let f = |a: f64, b: f64| -> f64 {a + b};`) is NOT f64(typeid=10)
	} else {
		if flags&mnem.Callable != 0 {
			// It is callable
			return flags &^ mnem.OpCodeMask, false, nil
		}
	}

	// No need to check the "Callable" flag. It has already been considered.
	if getMaybeResultType {
		if flags&mnem.Maybe == 0 {
			// It is NOT maybe
			return mnem.ReturnAny, false, errors.New(emsg.ExecErr00015)
		}
		flags &^= mnem.Maybe

		// TODO: BUG: Find and set new type id
		// typeInfo.Unwrap()
	} else {
		if flags&mnem.Maybe != 0 {
			// It is maybe
			return flags &^ mnem.OpCodeMask, false, nil
		}
	}

	// No need to check the "Callable" and "Maybe" flags. They have already been considered.
	if getIndexResultType {
		if flags&mnem.Indexable == 0 && flags&mnem.ReturnTypeMask == mnem.ReturnString {
			// `"str"[x]` is `u8`
			return TypeInfo_Primitive_U8.Flags, true, nil
		}
		if flags&mnem.Indexable == 0 {
			// It is NOT indexable
			return mnem.ReturnAny, false, errors.New(emsg.ExecErr00021)
		}
		flags &^= mnem.Indexable

		typeInfo.Unwrap()
	} else {
		if flags&(mnem.Callable|mnem.Maybe) == 0 && flags&mnem.Indexable != 0 {
			// It is indexable
			return flags &^ mnem.OpCodeMask, false, nil
		}
	}

	if mnem.ReturnInt <= flags&mnem.ReturnTypeMask && flags&mnem.ReturnTypeMask <= mnem.ReturnString {
		return flags &^ mnem.OpCodeMask, true, nil
	}

	return flags &^ mnem.OpCodeMask, false, nil
}
