package parser

import (
	emsg "github.com/shellyln/dust-lang/scripting/errors"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

//
func lineComment() ParserFn {
	return FlatGroup(
		Seq("//"),
		ZeroOrMoreTimes(CharClassN("\r", "\n")),
		First(CharClass("\r\n", "\r", "\n"), End()),
	)
}

//
func hashLineComment() ParserFn {
	return FlatGroup(
		Seq("#"),
		ZeroOrMoreTimes(CharClassN("\r", "\n")),
		First(CharClass("\r\n", "\r", "\n"), End()),
	)
}

//
func blockComment() ParserFn {
	return FlatGroup(
		Seq("/*"),
		ZeroOrMoreTimes(CharClassN("*/")),
		First(
			FlatGroup(End(), Error(emsg.ParserErr00005)),
			CharClass("*/"),
		),
	)
}

//
func comment() ParserFn {
	return First(lineComment(), blockComment())
}

//
func sp0() ParserFn {
	return erase(ZeroOrMoreTimes(First(Whitespace(), comment())))
}

//
func sp1() ParserFn {
	return erase(OneOrMoreTimes(First(Whitespace(), comment())))
}
