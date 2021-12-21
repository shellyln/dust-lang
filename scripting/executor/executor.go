package executor

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func traverseWayThereCallback(ctxif interface{}, ast Ast) (Ast, WayThereMode, error) {
	switch ast.OpCode & mnem.OpCodeMask {
	case mnem.Quote:
		return ast, WayThereMode_Lazy, nil
	case mnem.Last:
		return ast, WayThereMode_Last, nil
	case mnem.Scope:
		// ctxif.(*ExecutionContext).PushScope(ast.OpCode & mnem.ReturnTypeMask)
		(*ExecutionContext)((*rawInterface2)(unsafe.Pointer(&ctxif)).Ptr).PushScope(ast.OpCode & mnem.ReturnTypeMask)
		return ast, WayThereMode_None, nil

	// TODO: register labels

	default:
		return ast, WayThereMode_None, nil
	}
}

//
func traverseWayBackCallback(ctxif interface{}, ast Ast) (Ast, int16, TraverseOpcode, interface{}, error) {
	// ctx := ctxif.(*ExecutionContext)
	ctx := (*ExecutionContext)((*rawInterface2)(unsafe.Pointer(&ctxif)).Ptr)

	opcode := ast.OpCode & mnem.OpCodeMask
	opcodeAndType := ast.OpCode &^ mnem.FlagsMask

	if opcode == mnem.Symbol || (mnem.Imm_begin < opcodeAndType && opcodeAndType < mnem.Imm_end) {
		ret, ok, thrown, err := execValueOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Conv_begin < opcodeAndType && opcodeAndType < mnem.Conv_end {
		ret, ok, thrown, err := execConvertOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Arithmetic_begin < opcodeAndType && opcodeAndType < mnem.Arithmetic_end {
		ret, ok, thrown, err := execArithmeticOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.BitwiseAndLogical_begin < opcodeAndType && opcodeAndType < mnem.BitwiseAndLogical_end {
		ret, ok, thrown, err := execBitwiseAndLogicalOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Assign_begin < opcodeAndType && opcodeAndType < mnem.Assign_end {
		ret, ok, thrown, err := execAssignOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Control_begin < opcodeAndType && opcodeAndType < mnem.Control_end {
		ret, ok, thrown, err := execControlOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Call_begin < opcodeAndType && opcodeAndType < mnem.Call_end {
		ret, ok, thrown, err := execDynCallOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if mnem.Object_begin < opcodeAndType && opcodeAndType < mnem.Object_end {
		ret, ok, thrown, err := execObjectOp(ctx, &ast)
		if ok || thrown != nil || err != nil {
			return ret, 0, TraverseOpcode_Nop, thrown, err
		}
		goto END
	}

	if opcode == mnem.Nop {
		return ast, 0, TraverseOpcode_Nop, nil, nil
	}

END:
	return ast,
		0, TraverseOpcode_Nop, nil,
		errors.New(emsg.ExecErr00011 + strconv.FormatInt(int64(ast.OpCode), 10))
}

//
func traverseChildrenErrCallback(ctxif interface{}, ast Ast, child Ast, thrown interface{}, err error) {
	switch ast.OpCode & mnem.OpCodeMask {
	case mnem.Scope:
		// ctxif.(*ExecutionContext).PopScope()
		(*ExecutionContext)((*rawInterface2)(unsafe.Pointer(&ctxif)).Ptr).PopScope()

		// TODO: unregister labels
	}
}

//
func Execute(ctx *ExecutionContext, ast Ast) (interface{}, error) {
	ast, _, _, thrown, err := ast.Traverse(
		traverseWayThereCallback,
		traverseWayBackCallback,
		traverseChildrenErrCallback,
		ctx,
	)
	if thrown != nil {
		return nil, errors.New(emsg.ExecErr00012 + fmt.Sprintf("%v", thrown))
	}
	return ast.Value, err
}
