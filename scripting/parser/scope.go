package parser

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

//
func scope(parser ParserFn) ParserFn {
	fn := func(ctx ParserContext) (ParserContext, error) {
		xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
		if !ok {
			return ctx, errors.New(emsg.InternalErr00001)
		}

		xctx.PushScope(mnem.ReturnAny)
		defer xctx.PopScope()

		return parser(ctx)
	}
	return fn
}

//
func funcScope(parser ParserFn) ParserFn {
	fn := func(ctx ParserContext) (ParserContext, error) {
		xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
		if !ok {
			return ctx, errors.New(emsg.InternalErr00001)
		}

		xctx.PushFuncScope(mnem.ReturnAny)
		defer xctx.PopScope()

		return parser(ctx)
	}
	return fn
}

//
func lambdaScope(parser ParserFn) ParserFn {
	fn := func(ctx ParserContext) (ParserContext, error) {
		xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
		if !ok {
			return ctx, errors.New(emsg.InternalErr00001)
		}

		xctx.PushLambdaScope(mnem.ReturnAny)
		defer xctx.PopScope()

		return parser(ctx)
	}
	return fn
}
