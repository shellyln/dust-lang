package rules

import (
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/object"
)

//
func makeOpMatcher(className string, ops []string) func(c interface{}) bool {
	return func(c interface{}) bool {
		ast, ok := c.(Ast)
		if !ok || ast.ClassName != className {
			return false
		}
		val := ast.Value.(string)
		for _, op := range ops {
			if op == val {
				return true
			}
		}
		return false
	}
}

//
func unwrapOperandItem(ctx ParserContext, asts AstSlice) (AstSlice, error) {
	return AstSlice{asts[0].Value.(Ast)}, nil
}

//
func anyOperand() ParserFn {
	return Trans(Any(), unwrapOperandItem)
}

//
func isOperator(className string, ops []string) ParserFn {
	return Trans(ObjClassFn(makeOpMatcher(className, ops)), unwrapOperandItem)
}
