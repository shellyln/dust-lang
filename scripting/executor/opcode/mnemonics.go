package mnemomics

import (
	base "github.com/shellyln/takenoco/base"
)

type AstOpCodeType = base.AstOpCodeType

const (
	ReturnAny AstOpCodeType = 0

	// NOTE: Same value as AstType except `any`.
	//       AstOpCodeType(0): any
	//       AstType(0)      : nil
	//       AstType_Any     : any
	ReturnRune   = AstOpCodeType(base.AstType_Rune)
	ReturnInt    = AstOpCodeType(base.AstType_Int)
	ReturnUint   = AstOpCodeType(base.AstType_Uint)
	ReturnFloat  = AstOpCodeType(base.AstType_Float)
	ReturnBool   = AstOpCodeType(base.AstType_Bool)
	ReturnString = AstOpCodeType(base.AstType_String)

	ReturnTypeBits int           = 4
	ReturnTypeMask AstOpCodeType = (1 << ReturnTypeBits) - 1

	// Meaning of `Lvalue|Callable|Maybe|Indexable|Bits8|ReturnUint` is:
	// The address of the variable that can be assigned to and
	// It is function that returns
	//   Option<?> of
	//   Array<?> of
	//   Uint8

	FlagsBits         int           = 30
	Lvalue            AstOpCodeType = 1 << (ReturnTypeBits + 0)
	Callable          AstOpCodeType = 1 << (ReturnTypeBits + 1)
	Maybe             AstOpCodeType = 1 << (ReturnTypeBits + 2)
	Indexable         AstOpCodeType = 1 << (ReturnTypeBits + 3)
	MarkerMask        AstOpCodeType = 0x0f << (ReturnTypeBits + 0)
	StorageMarkerMask AstOpCodeType = 0x01 << (ReturnTypeBits + 0) // Lvalue
	TypeMarkerMask    AstOpCodeType = 0x0e << (ReturnTypeBits + 0) // Callable, Maybe, Indexable
	Bits8             AstOpCodeType = 1 << (ReturnTypeBits + 4)
	Bits16            AstOpCodeType = 2 << (ReturnTypeBits + 4)
	Bits32            AstOpCodeType = 3 << (ReturnTypeBits + 4)
	BitLenMask        AstOpCodeType = Bits32
	TypeIdOffset      int           = (ReturnTypeBits + 6)
	TypeIdMask        AstOpCodeType = 0x00ffffff << TypeIdOffset
	FlagsMask         AstOpCodeType = ((1 << (ReturnTypeBits + FlagsBits)) - 1) &^ ReturnTypeMask

	MetaInfoBits int           = ReturnTypeBits + FlagsBits
	MetaInfoMask AstOpCodeType = ReturnTypeMask | FlagsMask
	OpCodeMask   AstOpCodeType = ^(ReturnTypeMask | FlagsMask)
)

