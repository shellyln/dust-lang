package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

//
func TransReturnStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	var car Ast

	if asts[0].ClassName == clsz.Placeholder {
		car = NilAst
	} else {
		car = asts[0]
		// TODO: check type
		// TODO: cast
	}

	cdrType, _, err := getReturnTypeAndFlags(xctx, false, false, false, car)
	if err != nil {
		return nil, errors.New(emsg.RuleErr00020)
	}

	// TODO: BUG: type casting

	return AstSlice{{
		OpCode:    mnem.Ret | cdrType,
		ClassName: clsz.C_op_ret,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: car},
	}}, nil
}

//
func TransBreakStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	var car, cdr Ast

	if asts[0].ClassName == clsz.Placeholder {
		car = ZeroStrAst
	} else {
		car = Ast{
			OpCode:    mnem.Imm_str,
			ClassName: clsz.C_imm_str,
			Type:      AstType_String,
			Value:     asts[0].Value.(string),
		}
		// TODO: check label name
	}

	if asts[1].ClassName == clsz.Placeholder {
		cdr = NilAst
	} else {
		cdr = asts[1]
		// TODO: check type
		// TODO: cast
	}

	cdrType, _, err := getReturnTypeAndFlags(xctx, false, false, false, cdr)
	if err != nil {
		return nil, errors.New(emsg.RuleErr00020)
	}

	// TODO: BUG: type casting

	return AstSlice{{
		OpCode:    mnem.Break | cdrType,
		ClassName: clsz.C_op_break,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: car, Cdr: cdr},
	}}, nil
}

//
func TransContinueStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	var car Ast

	if asts[0].ClassName == clsz.Placeholder {
		car = ZeroStrAst
	} else {
		car = Ast{
			OpCode:    mnem.Imm_str,
			ClassName: clsz.C_imm_str,
			Type:      AstType_String,
			Value:     asts[0].Value.(string),
		}
		// TODO: check label name
	}

	return AstSlice{{
		OpCode:    mnem.Continue,
		ClassName: clsz.C_op_continue,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: car},
	}}, nil
}
