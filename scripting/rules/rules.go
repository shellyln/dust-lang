package rules

import (
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/object"
)

//
var ProdRule TransformerFn = ProductionRule(
	[]Precedence{
		precedence20,
		precedence18,
		precedence17,
		precedence16,
		precedence15,
		precedence14,
		precedence13,
		precedence12,
		precedence11,
		precedence10,
		precedence09,
		precedence08,
		precedence07,
		precedence06,
		precedence04b,
		precedence04a,
		precedence03,
		precedence01,
	},
	FlatGroup(Start(), Any(), End()),
)
