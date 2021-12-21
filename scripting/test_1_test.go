package scripting_test

import (
	"math"
	"reflect"
	"testing"

	. "github.com/shellyln/dust-lang/scripting"
	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	. "github.com/shellyln/takenoco/base"
)

func TestParse0(t *testing.T) {
	xctx := NewExecutionContext(map[string]*VariableInfo{})

	got, err := Execute(xctx, "1 + 2")
	if err != nil {
		t.Errorf("ExecuteAst() error = %v, wantErr %v", err, false)
	}

	want := int64(3)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExecuteAst() = %v, want %v", got, want)
	}
}

func TestParse02(t *testing.T) {
	type Foobar struct {
		XFoo string `csv:"Foo" json:"foo"`
		Bar  string
		XBaz int `csv:"Baz" json:"baz"`
		Qux  uint
		Quux float64
		Test int `csv:"Test" json:"test"`
	}

	var out []Foobar
	xctx := NewExecutionContext(map[string]*VariableInfo{})

	err := Unmarshal(xctx, &out, `[{foo: 'abc', Qux: 3},{foo: 'cde', Qux: 5}]`)
	if err != nil {
		t.Errorf("ExecuteAst() error = %v, wantErr %v", err, false)
	}

	want := []Foobar{{XFoo: "abc", Qux: 3}, {XFoo: "cde", Qux: 5}}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("ExecuteAst() = %v, want %v", out, want)
	}
}

