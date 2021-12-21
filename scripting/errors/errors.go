package errors

const (
	InternalErr00001 = "[X00001] Internal error: ExecutionContext is not set."

	ParserErr00001 = "[P00001] Syntax error."
	ParserErr00002 = "[P00002] An unexpected character has appeared in the expression."
	ParserErr00003 = "[P00003] An unexpected character has appeared in the actual argument."
	ParserErr00004 = "[P00004] An unexpected character has appeared in the list literal."
	ParserErr00005 = "[P00005] An unexpected termination has appeared in the comment."
	ParserErr00006 = "[P00006] An unexpected termination has appeared in the string literal."
	ParserErr00007 = "[P00007] An unexpected newline has appeared in the string literal."
	ParserErr00008 = "[P00008] An unexpected character has appeared in the object literal."
	// ParserErr00009 = "[P00009]"
	ParserErr00010 = "[P00010] An unexpected character has appeared in the multiple statement."
	// ParserErr00011 = "[P00011]"
	ParserErr00012 = "[P00012] An unexpected character has appeared in the lambda definition."
	ParserErr00013 = "[P00013] An unexpected character has appeared in the function definition."

	RuleErr00001 = "[R00001] generateBinaryOpResult; have no return type: "
	RuleErr00002 = "[R00002] The operand must be a lvalue: "
	RuleErr00003 = "[R00003] unary promotion failed"
	RuleErr00004 = "[R00004] The operand must be a bool: "
	RuleErr00005 = "[R00005] The operand must be a number: "
	RuleErr00006 = "[R00006] The type is ambiguous: "
	RuleErr00007 = "[R00007] Unknown operator: "
	RuleErr00008 = "[R00008] The operand must be a number or string: "
	RuleErr00009 = "[R00009] arguments length is not matched"
	RuleErr00010 = "[R00010] The operand must be a string: "
	RuleErr00011 = "[R00011] Division by zero has occurred: "
	RuleErr00012 = "[R00012] The object key must be a string."
	RuleErr00013 = "[R00013] The array index must be a number."
	RuleErr00014 = "[R00014] The array slice index must be a number."
	RuleErr00015 = "[R00015] Unknown keyword: "
	RuleErr00016 = "[R00016] Unknown assignment operator is specified: "
	RuleErr00017 = "[R00017] The type conversion is ambiguous: "
	RuleErr00018 = "[R00018] The variable is duplicated: "
	RuleErr00019 = "[R00019] unknown class token: "
	RuleErr00020 = "[R00020] GetReturnType; have no return type: "
	RuleErr00021 = "[R00021] The operand must be a primitive type: "

	ExecErr00001 = "[E00001] Division by zero has occurred: "
	ExecErr00002 = "[E00002] The variable is duplicated: "
	ExecErr00003 = "[E00003] cast failed: "
	ExecErr00004 = "[E00004] The variable can not defined: "
	ExecErr00005 = "[E00005] argument length is not matched: "
	ExecErr00006 = "[E00006] Panic! Type conversion (int) is failed: "
	ExecErr00007 = "[E00007] Panic! Type conversion (uint) is failed: "
	ExecErr00008 = "[E00008] Panic! Type conversion (float) is failed: "
	ExecErr00009 = "[E00009] Panic! Type conversion (bool) is failed: "
	ExecErr00010 = "[E00010] Panic! Type conversion (string) is failed: "
	ExecErr00011 = "[E00011] Unknown opcode: "
	ExecErr00012 = "[E00012] Uncaught exception: "
	ExecErr00013 = "[E00013] An undefined symbol is referenced: "
	ExecErr00014 = "[E00014] The operand must be callable: "
	ExecErr00015 = "[E00015] The operand must be maybve: "
	ExecErr00016 = "[E00016] unknown type: "
	ExecErr00017 = "[E00017] out of index: "
	// ExecErr00018 = "[E00018] key is not defined: "
	ExecErr00019 = "[E00019] unknown type: "
	ExecErr00020 = "[E00020] variable is not defined: "
	ExecErr00021 = "[E00021] The operand must be indexable: "
	// ExecErr00022 = "[E00022] The operand must be not indexable: "
)
