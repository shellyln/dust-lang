package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

// Return [0] full-width type flag that includes StorageMarker and TypeMarker.
// Return [1] `true` if type is primitive int|uint|bool|string.
func getReturnTypeAndFlags(
	ctx *xtor.ExecutionContext,
	getCallResultType bool, getMaybeResultType bool, getIndexResultType bool, op Ast) (AstOpCodeType, bool, error) {

	if op.OpCode&mnem.OpCodeMask == mnem.Symbol {
		vi, ok := ctx.GetVariableInfo(op.Value.(string))
		if !ok {
			return mnem.ReturnAny, false, errors.New(emsg.ExecErr00013 + op.Value.(string))
		}
		return tsys.ExtractReturnTypeAndFlags(ctx, getCallResultType, getMaybeResultType, getIndexResultType, vi.Flags)
	}

	if (op.OpCode & mnem.ReturnTypeMask) != mnem.ReturnAny {
		return tsys.ExtractReturnTypeAndFlags(ctx, getCallResultType, getMaybeResultType, getIndexResultType, op.OpCode)
	}

	return op.OpCode &^ mnem.OpCodeMask, false, nil
}

// TODO: from-to
// TODO: use type id or full AstOpCodeType
// NOTE: this function is only used at rules20 (`as` operator)
func getConversionOpCode(from AstOpCodeType) (opcode AstOpCodeType, className string) {
	switch from & mnem.ReturnTypeMask {
	case mnem.ReturnInt:
		opcode = mnem.Convi
		className = "$convi"
	case mnem.ReturnUint:
		opcode = mnem.Convu
		className = "$convu"
	case mnem.ReturnFloat:
		opcode = mnem.Convf
		className = "$convf"
	case mnem.ReturnBool:
		opcode = mnem.Convbool
		className = "$convbool"
	case mnem.ReturnString:
		opcode = mnem.Convstr
		className = "$convstr"
	default:
		opcode = mnem.Convdyn
		className = "$convdyn"
	}
	return
}
