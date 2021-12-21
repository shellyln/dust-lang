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
func TransDefForInSymbol(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: variables definition or placeholder
	// [1]: symbol
	// [2]: `Range` className or placeholder
	// [3]: range start or iterator
	// [4]: range end or placeholder
	// [5]: body

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	var defType AstOpCodeType

	if asts[2].ClassName == clsz.Range {
		ops, ok := promoteBinaryOperands(ctx, false, true, AstType_Float, asts[3], asts[4])
		if !ok {
			// TODO: BUG: case of [3] is placeholder
			return nil, errors.New(emsg.RuleErr00005 + "for")
		}

		defType, _, _ = getReturnTypeAndFlags(xctx, false, false, false, ops[0])
	} else {
		defType, _, _ = getReturnTypeAndFlags(xctx, false, false, true, asts[3])
	}

	_, ok = xctx.DefineVariable(asts[1].Value.(string), xtor.VariableInfo{
		Flags: defType,
	})
	if !ok {
		return nil, errors.New(emsg.RuleErr00018 + asts[1].Value.(string))
	}

	return asts, nil
}

//
// TODO: scoped def expression
func TransForInExpression(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: label or placeholder
	// [1]: expression

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	args := asts[1].Value.(AstSlice)
	// [0]: variables definition or placeholder
	// [1]: symbol
	// [2]: `Range` className or placeholder
	// [3]: range start or iterator
	// [4]: range end or placeholder
	// [5]: body

	if args[2].ClassName == clsz.Range {
		ops, ok := promoteBinaryOperands(ctx, false, true, AstType_Float, args[3], args[4])
		if !ok {
			return nil, errors.New(emsg.RuleErr00005 + "for")
		}

		defType, _, _ := getReturnTypeAndFlags(xctx, false, false, false, ops[0])
		opDef := Ast{
			OpCode:    mnem.DefVar | defType,
			ClassName: clsz.C_op_def,
			Type:      AstType_AstCons,
			Value: AstCons{
				Car: Ast{
					OpCode:    mnem.Quote,
					ClassName: clsz.C_op_quote,
					Type:      AstType_AstCons,
					Value:     AstCons{Car: args[1]},
				},
				Cdr: ops[0],
			},
		}

		if args[0].ClassName != clsz.Placeholder {
			opDef = Ast{
				OpCode:    mnem.Seq | defType,
				ClassName: clsz.C_op_seq,
				Type:      AstType_AstCons,
				Value: AstCons{
					Car: args[0],
					Cdr: opDef,
				},
			}
		}

		var condOpcode AstOpCodeType
		switch defType & mnem.ReturnTypeMask {
		case mnem.ReturnInt:
			condOpcode = mnem.CmpLEi_bool
		case mnem.ReturnUint:
			condOpcode = mnem.CmpLEu_bool
		case mnem.ReturnFloat:
			condOpcode = mnem.CmpLEf_bool
		}

		opCond := Ast{
			OpCode:         condOpcode,
			Type:           AstType_AstCons,
			Value:          AstCons{Car: args[1], Cdr: ops[1]},
			SourcePosition: asts[0].SourcePosition,
		}

		opIncl := Ast{
			OpCode: mnem.PreIncr | mnem.Lvalue | defType,
			Type:   AstType_AstCons,
			Value:  AstCons{Car: args[1]},
		}

		return TransTradForExpression(ctx, AstSlice{
			asts[0],
			Ast{
				ClassName: "Group",
				Type:      AstType_ListOfAst,
				Value: AstSlice{
					opDef, opCond, opIncl, args[5],
				},
				SourcePosition: asts[0].SourcePosition,
			},
		})
	}

	slice := make(AstSlice, 5, 5)
	// [0]: variables definition or placeholder
	// [1]: symbol
	// [2]: body
	// [3]: iterator/array
	// [4]: label

	slice[0] = args[0]
	slice[1] = Ast{
		OpCode:    mnem.Quote,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: args[1]},
	}

	retType, _, _ := getReturnTypeAndFlags(xctx, false, false, false, args[5])

	slice[2] = Ast{
		OpCode:    mnem.Quote | retType,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: args[5]},
	}

	// TODO: traverse children break statements and transform them to add type conversion operations.

	itemType, _, _ := getReturnTypeAndFlags(xctx, false, false, true, args[3])
	itemType = xtor.ToOpCodeType(AstType_Any)

	opcode := mnem.ForIn
	switch itemType & mnem.ReturnTypeMask {
	case mnem.ReturnInt:
		opcode = mnem.ForIni
	case mnem.ReturnUint:
		opcode = mnem.ForInu
	case mnem.ReturnFloat:
		opcode = mnem.ForInf
	case mnem.ReturnBool:
		opcode = mnem.ForInbool
	case mnem.ReturnString:
		opcode = mnem.ForInstr
	}

	slice[3] = args[3]
	slice[4] = asts[0]

	return AstSlice{Ast{
		OpCode:    mnem.Scope | retType,
		ClassName: clsz.C_op_scope,
		Type:      AstType_AstCons,
		Value: AstCons{Car: Ast{
			OpCode:    opcode | retType,
			ClassName: clsz.C_op_forin,
			Type:      AstType_ListOfAst,
			Value:     slice,
		}},
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}

//
func TransTradForExpression(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	// [0]: label or placeholder
	// [1]: expression

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return nil, errors.New(emsg.InternalErr00001)
	}

	args := asts[1].Value.(AstSlice)
	// [0]: variables definition
	// [1]: condition
	// [2]: inclement/declement/etc
	// [3]: body

	cond, ok := promoteUnaryOperand(ctx, true, false, AstType_Bool, args[1])
	if !ok {
		return nil, errors.New(emsg.RuleErr00004)
	}

	slice := make(AstSlice, 4, 4)
	// [0]: variables definition, condition
	// [1]: body
	// [2]: inclement/declement/etc
	// [3]: label

	slice[3] = asts[0]

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

	retType, _, _ := getReturnTypeAndFlags(xctx, false, false, false, args[3])

	// TODO: traverse children break statements and transform them to add type conversion operations.

	slice[1] = Ast{
		OpCode:    mnem.Quote | retType,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: args[3]},
	}

	slice[2] = Ast{
		OpCode:    mnem.Quote,
		ClassName: clsz.C_op_quote,
		Type:      AstType_AstCons,
		Value:     AstCons{Car: args[2]},
	}

	return AstSlice{Ast{
		OpCode:    mnem.Scope | retType,
		ClassName: clsz.C_op_scope,
		Type:      AstType_AstCons,
		Value: AstCons{Car: Ast{
			OpCode:    mnem.For | retType,
			ClassName: clsz.C_op_for,
			Type:      AstType_ListOfAst,
			Value:     slice,
		}},
		SourcePosition: asts[0].SourcePosition,
	}}, nil
}
