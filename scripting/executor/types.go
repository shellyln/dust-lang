package executor

import (
	. "github.com/shellyln/takenoco/base"
)

//
type returnPayload struct {
	Type  AstType
	Value interface{}
}

//
type breakPayload struct {
	Label string
	Type  AstType
	Value interface{}
}

//
type continuePayload struct {
	Label string
}

//
type thrownPayload struct {
	Value interface{}
}

//
type ExecutionScopeType int

const (
	//
	ExecutionScope_Normal ExecutionScopeType = iota
	//
	ExecutionScope_Function
	//
	ExecutionScope_Lambda
)

//
type ExecutionScope struct {
	Type       ExecutionScopeType
	ReturnType AstOpCodeType
	Dict       map[string]*VariableInfo
	Next       *ExecutionScope
	StackNext  *ExecutionScope
}
