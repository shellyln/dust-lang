//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	scripting "github.com/shellyln/dust-lang/scripting"
)

// NOTE: At this time, tinygo can't compile this.
func evalScript(this js.Value, args []js.Value) interface{} {
	src := ""
	if 0 < len(args) {
		src = args[0].String()
	}

	xctx := scripting.NewExecutionContext()
	ret, err := scripting.Execute(xctx, src)
	if err != nil {
		println(err)
	}

	return js.ValueOf(ret)
}

//export evalScriptWithTinyGo
func evalScriptWithTinyGo(src string) interface{} {
	// https://tinygo.org/docs/guides/webassembly/
	// NOTE: BUG: At this time, non-numeric type parameters are not accepted.
	xctx := scripting.NewExecutionContext()
	ret, err := scripting.Execute(xctx, src)
	if err != nil {
		println(err)
	}
	// NOTE: BUG: At this time, non-numeric type parameters are not accepted.
	return ret
}

func main() {
	ch := make(chan struct{}, 0)
	println("Go WebAssembly Initialized")

	// for golang/go
	js.Global().Set("evalScript", js.FuncOf(evalScript))

	// for golang/go
	<-ch
}
