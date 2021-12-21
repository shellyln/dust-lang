package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/dust-lang/scripting/zeros"
	. "github.com/shellyln/takenoco/base"
)

//
func promoteUnaryOperand(
	ctx ParserContext, forcePromoting bool, disallowNarrowing bool,
	maxType AstType, op Ast) (Ast, bool) {

	// Meaning of `Lvalue|Callable|Maybe|Indexable|Bits8|ReturnUint` is:
	// The address of the variable that can be assigned to and
	// It is function that returns
	//   Option<?> of
	//   Array<?> of
	//   Uint8

	orig := op
	var max AstType
	var isSymbol bool

	var opReturnType AstOpCodeType

	if op.OpCode&mnem.OpCodeMask == mnem.Symbol {
		xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
		if !ok {
			return orig, false // TODO: BUG: Error!
		}

		vi, ok := xctx.GetVariableInfo(op.Value.(string))
		if !ok {
			return orig, false // TODO: BUG: Error!
		}

		isSymbol = true

		// NOTE: Promoted PRIMITIVE operands SHOULD NOT have bit length info.
		op.OpCode = mnem.Symbol | vi.Flags&(mnem.ReturnTypeMask|mnem.MarkerMask|mnem.TypeIdMask)

		if vi.Flags&mnem.TypeMarkerMask == 0 {
			opReturnType = vi.Flags & mnem.ReturnTypeMask
		}
	} else {
		if op.OpCode&mnem.TypeMarkerMask == 0 {
			opReturnType = op.OpCode & mnem.ReturnTypeMask
		}
	}

	if forcePromoting {
		max = maxType
	} else {
		max = xtor.ToAstType(opReturnType)
		if maxType < max {
			if disallowNarrowing {
				return orig, false
			}
			max = maxType
		}
	}

	ok := true

	switch max {
	case AstType_Int:
		if opReturnType == mnem.ReturnInt {

			// nothing to do

		} else if !isSymbol &&
			mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(opReturnType == mnem.ReturnUint ||
				opReturnType == mnem.ReturnFloat ||
				opReturnType == mnem.ReturnBool ||
				opReturnType == mnem.ReturnString) {

			op, ok = xtor.CastNumberLiteral(op, AstType_Int)
			if ok {
				op.OpCode = mnem.Imm_i64 | AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id)<<mnem.TypeIdOffset
				op.ClassName = clsz.C_imm_i64
			}

		} else if opReturnType == mnem.ReturnUint {
			op = Ast{
				OpCode:    mnem.Convu_i | AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convu_i,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnFloat {
			op = Ast{
				OpCode:    mnem.Convf_i | AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convf_i,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnBool {
			op = Ast{
				OpCode:    mnem.Convbool_i | AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convbool_i,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnString {
			op = Ast{
				OpCode:    mnem.Convstr_i | AstOpCodeType(tsys.TypeInfo_Primitive_I64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convstr_i,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else {
			return orig, false
		}

	case AstType_Uint:
		if opReturnType == mnem.ReturnUint {

			// nothing to do

		} else if !isSymbol &&
			mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(opReturnType == mnem.ReturnInt ||
				opReturnType == mnem.ReturnFloat ||
				opReturnType == mnem.ReturnBool ||
				opReturnType == mnem.ReturnString) {

			op, ok = xtor.CastNumberLiteral(op, AstType_Uint)
			if ok {
				op.OpCode = mnem.Imm_u64 | AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id)<<mnem.TypeIdOffset
				op.ClassName = clsz.C_imm_u64
			}

		} else if opReturnType == mnem.ReturnInt {
			op = Ast{
				OpCode:    mnem.Convi_u | AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convi_u,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnFloat {
			op = Ast{
				OpCode:    mnem.Convf_u | AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convf_u,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnBool {
			op = Ast{
				OpCode:    mnem.Convbool_u | AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convbool_u,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnString {
			op = Ast{
				OpCode:    mnem.Convstr_u | AstOpCodeType(tsys.TypeInfo_Primitive_U64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convstr_u,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else {
			return orig, false
		}

	case AstType_Float:
		if opReturnType == mnem.ReturnFloat {

			// nothing to do

		} else if !isSymbol &&
			mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(opReturnType == mnem.ReturnInt ||
				opReturnType == mnem.ReturnUint ||
				opReturnType == mnem.ReturnBool ||
				opReturnType == mnem.ReturnString) {

			op.OpCode = mnem.Imm_f64 | AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id)<<mnem.TypeIdOffset
			if ok {
				op.ClassName = clsz.C_imm_f64
				op, ok = xtor.CastNumberLiteral(op, AstType_Float)
			}

		} else if opReturnType == mnem.ReturnInt {
			op = Ast{
				OpCode:    mnem.Convi_f | AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convi_f,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnUint {
			op = Ast{
				OpCode:    mnem.Convu_f | AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convu_f,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnBool {
			op = Ast{
				OpCode:    mnem.Convbool_f | AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convbool_f,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnString {
			op = Ast{
				OpCode:    mnem.Convstr_f | AstOpCodeType(tsys.TypeInfo_Primitive_F64.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convstr_f,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else {
			return orig, false
		}

	case AstType_Bool:
		if opReturnType == mnem.ReturnBool {

			// nothing to do

		} else if !isSymbol &&
			mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(opReturnType == mnem.ReturnInt ||
				opReturnType == mnem.ReturnUint ||
				opReturnType == mnem.ReturnFloat ||
				opReturnType == mnem.ReturnString) {

			op.OpCode = mnem.Imm_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id)<<mnem.TypeIdOffset
			if ok {
				op.ClassName = clsz.C_imm_bool
				op, ok = xtor.CastNumberLiteral(op, AstType_Bool)
			}

		} else if mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(op.OpCode&mnem.OpCodeMask == mnem.Imm_nil ||
				op.OpCode&mnem.OpCodeMask == mnem.Imm_unitval) {

			// NOTE: nil and unit value are to be false.
			op = FalseAst

		} else if opReturnType == mnem.ReturnInt {
			op = Ast{
				OpCode:    mnem.Convi_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convi_bool,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnUint {
			op = Ast{
				OpCode:    mnem.Convu_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convu_bool,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnFloat {
			op = Ast{
				OpCode:    mnem.Convf_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convf_bool,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnString {
			op = Ast{
				OpCode:    mnem.Convstr_bool | AstOpCodeType(tsys.TypeInfo_Primitive_Bool.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convstr_bool,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else {
			return orig, false
		}

	case AstType_String:
		if opReturnType == mnem.ReturnString {

			// nothing to do

		} else if !isSymbol &&
			mnem.Imm_begin < op.OpCode && op.OpCode < mnem.Imm_end &&
			(opReturnType == mnem.ReturnInt ||
				opReturnType == mnem.ReturnUint ||
				opReturnType == mnem.ReturnFloat ||
				opReturnType == mnem.ReturnBool) {

			op.OpCode = mnem.Imm_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id)<<mnem.TypeIdOffset
			if ok {
				op.ClassName = clsz.C_imm_str
				op, ok = xtor.CastNumberLiteral(op, AstType_String)
			}

		} else if opReturnType == mnem.ReturnInt {
			op = Ast{
				OpCode:    mnem.Convi_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convi_str,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnUint {
			op = Ast{
				OpCode:    mnem.Convu_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convu_str,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnFloat {
			op = Ast{
				OpCode:    mnem.Convf_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convf_str,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else if opReturnType == mnem.ReturnBool {
			op = Ast{
				OpCode:    mnem.Convbool_str | AstOpCodeType(tsys.TypeInfo_Primitive_String.Id)<<mnem.TypeIdOffset,
				ClassName: clsz.C_op_convbool_str,
				Type:      AstType_AstCons,
				Value:     AstCons{Car: op},
			}

		} else {
			return orig, false
		}
	}
	return op, ok
}

//
func promoteBinaryOperands(
	ctx ParserContext, forcePromoting bool, disallowNarrowing bool,
	maxType AstType, op1 Ast, op2 Ast) ([2]Ast, bool) {

	orig := [2]Ast{op1, op2}
	ret := orig

	var max AstType

	if forcePromoting {
		max = maxType
	} else {
		for i := 0; i < 2; i++ {
			opReturnType := ret[i].OpCode & mnem.ReturnTypeMask

			switch ret[i].OpCode & mnem.OpCodeMask {
			case mnem.Symbol:
				if xctx, ok := ctx.Tag.(*xtor.ExecutionContext); ok {
					if vi, ok := xctx.GetVariableInfo(ret[i].Value.(string)); ok {
						if z := xtor.ToAstType(vi.Flags); max < z {
							max = z
						}
					}
				}
			case mnem.Imm_data: // TODO: ptr, nil, unit
				// TODO: ?
				// nothing to do
			default:
				if ret[i].OpCode&mnem.TypeMarkerMask == 0 {
					if opReturnType == mnem.ReturnInt {
						if max < AstType_Int {
							max = AstType_Int
						}
					} else if opReturnType == mnem.ReturnUint {
						if max < AstType_Uint {
							max = AstType_Uint
						}
					} else if opReturnType == mnem.ReturnFloat {
						if max < AstType_Float {
							max = AstType_Float
						}
					} else if opReturnType == mnem.ReturnBool {
						if max < AstType_Bool {
							max = AstType_Bool
						}
					} else if opReturnType == mnem.ReturnString {
						if max < AstType_String {
							max = AstType_String
						}
					}
				}
			}
		}
		if maxType < max {
			if disallowNarrowing {
				return orig, false
			}
			max = maxType
		}
	}

	for i := 0; i < 2; i++ {
		var ok bool
		ret[i], ok = promoteUnaryOperand(ctx, true, disallowNarrowing, max, ret[i])
		if !ok {
			return orig, false
		}
	}
	return ret, true
}

//
func promoteUnaryOperandForBitwise(ctx ParserContext, op Ast) (Ast, bool) {
	if opx, ok := promoteUnaryOperand(ctx, false, true, AstType_Uint, op); ok {
		return opx, true
	}

	opxf, ok := promoteUnaryOperand(ctx, true, true, AstType_Float, op)
	if !ok {
		return op, false
	}

	opx, ok := promoteUnaryOperand(ctx, true, false, AstType_Int, opxf)
	if !ok {
		return op, false
	}

	return opx, true
}

//
func promoteBinaryOperandsForBitwise(ctx ParserContext, op1 Ast, op2 Ast) ([2]Ast, bool) {
	orig := [2]Ast{op1, op2}

	ops, ok := promoteBinaryOperands(ctx, false, true, AstType_Uint, op1, op2)
	if ok {
		return ops, true
	} else {
		op1f, ok1 := promoteUnaryOperand(ctx, true, true, AstType_Float, op1)
		if ok1 {
			op1, ok1 = promoteUnaryOperand(ctx, true, false, AstType_Int, op1f)
			if !ok1 {
				return orig, false
			}
		} else {
			return orig, false
		}

		op2f, ok2 := promoteUnaryOperand(ctx, true, true, AstType_Float, op2)
		if ok2 {
			op2, ok2 = promoteUnaryOperand(ctx, true, false, AstType_Int, op2f)
			if !ok2 {
				return orig, false
			}
		} else {
			return orig, false
		}

		return [2]Ast{op1, op2}, true
	}
}

// op1 and op2 should have the same type. (promoted)
func generateBinaryOpResult(ctx ParserContext, opcode AstOpCodeType, op1 Ast, op2 Ast) (Ast, error) {
	ret := Ast{
		Type:           AstType_AstCons,
		Value:          AstCons{Car: op1, Cdr: op2},
		SourcePosition: op1.SourcePosition,
	}

	xctx, ok := ctx.Tag.(*xtor.ExecutionContext)
	if !ok {
		return ret, errors.New(emsg.InternalErr00001)
	}

	retType, matched, err := getReturnTypeAndFlags(xctx, false, false, false, op1)
	if err != nil {
		return ret, err
	}
	if !matched {
		return ret, errors.New(emsg.RuleErr00001)
	}

	// NOTE: Keep `opcode` side storage markers
	ret.OpCode = opcode&^(mnem.ReturnTypeMask|mnem.BitLenMask|mnem.TypeMarkerMask|mnem.TypeIdMask) |
		retType&(mnem.ReturnTypeMask|mnem.BitLenMask|mnem.TypeMarkerMask|mnem.TypeIdMask)

	return ret, nil
}

//
func generateUnaryOpResult(ctx ParserContext, opcode AstOpCodeType, op1 Ast) (Ast, error) {
	return generateBinaryOpResult(ctx, opcode, op1, Ast{})
}