func TestParse(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "0",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "1qqqqqqa",
		args:             args{s: "let p:i8 = 1i8; p+3i8"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(4),
		execWantErr:      false,
	}, {
		name:             "1qqqqqqb",
		args:             args{s: "let p:u8 = 1u8; p+3u8"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(4),
		execWantErr:      false,
	}, {
		name:             "1qqqqqqc",
		args:             args{s: "let p:f32 = 1.0f32; p + 3.0f32"}, // TODO: BUG: syntax error on `1.0_f32`
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         float64(4),
		execWantErr:      false,
		// }, {
		// 	name:             "1ya1",
		// 	args:             args{s: "let p = u8![2,3,5,7,11] as any; p[1] + 1u8"}, // TODO: BUG: `as any have no effect`
		// 	checkCompileWant: false,
		// 	compileWantErr:   false,
		// 	execWant:         uint64(4),
		// 	execWantErr:      false,
		// }, {
		// 	name:             "1ya2",
		// 	args:             args{s: "let p:any = u8![2,3,5,7,11] as any; p[1] + 1u8"}, // TODO: p can't be cast to `[u8]``
		// 	checkCompileWant: false,
		// 	compileWantErr:   false,
		// 	execWant:         uint64(4),
		// 	execWantErr:      false,
	}, {
		name:             "1ya2a",
		args:             args{s: "let p:any = {x:11}; p.x = u8![2,3,5,7,11] as any; (p.x[1])as int + 1u8"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(4),
		execWantErr:      false,
	}, {
		name:             "1ya2b",
		args:             args{s: "let p:any = {}; p.x = u8![2,3,5,7,11] as any; (p.x[1])as int + 1u8"}, // TODO: BUG: Can't add new key `x`
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(4),
		execWantErr:      false,
	}, {
		name:             "1aa",
		args:             args{s: "let mut p: [any] = [1,2,3]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{int64(1), int64(2), int64(3)},
		execWantErr:      false,
	}, {
		name:             "1e1",
		args:             args{s: "let mut p = for let mut z = 0; n in i64![3, 5, 7] {z = z + n; z}; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
	}, {
		name:             "1e2",
		args:             args{s: "let s = i64![3, 5, 7]; let mut p = for let mut z = 0; n in s {z = z + n; z}; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
	}, {
		name:             "1e3",
		args:             args{s: "let mut p = for let mut z = 0; n in 1..3 {z = z + n; z}; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(6),
		execWantErr:      false,
	}, {
		name:             "1qa1",
		args:             args{s: "let mut r#p = 0; ++r#p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1),
		execWantErr:      false,
	}, {
		name:             "1qa1",
		args:             args{s: "let mut p = 0; ++r#p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1),
		execWantErr:      false,
	}, {
		name:             "1qa1",
		args:             args{s: "let mut p = 0; ++p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1),
		execWantErr:      false,
	}, {
		name:             "1qa2",
		args:             args{s: "{let mut p = 0; ++p}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1),
		execWantErr:      false,
	}, {
		name:             "1qb",
		args:             args{s: "let mut p = 0; p++; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1),
		execWantErr:      false,
	}, {
		name:             "1yb",
		args:             args{s: "let p = u8![2,3,5,7,11]; p[1] + 1u8"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(4),
		execWantErr:      false,
	}, {
		name:             "1yc",
		args:             args{s: "u8![2,3,5,7,11][1] + 1u8"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(4),
		execWantErr:      false,
	}, {
		name:             "1x",
		args:             args{s: "`abcde`[1]+100"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(198),
		execWantErr:      false,
	}, {
		name:             "1x",
		args:             args{s: "`abcde`[1:3]+`fghij`"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         "bcfghij",
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "#aaa\n\n/*aaa*/0?0?\n//aaa\n-(-(- -(-3)))/2+1:23:1?+17.00e+00:19/*bbb*/\n//aaa"},
		checkCompileWant: true,
		compileWant:      Ast{OpCode: mnem.Imm_f64, ClassName: "%imm_f64", Type: AstType_Float, Value: 17.0},
		compileWantErr:   false,
		execWant:         17.0,
		execWantErr:      false,
	}, {
		name:             "1b",
		args:             args{s: "\n#aaa\n/*aaa*/0?0?\n#aaa\n-(-(- -(-3)))/2+1:23:1?+17.00e+00:19/*bbb*/\n#aaa"},
		checkCompileWant: false,
		compileWantErr:   true,
		execWant:         17.0,
		execWantErr:      false,
	}, {
		name:             "2",
		args:             args{s: "((-(1+(2))+3**(-4.01/5+6)+7))?8+9*10:(11+(12-13)+14  )  "},
		checkCompileWant: true,
		compileWant:      Ast{OpCode: mnem.Imm_i64, ClassName: "%imm_i64", Type: AstType_Int, Value: int64(98)},
		compileWantErr:   false,
		execWant:         int64(98),
		execWantErr:      false,
	}, {
		name:             "3",
		args:             args{s: "11*13+(7+5)*-2**3"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         47.0,
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[3i32;5]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{int64(3), int64(3), int64(3), int64(3), int64(3)}, // NOTE: promotion
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[1,2,3,]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{int64(1), int64(2), int64(3)},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "u8![3i32;5]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []uint8{3, 3, 3, 3, 3},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "vec![1,2,3,]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{int64(1), int64(2), int64(3)},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "u8![1,2,3,]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []uint8{1, 2, 3},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[11,13,17][2]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(17),
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[0x01_u8,0b1010_1010,true,-Infinity-Infinity,]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{uint64(1), int64(170), true, math.Inf(-1)}, // NOTE: promotion
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "[[1],[2]]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{[]interface{}{int64(1)}, []interface{}{int64(2)}},
		execWantErr:      false,
	}, {
		name:             "4",
		args:             args{s: "vec![vec![1],[2]]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{[]interface{}{int64(1)}, []interface{}{int64(2)}},
		execWantErr:      false,
	}, {
		name:             "5",
		args:             args{s: "[[1],[2]][1][0]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(2),
		execWantErr:      false,
	}, {
		name:             "7",
		args:             args{s: "a=aa=12345;aa"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(12345),
		execWantErr:      false,
	}, {
		name:             "7",
		args:             args{s: "a=aa=12345;aa;"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(12345),
		execWantErr:      false,
	}, {
		name:             "7",
		args:             args{s: "a=aa=12345;;aa;"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(12345),
		execWantErr:      false,
	}, {
		name:             "12",
		args:             args{s: "aaa+++++17"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(18),
		execWantErr:      false,
	}, {
		name:             "13",
		args:             args{s: "aaa+++++17,aaa"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(3),
		execWantErr:      false,
	}, {
		name:             "14",
		args:             args{s: "aaa+++++17;aaa"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(3),
		execWantErr:      false,
	}, {
		name:             "16",
		args:             args{s: "let p: int = 11;"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(11),
		execWantErr:      false,
	}, {
		name:             "16",
		args:             args{s: "let mut p: int = 11; p = 13; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(13),
		execWantErr:      false,
	}, {
		name:             "16",
		args:             args{s: "let p: int = 11; p = 13; p"},
		checkCompileWant: false,
		compileWantErr:   true,
	}, {
		name:             "16",
		args:             args{s: "const p: int = 11; p = 13; p"},
		checkCompileWant: false,
		compileWantErr:   true,
	}, {
		name:             "17",
		args:             args{s: "let p: int = 11; let q: uint = 13; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(11),
		execWantErr:      false,
	}, {
		name:             "17",
		args:             args{s: "let p: int = 11; let q: uint = 13; p+17"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(28),
		execWantErr:      false,
	}, {
		name:             "18",
		args:             args{s: "let p: int = 11, q: uint = 13"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(13),
		execWantErr:      false,
	}, {
		name:             "19",
		args:             args{s: "let p = [11,13]; p[0] = 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         []interface{}{int64(17), int64(13)},
		execWantErr:      false,
	}, {
		name:             "36a",
		args:             args{s: "let p = 3 as float"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         3.0,
		execWantErr:      false,
	}, {
		name:             "36a",
		args:             args{s: "let p = (3 as float)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         3.0,
		execWantErr:      false,
	}, {
		name:             "36a",
		args:             args{s: "let p = (3)as float"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         3.0,
		execWantErr:      false,
	}, {
		name:             "36b",
		args:             args{s: "let p = 3 as int"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(3),
		execWantErr:      false,
	}, {
		name:             "37a",
		args:             args{s: "[2,3,5,7,11,13,17][1:2][0]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(3),
		execWantErr:      false,
	}, {
		name:             "37a",
		args:             args{s: "[2,3,5,7,11,13,17][1..2][0]"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(3),
		execWantErr:      false,
	}, {
		name:             "37b",
		args:             args{s: "let p: any = [1,2,3]; let q: any = p; p == q"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         true,
		execWantErr:      false,
	}, {
		name:             "37b",
		args:             args{s: "let p: any = [1,2,3]; let q: any = p[0..]; p == q"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         true,
		execWantErr:      false,
	}, {
		name:             "37b",
		args:             args{s: "let p: any = [1,2,3]; let q: any = [1,2,3]; p == q"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         false,
		execWantErr:      false,
	}, {
		name:             "41",
		args:             args{s: "a%aa"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestParse2(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestParse3(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestRandomObjects(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "7",
		args:             args{s: "a=aa=12345;{aa}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(12345),
		execWantErr:      false,
	}, {
		name:             "7",
		args:             args{s: "a=aa=12345;{aa}a"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(12345),
		execWantErr:      false,
	}, {
		name:             "8",
		args:             args{s: "{az:qux}.az()"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         17.0,
		execWantErr:      false,
	}, {
		name:             "9",
		args:             args{s: "{a:qux}['a']()as string as int+19"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(36),
		execWantErr:      false,
	}, {
		name:             "10",
		args:             args{s: "{a:qux}['a']()"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         17.0,
		execWantErr:      false,
	}, {
		name:             "11",
		args:             args{s: "{az:qux}[''+'az']()"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         17.0,
		execWantErr:      false,
	}, {
		name:             "15",
		args:             args{s: "{}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         map[string]interface{}{},
		execWantErr:      false,
	}, {
		name:             "15",
		args:             args{s: "{z1:1,z2:2}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         map[string]interface{}{"z1": int64(1), "z2": int64(2)},
		execWantErr:      false,
	}, {
		name:             "15",
		args:             args{s: "hashmap!{z1=>1,z2=>2}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         map[string]interface{}{"z1": int64(1), "z2": int64(2)},
		execWantErr:      false,
	}, {
		name:             "20",
		args:             args{s: "let p = {az: 3}; p['az'] = 5; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         map[string]interface{}{"az": int64(5)},
		execWantErr:      false,
	}, {
		name:             "20",
		args:             args{s: "let p = {az: 3}; p.az = 5; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         map[string]interface{}{"az": int64(5)},
		execWantErr:      false,
	}, {
		name:             "37c",
		args:             args{s: "let p: any = {ax: 1}; let q: any = p; p == q"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         true,
		execWantErr:      false,
	}, {
		name:             "37c",
		args:             args{s: "let p: any = {ax: 1}; let q: any = {ax: 1}; p == q"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         false,
		execWantErr:      false,
		// }, {
		// 	name:             "42",
		// 	args:             args{s: "{a:1}.a as int/0"},
		// 	checkCompileWant: false,
		// 	compileWantErr:   false,
		// 	execWantErr:      true,
	}, {
		name:             "43",
		args:             args{s: "let a = {let b = 3; b + 5}; a"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestRandomCalls(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "6",
		args:             args{s: "sum([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         26.0,
		execWantErr:      false,
	}, {
		name:             "6",
		args:             args{s: "sum(u8![1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         26.0,
		execWantErr:      false,
	}, {
		name:             "6",
		args:             args{s: "let p = sum; p([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         26.0,
		execWantErr:      false,
	}, {
		name:             "6",
		args:             args{s: "let p = sum; p([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         26.0,
		execWantErr:      false,
	}, {
		name:             "6",
		args:             args{s: "let p:int = sum; p([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   true,
	}, {
		name:             "6",
		args:             args{s: "{az:sum}.az([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   true,
	}, {
		name:             "6",
		args:             args{s: "[sum][0]([1.0],{a:'www'},'',3_i16,5_i32,7_u32)+11"},
		checkCompileWant: false,
		compileWantErr:   true,
	}, {
		name:             "39",
		args:             args{s: "let f = |a: int, b: float|{a * b}; let p: int = 10, q: float = 20; p + q; f(p,q)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         200.0,
		execWantErr:      false,
	}, {
		name:             "25",
		args:             args{s: "|a: int, b: float|{a + b}(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "25",
		args:             args{s: "|a: int, b: float|{return a + b}(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "25",
		args:             args{s: "|a: int, b: float|->f64{return a + b}(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "25",
		args:             args{s: "|a: int, b: float|{let c: float = a + b; c + 1}(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         9.0,
		execWantErr:      false,
	}, {
		name:             "25c",
		args:             args{s: "|a: int, b: float|{a + b}(3,5)as float+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         15.0,
		execWantErr:      false,
	}, {
		name:             "26",
		args:             args{s: "|a: int, b: float| -> int {a + b}(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "26b",
		args:             args{s: "|a: int, b: float| -> int {a + b}(3,5)+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
		breakpoint:       true,
	}, {
		name:             "27",
		args:             args{s: "let f = |a: int, b: float|{a + b}; f(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "27b",
		args:             args{s: "let f = |a: int, b: float|{a + b}; f(3,5)as float+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         15.0,
		execWantErr:      false,
	}, {
		name:             "28",
		args:             args{s: "let f = |a: int, b: float| -> int {a + b}; f(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "28b",
		args:             args{s: "let f = |a: int, b: float| -> int {a + b}; f(3,5)+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
	}, {
		name:             "29",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}qwerty(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "29b",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}qwerty(3,5)as float+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         15.0,
		execWantErr:      false,
	}, {
		name:             "30",
		args:             args{s: "fn qwerty(a: int, b: float) -> int {a + b}qwerty(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "30",
		args:             args{s: "fn qwerty(a: int, b: float) -> int {return a + b}qwerty(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "30b",
		args:             args{s: "fn qwerty(a: int, b: float) -> int {a + b}qwerty(3,5)+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
	}, {
		name:             "31",
		args:             args{s: "fn qwerty(a: int, b: float){a + b} let f = qwerty; f(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "31b",
		args:             args{s: "fn qwerty(a: int, b: float){a + b} let f = qwerty; f(3,5)as float+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         15.0,
		execWantErr:      false,
	}, {
		name:             "32",
		args:             args{s: "fn qwerty(a: int, b: float) -> int {a + b} let f = qwerty; f(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "32b",
		args:             args{s: "fn qwerty(a: int, b: float) -> int {a + b} let f = qwerty; f(3,5)+7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(15),
		execWantErr:      false,
	}, {
		name:             "33a",
		args:             args{s: "{az:|a: int, b: float|{a + b}}.az(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "33b",
		args:             args{s: "let f = |a: int, b: float|{a + b};{az:f}.az(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "33c",
		args:             args{s: "[|a: int, b: float|{a + b}][0](3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "33d",
		args:             args{s: "let f = |a: int, b: float|{a + b};[f][0](3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "33e",
		args:             args{s: "let f = |a: int, b: float|{a + b};let g = f;[g][0](3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "34a",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}{az:qwerty}.az(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "34b",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}let f = qwerty;{az:f}.az(3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "34c",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}[qwerty][0](3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "34d",
		args:             args{s: "fn qwerty(a: int, b: float){a + b}let f = qwerty;[f][0](3,5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         8.0,
		execWantErr:      false,
	}, {
		name:             "35a",
		args:             args{s: "fn add(a: float, b: float)->float{a + b} add(1,2) |> (add)(3) |> add(5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         11.0,
		execWantErr:      false,
	}, {
		name:             "35a",
		args:             args{s: "fn add(a: float, b: float)->float{a + b} add(1,2) |> |a: f64, b: f64| -> f64 {a + b}(3) |> add(5)"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         11.0,
		execWantErr:      false,
	}, {
		name:             "35b",
		args:             args{s: "fn add(a: float, b: float)->float{a + b} (add(1,2) |> (add)(3) |> add(5)) + 7"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         18.0,
		execWantErr:      false,
	}, {
		name:             "35c",
		args:             args{s: "fn add(a: float, b: float)->float{a + b} let p = add(1,2) |> (add)(3) |> add(5), q = 7; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         11.0,
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestRandomConditions(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x:int = 1; x {11} else if let y: int = 0; y {13} else {17}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(11),
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0; x {11} else if let y: int = 1; y {13} else {17}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(13),
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0; x {11} else if let y: int = 0; y {13} else {17}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(17),
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0; x {11} else {17}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(17),
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 1; x {} else {17}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0; x {11} else {}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0; x {11} else if let y: int = 0; y {13}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "21",
		args:             args{s: "if let x: int = 0, y: int = 0; x {11}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "24a",
		args:             args{s: "let p: int = (if let x: int = 0, y: int = 0; x {11} else {13})as float+ 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(30),
		execWantErr:      false,
	}, {
		name:             "24b",
		args:             args{s: "let p: int = (if let x: int = 1, y: int = 0; x {11})as float+ 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(28),
		execWantErr:      false,
	}, {
		name: "24c",
		// NOTE: runtime error (ok)
		//       nil -> float conversion is not allowed
		args:             args{s: "let p: int = (if let x: int = 0, y: int = 0; x {11})as float+ 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      true,
	}, {
		name: "24d",
		// TODO: BUG: panic (mnem.Convdyn_? family can't set type correctly if source is imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef)
		args:             args{s: "let p: int = (if let x: int = 1, y: int = 0; x {11})as string+ 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(1117),
		execWantErr:      false,
	}, {
		name: "24e",
		// TODO: BUG: panic (mnem.Convdyn_? family can't set type correctly if source is imm_?32, imm_?16, imm_?8, Imm_ptr, Imm_data, Imm_nil, Imm_ndef)
		args:             args{s: "let p: int = (if let x: int = 0, y: int = 0; x {11})as string+ 17; p"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestRandomLoops(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "while let mut x: int = 1, y: int = 0; x, x % 10 {x;x++}\nwhile let mut x: int = 1, y: int = 0; x, x % 10 {x;x++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "while let mut x: int = 1, y: int = 0; x, x % 10 {x;x++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: while let mut x: int = 1, y: int = 0; x, x % 10 {x++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "do let mut x: int = 1, y: int = 0 {x++} while x % 10;\ndo let mut x: int = 1, y: int = 0 {x++} while x % 10"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "do let mut x: int = 1, y: int = 0 {x++} while x % 10;"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: do let mut x: int = 1, y: int = 0 {x++} while x % 10;"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(9),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "for let mut i: int = 0, j: int = 0; i % 10; i++ {j++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "for let mut i: int = 1, j: int = 0; i % 10; i++ {j++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
		times:            1,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: for let mut i: int = 1, j: int = 0; i % 10; i++ {j++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: for let mut i: int = 1, j: int = 0; i % 10; i++ {j++;break 'mylabel}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: for let mut i: int = 1, j: int = 0; i % 10; i++ {j++;continue}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "'mylabel: for let mut i: int = 1, j: int = 0; i % 10; i++ {j++;continue 'mylabel}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "22",
		args:             args{s: "let mut i: int = 7, j: int = 0; for i = 1; i % 10; i++ {j++}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "loop let mut i: int = 1; {i++; if i % 5 {break}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "while let mut i: int = 1; true {i++; if i % 5 {break}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "'fooooo: loop {loop let mut i: int = 1; {i++; if i % 5 {break 'fooooo}}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "'fooooo: while 1 {while let mut i: int = 1; true {i++; if i % 5 {break 'fooooo}}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         nil,
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "loop let mut i: int = 1; {i++; if i % 5 {0} else {break returning i}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(5),
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "while let mut i: int = 1; true {i++; if i % 5 {0} else {break returning i}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(5),
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "'fooooo: loop {loop let mut i: int = 1; {i++; if i % 5 {0} else {break 'fooooo returning i}}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(5),
		execWantErr:      false,
	}, {
		name:             "23",
		args:             args{s: "'fooooo: while 1 {while let mut i: int = 1; true {i++; if i % 5 {0} else {break 'fooooo returning i}}}"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(5),
		execWantErr:      false,
	}, {
		name:             "40",
		args:             args{s: "let mut z = 0; 'mylabel: for n in 1..6 {z = z + n} z"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(21),
		execWantErr:      false,
	}, {
		name:             "40b",
		args:             args{s: "let mut z = 0u64; 'mylabel: for n in u64![1,2,3,4,5,6] {z = z + n as uint} z"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         uint64(21),
		execWantErr:      false,
	}, {
		name:             "40c",
		args:             args{s: "let mut z = 0; 'mylabel: for n in u64![1,2,3,4,5,6] {z = z + n as int} z"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(21),
		execWantErr:      false,
	}, {
		name:             "40d",
		args:             args{s: "let mut z = 0; 'mylabel: for n in u64![1,2,3,4,5,6] {z = z + n} z"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(21),
		execWantErr:      false,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}

func TestRandomFuncs(t *testing.T) {
	tests := []testMatrixItem{{
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}, {
		name: "38a",
		args: args{s: `
			fn tarai(x: i64, y: i64, z: i64) -> i64 {
				if x <= y {
					y
				} else {
					recurse(
						recurse(x - 1, y, z),
						recurse(y - 1, z, x),
						recurse(z - 1, x, y),
					)
				} as i64
			}
			tarai(7, 6, 0)
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(7),
		execWantErr:      false,
	}, {
		name: "38b",
		args: args{s: `
			fn tarai(x: i64, y: i64, z: i64) -> i64 {
				if x <= y {
					y
				} else {
					tarai(
						tarai(x - 1, y, z),
						tarai(y - 1, z, x),
						tarai(z - 1, x, y),
					)
				} as i64
			}
			tarai(7, 6, 0)
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(7),
		execWantErr:      false,
	}, {
		name: "38c",
		args: args{s: `
			fn add(x: i64, y: i64) -> i64 {
				x + y
			}

			fn zzz(f: any, x: i64, y: i64) -> i64 {
				f(x, y) as i64
			}

			zzz(add, 3, 5)
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name: "38d",
		args: args{s: `
			fn zzz(x: i64, y: i64) -> i64 {
				add(x, y) // NOTE: backref
			}

			fn add(x: i64, y: i64) -> i64 {
				x + y
			}

			zzz(3, 5)
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
	}, {
		name: "38e",
		args: args{s: `
			fn concat(x: String, y: String) -> String {
				x + y
			}

			concat('foo', 'bar')
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         "foobar",
		execWantErr:      false,
	}, {
		name: "38fa",
		args: args{s: `
		let mut count = 0;
		fn det_mandelbrot(re: f64, im: f64, limit: usize) -> any {
			// z_(n+1) = z_n**2 + c
			//         = (a + bi)**2 + c
			//         = a**2 - b**2 + 2abi + c
		
			count++;
			let mut z_re = 0.0, z_im = 0.0;
		
			for let mut i: u64 = 0.0; i <= limit; i++ {
				let w_re = (z_re * z_re - z_im * z_im) + re;
				let w_im = (z_re * z_im * 2) + im;
				z_re = w_re;
				z_im = w_im;
		
				// square of norm (re**2 + im**2)
				if z_re * z_re + z_im * z_im > 4.0 {
					return Some(i);
				}
			}
			None
		}
		
		fn plot_mandelbrot() -> any {
			let pixcel_w = 20, pixcel_h = 15;
			let nw_re = -1.2, nw_im = 0.35; // top-left     (north-west)
			let se_re = -1.0, se_im = 0.20; // bottom-right (south-east)
			let width = se_re - nw_re, height = nw_im - se_im;
		
			let mut re = 0.0, im = 0.0;
			let mut v: any = None;
			let buf = [0_u8; pixcel_w * pixcel_h];
		
			for let mut x = 0; x < pixcel_w; x++ {
				for let mut y = 0; y < pixcel_h; y++ {
					re = nw_re + x as f64 * width  / pixcel_w as f64;
					im = nw_im - y as f64 * height / pixcel_h as f64;
		
					v = det_mandelbrot(re, im, 32);
					if v !== None {
						buf[pixcel_w * y + x] = (32 - v as usize) as u8;
					} else {
						buf[pixcel_w * y + x] = 0_u8;
					}
				}
			}
		
			buf
		}
		plot_mandelbrot();
		count
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(300),
		execWantErr:      false,
	}, {
		name: "38fb",
		args: args{s: `
		let mut count = 0;
		fn det_mandelbrot(re: f64, im: f64, limit: usize) -> any {
			// z_(n+1) = z_n**2 + c
			//         = (a + bi)**2 + c
			//         = a**2 - b**2 + 2abi + c
		
			count++;
			let mut z_re = 0.0, z_im = 0.0;
		
			for i in 0..limit {
				let w_re = (z_re * z_re - z_im * z_im) + re;
				let w_im = (z_re * z_im * 2) + im;
				z_re = w_re;
				z_im = w_im;
		
				// square of norm (re**2 + im**2)
				if z_re * z_re + z_im * z_im > 4.0 {
					return Some(i);
				}
			}
			None
		}
		
		fn plot_mandelbrot() -> any {
			let pixcel_w = 20, pixcel_h = 15;
			let nw_re = -1.2, nw_im = 0.35; // top-left     (north-west)
			let se_re = -1.0, se_im = 0.20; // bottom-right (south-east)
			let width = se_re - nw_re, height = nw_im - se_im;
		
			let mut re = 0.0, im = 0.0;
			let mut v: any = None;
			let buf = u8![0_u8; pixcel_w * pixcel_h];
		
			for x in 0..(pixcel_w - 1) {
				for y in 0..(pixcel_h - 1) {
					re = nw_re + x as f64 * width  / pixcel_w as f64;
					im = nw_im - y as f64 * height / pixcel_h as f64;
		
					v = det_mandelbrot(re, im, 32);
					if v !== None {
						buf[pixcel_w * y + x] = (32 - v as usize) as u8;
					} else {
						buf[pixcel_w * y + x] = 0_u8;
					}
				}
			}
			buf
		}
		plot_mandelbrot();
		count
		`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(300),
		execWantErr:      false,
	}, {
		name: "38cccc0",
		args: args{s: `
				//{
					let qwerty = 1;
					fn add(x: i64, y: i64) -> i64 {
						x + y
					}
					fn zzz(x: i64, y: i64) -> i64 {
						add(x, y)
					}
					zzz(3, 5)
				//}
			`},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(8),
		execWantErr:      false,
		breakpoint:       false,
		// }, {
		// 	name: "38cccc1", // BUG: symbol `add` is not defined in the module scope, so it is invisible.
		// 	args: args{s: `
		// 		{
		// 			let qwerty = 1;
		// 			fn add(x: i64, y: i64) -> i64 {
		// 				x + y
		// 			}
		// 			fn zzz(x: i64, y: i64) -> i64 {
		// 				add(x, y)
		// 			}
		// 			zzz(3, 5)
		// 		}
		// 	`},
		// 	checkCompileWant: false,
		// 	compileWantErr:   false,
		// 	execWant:         int64(8),
		// 	execWantErr:      false,
		// 	breakpoint:       true,
		// }, {
		// 	name: "38cccc2", // BUG: symbol `add` in zzz() should be resolved at compile time
		// 	args: args{s: `
		// 		{
		// 			let qwerty = 1;
		// 			fn add(x: i64, y: i64) -> i64 {
		// 				x + y
		// 			}
		// 			fn zzz(x: i64, y: i64) -> i64 {
		// 				add(x, y) as i64
		// 			}
		// 			zzz(3, 5)
		// 		}
		// 	`},
		// 	checkCompileWant: false,
		// 	compileWantErr:   false,
		// 	execWant:         int64(8),
		// 	execWantErr:      false,
		// 	breakpoint:       true,
	}, {
		name:             "1",
		args:             args{s: "0"},
		checkCompileWant: false,
		compileWantErr:   false,
		execWant:         int64(0),
		execWantErr:      false,
	}}

	runMatrix(t, tests)
}
