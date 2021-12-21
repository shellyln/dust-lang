package scripting_test

import (
	"reflect"
	"testing"

	. "github.com/shellyln/dust-lang/scripting"
	xtor "github.com/shellyln/dust-lang/scripting/executor"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

func astSliceEquals(a AstSlice, b AstSlice) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b == nil || a == nil && b != nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a.ItemEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

func makeXctx() *xtor.ExecutionContext {
	xctx := NewExecutionContext(
		map[string]*VariableInfo{
			"a": {
				Flags: mnem.ReturnInt | mnem.Lvalue,
				Value: int64(1),
			},
			"aa": {
				Flags: mnem.ReturnInt | mnem.Lvalue,
				Value: int64(1),
			},
			"aaa": {
				Flags: mnem.ReturnInt | mnem.Lvalue,
				Value: int64(1),
			},
			"foo": {
				Flags: mnem.ReturnInt | mnem.Lvalue,
				Value: int64(1),
			},
			"bar": {
				Flags: mnem.ReturnInt | mnem.Lvalue,
				Value: int64(1),
			},
			"baz": {
				Flags: mnem.ReturnInt,
				Value: int64(1),
			},
			"qux": {
				Flags: mnem.ReturnFloat | mnem.Callable,
				Value: func() (float64, int, error) {
					return 17, 0, nil
				},
			},
			"sum": {
				Flags: mnem.ReturnFloat | mnem.Callable,
				Value: func(q []int64, o map[string]interface{}, s string, a float64, b ...float64) (float64, int, error) {
					v := a
					for _, w := range b {
						v += w
					}
					return v, 0, nil
				},
			},
		},
	)
	return xctx
}

var dummySlice AstSlice = AstSlice{}

type args struct {
	s string
}

type testMatrixItem struct {
	name             string
	args             args
	checkCompileWant bool
	compileWant      Ast
	compileWantErr   bool
	execWant         interface{}
	execWantErr      bool
	times            int
	breakpoint       bool
}

func runMatrix(t *testing.T, tests []testMatrixItem) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.breakpoint {
				zzz := 0
				zzz++
			}

			xctx := makeXctx()

			compileGot, err := Compile(xctx, tt.args.s)
			if (err != nil) != tt.compileWantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.compileWantErr)
				return
			}
			if tt.compileWantErr {
				return
			}
			if tt.checkCompileWant {
				if !dummySlice.ItemEquals(compileGot, tt.compileWant) {
					t.Errorf("Compile() = %v, want %v", compileGot, tt.compileWant)
					return
				}
			}

			n := tt.times
			if n == 0 {
				n++
			}

			for i := 0; i < n; i++ {
				execGot, err := ExecuteAst(xctx, compileGot)
				if (err != nil) != tt.execWantErr {
					t.Errorf("ExecuteAst() error = %v, wantErr %v", err, tt.execWantErr)
					return
				}
				if !reflect.DeepEqual(execGot, tt.execWant) {
					t.Errorf("ExecuteAst() = %v, want %v", execGot, tt.execWant)
				}
			}
		})
	}
}
