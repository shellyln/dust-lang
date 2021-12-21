package rules

import (
	"errors"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	clsz "github.com/shellyln/dust-lang/scripting/parser/classes"
	. "github.com/shellyln/takenoco/base"
)

// Logical NOT `! op1`
// Bitwise NOT `~ op1`
// Unary plus `+ op1`
// Unary negation `- op1`
// Prefix Increment `++ op1`
// Prefix Decrement `-- op1`
var precedence17 = Precedence{
	Rules: []ParserFn{
		Trans(
			FlatGroup(
				isOperator(clsz.UnaryPrefixOp, []string{"!"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok1 := promoteUnaryOperand(ctx, true, false, AstType_Bool, asts[1])

				// NOTE: ptr, any, nil, unit value are not promoted.

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnBool {

					return AstSlice{{
						OpCode:         mnem.Imm_bool,
						ClassName:      clsz.C_imm_bool,
						Type:           AstType_Bool,
						Value:          !op1.Value.(bool),
						SourcePosition: op1.SourcePosition,
					}}, nil

				} else {
					// NOTE: Imm_ptr and Imm_data will be compared at runtime

					if !ok1 {
						op1 = Ast{
							OpCode:    mnem.Convdyn_bool,
							ClassName: "$convdyn_bool",
							Type:      AstType_AstCons,
							Value:     AstCons{Car: op1},
						}
					}

					return AstSlice{{
						OpCode:         mnem.LogicalNotbool_bool,
						ClassName:      "$logicalnot_bool",
						Type:           AstType_AstCons,
						Value:          AstCons{Car: op1},
						SourcePosition: op1.SourcePosition,
					}}, nil
				}
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.UnaryPrefixOp, []string{"~"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperandForBitwise(ctx, asts[1])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[0].Value.(string))
				}

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt {
					return AstSlice{{
						OpCode:         mnem.Imm_i64,
						ClassName:      clsz.C_imm_i64,
						Type:           AstType_Int,
						Value:          ^op1.Value.(int64),
						SourcePosition: op1.SourcePosition,
					}}, nil
				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {
					return AstSlice{{
						OpCode:         mnem.Imm_u64,
						ClassName:      clsz.C_imm_u64,
						Type:           AstType_Uint,
						Value:          ^op1.Value.(uint64),
						SourcePosition: op1.SourcePosition,
					}}, nil
				}

				ret, err := generateUnaryOpResult(ctx, mnem.BitwiseNot, op1)
				if err != nil {
					return nil, err
				}

				if ret.OpCode&mnem.ReturnTypeMask != mnem.ReturnAny {
					return AstSlice{ret}, nil
				} else {
					return nil, errors.New(emsg.RuleErr00006 + asts[0].Value.(string))
				}
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.UnaryPrefixOp, []string{"++", "--"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperand(ctx, false, true, AstType_Float, asts[1])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[0].Value.(string))
				}

				if op1.OpCode&mnem.Lvalue == 0 {
					return nil, errors.New(emsg.RuleErr00002 + asts[0].Value.(string))
				}

				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}

				var opcode AstOpCodeType
				var err error

				if asts[0].Value.(string) == "++" {
					opcode = mnem.PreIncr | mnem.Lvalue
				} else {
					opcode = mnem.PreDecr | mnem.Lvalue
				}

				ret, err = generateUnaryOpResult(ctx, opcode, op1)
				if err != nil {
					return nil, err
				}

				if ret.OpCode&mnem.ReturnTypeMask != mnem.ReturnAny {
					return AstSlice{ret}, nil
				} else {
					return nil, errors.New(emsg.RuleErr00006 + asts[0].Value.(string))
				}
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.UnaryPrefixOp, []string{"+"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperand(ctx, false, true, AstType_Float, asts[1])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[0].Value.(string))
				}

				op1ReturnType := op1.OpCode & mnem.ReturnTypeMask

				if op1ReturnType == mnem.ReturnInt ||
					op1ReturnType == mnem.ReturnUint ||
					op1ReturnType == mnem.ReturnFloat {

					return asts[1:2], nil
				} else {
					_, err := generateUnaryOpResult(ctx, mnem.Add, op1)
					if err != nil {
						return nil, err
					}
					return asts[1:2], nil
				}
			},
		),
		Trans(
			FlatGroup(
				isOperator(clsz.UnaryPrefixOp, []string{"-"}),
				anyOperand(),
			),
			func(ctx ParserContext, asts AstSlice) (AstSlice, error) {

				// TODO: BUG: Set type id.

				op1, ok := promoteUnaryOperand(ctx, false, true, AstType_Float, asts[1])
				if !ok {
					return nil, errors.New(emsg.RuleErr00005 + asts[0].Value.(string))
				}

				ret := Ast{
					SourcePosition: op1.SourcePosition,
				}

				opcode := mnem.Neg

				if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnInt {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Int
					ret.Value = -op1.Value.(int64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnUint {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Uint
					ret.Value = -op1.Value.(uint64)

				} else if mnem.Imm_begin < op1.OpCode && op1.OpCode < mnem.Imm_end &&
					op1.OpCode&mnem.ReturnTypeMask == mnem.ReturnFloat {

					ret.OpCode = op1.OpCode
					ret.ClassName = op1.ClassName
					ret.Type = AstType_Float
					ret.Value = -op1.Value.(float64)

				} else {
					var err error
					ret, err = generateUnaryOpResult(ctx, opcode, op1)
					if err != nil {
						return nil, err
					}
				}

				if ret.OpCode&mnem.ReturnTypeMask != mnem.ReturnAny {
					return AstSlice{ret}, nil
				} else {
					return nil, errors.New(emsg.RuleErr00006 + asts[0].Value.(string))
				}
			},
		),
	},
	Rtol: true,
}
