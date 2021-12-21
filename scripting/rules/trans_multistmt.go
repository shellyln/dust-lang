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
func TransMultipleStatement(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	args := asts[0].Value.(AstSlice)
	if len(args) == 0 {
		return AstSlice{Ast{
			OpCode: mnem.Imm_nil,
		}}, nil
	} else {
		slice := make(AstSlice, 0, len(args))
		slice = append(slice, args...)

		retType, _, _ := getReturnTypeAndFlags(xctx, false, false, false, slice[len(slice)-1])
		retType &^= mnem.StorageMarkerMask

		if len(args) == 1 {
			c := args[0].OpCode & mnem.OpCodeMask
			if mnem.Imm_begin < c && c < mnem.Imm_end {
				return AstSlice{args[0]}, nil
			} else {
				return AstSlice{{
					OpCode:         mnem.Scope | retType,
					ClassName:      clsz.C_op_scope,
					Type:           AstType_AstCons,
					Value:          AstCons{Car: args[0]},
					SourcePosition: asts[0].SourcePosition,
				}}, nil
			}
		} else {
			return AstSlice{{
				OpCode:    mnem.Scope | retType,
				ClassName: clsz.C_op_scope,
				Type:      AstType_AstCons,
				Value: AstCons{Car: Ast{
					OpCode:         mnem.Last | retType,
					ClassName:      clsz.C_op_last,
					Type:           AstType_ListOfAst,
					Value:          slice,
					SourcePosition: asts[0].SourcePosition,
				}},
				SourcePosition: asts[0].SourcePosition,
			}}, nil
		}
	}
}
