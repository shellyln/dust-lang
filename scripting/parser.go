package scripting

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	emsg "github.com/shellyln/dust-lang/scripting/errors"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	. "github.com/shellyln/dust-lang/scripting/parser"
	. "github.com/shellyln/takenoco/base"
	. "github.com/shellyln/takenoco/string"
)

var (
	programParser ParserFn
)

func init() {
	programParser = Program()
}

//
type VariableInfoMap = map[string]*xtor.VariableInfo

//
func CloneDict(src VariableInfoMap) VariableInfoMap {
	dest := make(VariableInfoMap)
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

//
func MergeDict(dest VariableInfoMap, srcDicts ...VariableInfoMap) VariableInfoMap {
	for _, src := range srcDicts {
		for key, value := range src {
			dest[key] = value
		}
	}
	return dest
}

//
func NewExecutionContext(dicts ...VariableInfoMap) *ExecutionContext {
	dict := MergeDict(VariableInfoMap{}, dicts...)
	xctx := xtor.NewExecutionContext(&ExecutionScope{
		Dict: dict,
	})
	xctx.ResetScope()
	return xctx
}

//
func Compile(xctx *ExecutionContext, s string) (Ast, error) {
	out, err := programParser(*NewStringParserContextWithTag(s, xctx))
	if err != nil {
		return Ast{}, errors.New(
			err.Error() +
				"\nParse failed at Col " +
				strconv.Itoa(out.SourcePosition.Position) + "\n" +
				out.Str[out.SourcePosition.Position:],
		)
	} else {
		if out.MatchStatus == MatchStatus_Matched {
			return out.AstStack[0], nil
		} else {
			return Ast{}, errors.New(
				emsg.ParserErr00001 + "\nParse failed at Col " +
					strconv.Itoa(out.SourcePosition.Position) + "\n" +
					out.Str[out.SourcePosition.Position:],
			)
		}
	}
}

// TODO: RegisterModule(xctx, ast, name)

//
func ExecuteAst(xctx *ExecutionContext, ast Ast) (interface{}, error) {
	ret, err := xtor.Execute(xctx, ast)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//
func Execute(xctx *ExecutionContext, s string) (interface{}, error) {
	ast, err := Compile(xctx, s)
	if err != nil {
		return nil, err
	}
	ret, err := xtor.Execute(xctx, ast)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//
func transformObject(outPtr interface{}, in interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(outPtr)
}

//
func Unmarshal(xctx *ExecutionContext, outPtr interface{}, s string) error {
	v, err := Execute(xctx, s)
	if err != nil {
		return err
	}
	transformObject(outPtr, v)
	return nil
}

// TODO: Call the function from host
func Call(xctx *ExecutionContext, fn string, params ...interface{}) (interface{}, error) {
	return nil, nil
}
