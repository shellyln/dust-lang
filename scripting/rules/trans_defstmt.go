package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

// AST transformation of variable/constant definition statement
func TransDefStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: var/const keyword symbol
	// [1]: multiple definitions

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	args := asts[1].Value.(AstSlice)
	// [i+0]: symbol
	// [i+1]: type (TypeInfo) / placeholder
	// [i+2]: assignment operator
	// [i+3]: value expression

	slice := make(AstSlice, len(args)/4, len(args)/4)

	isConst := false

	switch asts[0].Value.(string) {
	case "let mut":
		// nothing to do
	case "let":
		isConst = true
	case "const":
		isConst = true
	default:
		return nil, errors.New(emsg.RuleErr00015)
	}

	for i := 0; i < len(args); i += 4 {
		var opcode AstOpCodeType

		retTypeAndFlags := mnem.ReturnAny

		typeInference := false

		if args[i+1].ClassName == clsz.Placeholder {
			typeInference = true
		} else {
			typeInfo := args[i+1].Value.(tsys.TypeInfo)
			retTypeAndFlags = typeInfo.Flags
		}

		switch args[i+2].Value.(string) {
		case "=":
			if isConst {
				opcode = mnem.DefConst
			} else {
				opcode = mnem.DefVar
			}
		default:
			return nil, errors.New(emsg.RuleErr00016 + args[i+2].Value.(string))
		}

		if isConst {
			// TODO: check const expr
		}

		// TODO: BUG: typeid is not set (retTypeAndFlags)

		// TODO: BUG: Callable, Maybe, Indexable flags should be passed
		// TODO: BUG: Complex type equality checking
		op2, ok := promoteUnaryOperand(ctx, true, false, xtor.ToAstType(retTypeAndFlags), args[i+3])
		if !ok {
			return nil, errors.New(emsg.RuleErr00017 + args[i+2].Value.(string))
		}

		if typeInference {
			retTypeAndFlags = op2.OpCode &^ mnem.OpCodeMask
		}

		if isConst {
			retTypeAndFlags &^= mnem.Lvalue
		} else {
			retTypeAndFlags |= mnem.Lvalue
		}

		varName := args[i].Value.(string)

		if varName != "_" {
			_, ok = xctx.DefineVariable(args[i].Value.(string), xtor.VariableInfo{
				// TODO: If isCallable, set argument types information for compiling the call operator.
				Flags: retTypeAndFlags & (mnem.ReturnTypeMask | mnem.FlagsMask),
				Value: nil,
			})
			if !ok {
				return nil, errors.New(emsg.RuleErr00018 + args[i].Value.(string))
			}

			slice[i/4] = Ast{
				OpCode:    opcode | retTypeAndFlags,
				ClassName: clsz.C_op_def,
				Type:      AstType_AstCons,
				Value: AstCons{
					Car: Ast{
						OpCode:    mnem.Quote,
						ClassName: clsz.C_op_quote,
						Type:      AstType_AstCons,
						Value:     AstCons{Car: args[i]},
					},
					Cdr: op2,
				},
			}
		} else {
			slice[i/4] = op2
		}
	}

	if len(slice) == 1 {
		return AstSlice{slice[0]}, nil
	} else {
		return AstSlice{{
			OpCode:         mnem.Last | (slice[len(slice)-1].OpCode &^ mnem.OpCodeMask),
			ClassName:      clsz.C_op_last,
			Type:           AstType_ListOfAst,
			Value:          slice,
			SourcePosition: asts[0].SourcePosition,
		}}, nil
	}
}
