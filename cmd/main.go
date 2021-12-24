package dust

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/shellyln/dust-lang/scripting"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	parser "github.com/shellyln/takenoco/base"
)

type CommandFlags struct {
	Eval bool
}

// TODO: divide map
var functions = scripting.VariableInfoMap{
	"print": {
		Flags: mnem.ReturnAny | mnem.Callable,
		Value: func(x ...interface{}) {
			fmt.Print(x...)
		},
	},
	"println": {
		Flags: mnem.ReturnAny | mnem.Callable,
		Value: func(x ...interface{}) {
			fmt.Println(x...)
		},
	},
	"json_stringify": {
		Flags: mnem.ReturnString | mnem.Callable,
		Value: func(x interface{}) (string, error) {
			bytes, err := json.Marshal(x)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
	},
	"json_parse": {
		Flags: mnem.ReturnAny | mnem.Callable,
		Value: func(s string) (interface{}, error) {
			var z interface{}
			err := json.Unmarshal([]byte(s), &z)
			if err != nil {
				return nil, err
			}
			return z, nil
		},
	},
	"file_read_as_string": {
		Flags: mnem.ReturnString | mnem.Callable,
		Value: func(path string) (string, error) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		},
	},
	"file_read_as_bytes": {
		Flags: mnem.ReturnUint | mnem.Indexable | mnem.Bits8 | mnem.Callable,
		Value: func(path string) ([]byte, error) {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return bytes, nil
		},
	},
	// TODO: math constants
	"math_rand_u64": {
		Flags: mnem.ReturnUint | mnem.Callable,
		Value: func() uint64 {
			return rand.Uint64()
		},
	},
	"math_rand_f64": {
		Flags: mnem.ReturnFloat | mnem.Callable,
		Value: func() float64 {
			return rand.Float64()
		},
	},
	// TODO: assert_eq
	// TODO: len, cap, append // Needs generics
	// TODO: keys, values, items // Needs generics
	// TODO: slice(w/ cloning), map, filter, reduce, ... // Needs generics, calling ast from host func
	// TODO: http_fetch
	// TODO: regexp
	// TODO: execute
	// TODO: eval
	// TODO: csv
}

func SubcommandMain(flags CommandFlags, args []string) {
	rand.Seed(time.Now().UnixNano())

	// TODO: check permitions; `--allow="net,file,eval"`
	xctx := scripting.NewExecutionContext(functions)

	var ret interface{}
	var err error

	if flags.Eval {
		ret, err = scripting.Execute(xctx, strings.Join(args, " "))
	} else {
		slice := make(parser.AstSlice, len(args))
		modules := parser.Ast{
			OpCode: mnem.Last,
			Type:   parser.AstType_ListOfAst,
			Value:  slice,
		}
		for i := 0; i < len(args); i++ {
			bytes, err := ioutil.ReadFile(args[i])
			if err != nil {
				panic(err)
			}
			slice[i], err = scripting.Compile(xctx, string(bytes))
			if err != nil {
				panic(err)
			}
		}
		ret, err = scripting.ExecuteAst(xctx, modules)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else {
		fmt.Println(ret)
	}
}