const (
	Nop AstOpCodeType = (iota << MetaInfoBits)

	Symbol,
	Symbol_i,
	Symbol_u,
	Symbol_f,
	Symbol_bool,
	Symbol_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Assign_begin AstOpCodeType = (iota << MetaInfoBits)

	Assign,
	Assign_i,
	Assign_u,
	Assign_f,
	Assign_bool,
	Assign_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	DefVar,
	DefVar_i,
	DefVar_u,
	DefVar_f,
	DefVar_bool,
	DefVar_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	DefConst,
	DefConst_i,
	DefConst_u,
	DefConst_f,
	DefConst_bool,
	DefConst_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Arg,
	Arg_i,
	Arg_u,
	Arg_f,
	Arg_bool,
	Arg_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Assign_end    AstOpCodeType = (iota << MetaInfoBits)
	Control_begin AstOpCodeType = Assign_end

	Quote,
	Quote_i,
	Quote_u,
	Quote_f,
	Quote_bool,
	Quote_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Scope,
	Scope_i,
	Scope_u,
	Scope_f,
	Scope_bool,
	Scope_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Jmp AstOpCodeType = (iota << MetaInfoBits)

	JmpT AstOpCodeType = (iota << MetaInfoBits)

	If,
	If_i,
	If_u,
	If_f,
	If_bool,
	If_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	While,
	While_i,
	While_u,
	While_f,
	While_bool,
	While_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	DoWhile,
	DoWhile_i,
	DoWhile_u,
	DoWhile_f,
	DoWhile_bool,
	DoWhile_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	For,
	For_i,
	For_u,
	For_f,
	For_bool,
	For_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForIn,
	ForIn_i,
	ForIn_u,
	ForIn_f,
	ForIn_bool,
	ForIn_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForIni,
	ForIni_i,
	ForIni_u,
	ForIni_f,
	ForIni_bool,
	ForIni_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForInu,
	ForInu_i,
	ForInu_u,
	ForInu_f,
	ForInu_bool,
	ForInu_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForInf,
	ForInf_i,
	ForInf_u,
	ForInf_f,
	ForInf_bool,
	ForInf_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForInbool,
	ForInbool_i,
	ForInbool_u,
	ForInbool_f,
	ForInbool_bool,
	ForInbool_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	ForInstr,
	ForInstr_i,
	ForInstr_u,
	ForInstr_f,
	ForInstr_bool,
	ForInstr_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Break,
	Break_i,
	Break_u,
	Break_f,
	Break_bool,
	Break_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Continue AstOpCodeType = (iota << MetaInfoBits)

	Ret,
	Ret_i,
	Ret_u,
	Ret_f,
	Ret_bool,
	Ret_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Seq,
	Seq_i,
	Seq_u,
	Seq_f,
	Seq_bool,
	Seq_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Last,
	Last_i,
	Last_u,
	Last_f,
	Last_bool,
	Last_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Control_end AstOpCodeType = (iota << MetaInfoBits)
	Call_begin  AstOpCodeType = Control_end

	Call,
	Call_i,
	Call_u,
	Call_f,
	Call_bool,
	Call_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	NativeCall,
	NativeCall_i,
	NativeCall_u,
	NativeCall_f,
	NativeCall_bool,
	NativeCall_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	DynCall,
	DynCall_i,
	DynCall_u,
	DynCall_f,
	DynCall_bool,
	DynCall_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Call_end     AstOpCodeType = (iota << MetaInfoBits)
	Object_begin AstOpCodeType = Call_end

	Func,
	Func_i,
	Func_u,
	Func_f,
	Func_bool,
	Func_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Lambda,
	Lambda_i,
	Lambda_u,
	Lambda_f,
	Lambda_bool,
	Lambda_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	List,
	List_i,
	List_u,
	List_f,
	List_bool,
	List_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	FilledList,
	FilledList_i,
	FilledList_u,
	FilledList_f,
	FilledList_bool,
	FilledList_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Index,
	Index_i,
	Index_u,
	Index_f,
	Index_bool,
	Index_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Slice,
	Slice_i,
	Slice_u,
	Slice_f,
	Slice_bool,
	Slice_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Object AstOpCodeType = (iota << MetaInfoBits)

	Mapindex,
	Mapindex_i,
	Mapindex_u,
	Mapindex_f,
	Mapindex_bool,
	Mapindex_str = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Object_end       AstOpCodeType = (iota << MetaInfoBits)
	Arithmetic_begin AstOpCodeType = Object_end

	PreIncr,
	PreIncr_i,
	PreIncr_u,
	PreIncr_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	PreDecr,
	PreDecr_i,
	PreDecr_u,
	PreDecr_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	PostIncr,
	PostIncr_i,
	PostIncr_u,
	PostIncr_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	PostDecr,
	PostDecr_i,
	PostDecr_u,
	PostDecr_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Neg,
	Neg_i,
	Neg_u,
	Neg_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Pow,
	Pow_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnFloat

	Mul,
	Mul_i,
	Mul_u,
	Mul_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Div,
	Div_i,
	Div_u,
	Div_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Mod,
	Mod_i,
	Mod_u,
	Mod_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Add,
	Add_i,
	Add_u,
	Add_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Sub,
	Sub_i,
	Sub_u,
	Sub_f AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Concat,
	Concat_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnString

	Arithmetic_end          AstOpCodeType = (iota << MetaInfoBits)
	BitwiseAndLogical_begin AstOpCodeType = Arithmetic_end

	LogicalNotbool,
	LogicalNotbool_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	LogicalAndbool,
	LogicalAndbool_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	LogicalOrbool,
	LogicalOrbool_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	BitwiseNot,
	BitwiseNot_i,
	BitwiseNot_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseLShift,
	BitwiseLShift_i,
	BitwiseLShift_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseSRShift,
	BitwiseSRShift_i,
	BitwiseSRShift_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseURShift,
	BitwiseURShift_i,
	BitwiseURShift_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseAnd,
	BitwiseAnd_i,
	BitwiseAnd_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseXor,
	BitwiseXor_i,
	BitwiseXor_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	BitwiseOr,
	BitwiseOr_i,
	BitwiseOr_u AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	CmpEqi,
	CmpEqi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpEqu,
	CmpEqu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpEqf,
	CmpEqf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpEqbool,
	CmpEqbool_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpEqstr,
	CmpEqstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpEqdyn,
	CmpEqdyn_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpNotEqi,
	CmpNotEqi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpNotEqu,
	CmpNotEqu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpNotEqf,
	CmpNotEqf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpNotEqbool,
	CmpNotEqbool_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpNotEqstr,
	CmpNotEqstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpNotEqdyn,
	CmpNotEqdyn_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpStrictEqdyn,
	CmpStrictEqdyn_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpStrictNotEqdyn,
	CmpStrictNotEqdyn_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpLTi,
	CmpLTi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLTu,
	CmpLTu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLTf,
	CmpLTf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLTstr,
	CmpLTstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpLEi,
	CmpLEi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLEu,
	CmpLEu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLEf,
	CmpLEf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpLEstr,
	CmpLEstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpGTi,
	CmpGTi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGTu,
	CmpGTu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGTf,
	CmpGTf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGTstr,
	CmpGTstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	CmpGEi,
	CmpGEi_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGEu,
	CmpGEu_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGEf,
	CmpGEf_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool
	CmpGEstr,
	CmpGEstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnBool

	BitwiseAndLogical_end AstOpCodeType = (iota << MetaInfoBits)
	Conv_begin            AstOpCodeType = BitwiseAndLogical_end

	Convi,
	Convi_u,
	Convi_f,
	Convi_bool,
	Convi_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Convu,
	Convu_i,
	Convu_f,
	Convu_bool,
	Convu_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Convf,
	Convf_i,
	Convf_u,
	Convf_bool,
	Convf_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Convbool,
	Convbool_i,
	Convbool_u,
	Convbool_f,
	Convbool_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnString

	Convstr,
	Convstr_i,
	Convstr_u,
	Convstr_f,
	Convstr_bool AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool

	Convdyn,
	Convdyn_i,
	Convdyn_u,
	Convdyn_f,
	Convdyn_bool,
	Convdyn_str AstOpCodeType = (iota << MetaInfoBits),
		(iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Conv_end  AstOpCodeType = (iota << MetaInfoBits)
	Imm_begin AstOpCodeType = Conv_end

	Imm_i64,
	Imm_u64,
	Imm_f64,
	Imm_bool,
	Imm_str AstOpCodeType = (iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat,
		(iota << MetaInfoBits) + ReturnBool,
		(iota << MetaInfoBits) + ReturnString

	Imm_i32,
	Imm_u32,
	Imm_f32 AstOpCodeType = (iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint,
		(iota << MetaInfoBits) + ReturnFloat

	Imm_i16,
	Imm_u16 AstOpCodeType = (iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	Imm_i8,
	Imm_u8 AstOpCodeType = (iota << MetaInfoBits) + ReturnInt,
		(iota << MetaInfoBits) + ReturnUint

	Imm_ptr     AstOpCodeType = (iota << MetaInfoBits) // Value is unsafe.Pointer
	Imm_data    AstOpCodeType = (iota << MetaInfoBits) // Value is any interface{} (include nil)
	Imm_nil     AstOpCodeType = (iota << MetaInfoBits) // Value is nil
	Imm_unitval AstOpCodeType = (iota << MetaInfoBits) // Value is &UnitSingleton

	Imm_end AstOpCodeType = (iota << MetaInfoBits)
)
