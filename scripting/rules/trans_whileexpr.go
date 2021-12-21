package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

//
func TransWhileExpression(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: label or placeholder
	// [1]: expression

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	args := asts[1].Value.(AstSlice)
	// [0]: variables definition
	// [1]: condition
	// [2]: body

	cond, ok := promoteUnaryOperand(ctx, true, false, AstType_Bool, args[1])
	if !ok {
		return nil, errors.New(emsg.RuleErr00004 + args[1].Value.(string))
	}

	slice := make(AstSlice, 3, 3)
	// [0]: variables definition, condition
	// [1]: body
	// [2]: label

	slice[2] = asts[0]

	if args[0].ClassName == clsz.Placeholder {
		slice[0] = Ast{
			OpCode:    mnem.Quote,
			ClassName: clsz.C_op_quote,
			Type:      AstType_AstCons,
			Value:     AstCons{Car: cond},
		}
	} else {
		slice[0] = Ast{
			OpCode:    mnem.Seq,
			ClassName: clsz.C_op_seq,
			Type:      AstType_AstCons,
			Value: AstCons{
				Car: args[0],
				Cdr: Ast{
					OpCode:    mnem.Quote,
					ClassName: clsz.C_op_quote,
					Type:      AstType_AstCons,
					Value:     AstCons{Car: cond},
				},
			},
		}
	}

	retType, _, _ := getReturnTypeAndFlags(xctx, false, false, false, args[2])

	// TODO: traverse children break statements and transform them to add type conversion operations.

	slice[1] = Ast{
		OpCode:    mnem.Quote | retType,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: args[2]},
	}

	return AstSlice{Ast{
		OpCode:    mnem.Scope | retType,
		ClassName: clsz.C_op_scope,
		Type:      AstType_AstCons,
		Value: AstCons{Car: Ast{
			OpCode:    mnem.While | retType,
			ClassName: clsz.C_op_while,
			Type:      AstType_ListOfAst,
			Value:     slice,
		}},
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}
