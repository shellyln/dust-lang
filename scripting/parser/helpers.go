package parser

import (
	. "github.com/shellyln/takenoco/base"
)

//
func erase(fn ParserFn) ParserFn {
	return Trans(fn, Erase)
}

//
func ref(fn func() ParserFn) ParserFn {
	return Indirect(fn)
}

//
func setStr(s string) TransformerFn {
	return func(ctx ParserContext, asts AstSlice) (AstSlice, error) {
		return AstSlice{{
			ClassName: "setStr",
			Type:      AstType_String,
			Value:     s,
		}}, nil
	}
}

//
func replaceStr(fn ParserFn, s string) ParserFn {
	return Trans(fn, setStr(s))
}
