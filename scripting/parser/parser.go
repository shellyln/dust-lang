package parser

import (
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func Program() ParserFn {
	return FlatGroup(
		Start(),
		ZeroOrOnce(erase(hashLineComment())),
		sp0(),
		topLevelStatement(),
		End(),
	)
}
