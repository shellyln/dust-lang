# Dust - toy scripting language

<img src="https://raw.githubusercontent.com/shellyln/dust-lang/master/_assets/dust-lang-logo.svg" alt="logo" style="width:250px;" width="250">


* 👍 Syntax similar to Rust
* 👍 Loose JSON parsing
* 👍 Calling host functions

<!--  -->

* 👎 Loose and poor type system
* 👎 Super very poor performance


## 🚧 ToDo

* Traditional switch expression  
  (`switch expr {expr => expr, ..., _ => expr}`)
* Casting to complex type  
  (`a as [u8]`, `a as {x: u8, y: f64}`, `a as |u8|->u8`)
* Keeping complex type information in the expression
* If-let not-null conditional expression  
  (`if let Some(z) = some_nullable_value {...} else {...}`)
* Defining and Instantiating Structs
* Calling function from host
* Deep marshalling in host function calls
* Embedding user defined parsers (macro)  
  (e.g. `sql![select * from user]`, <code>heredoc!{&#x60;foo&#x60; x &#x60;baz&#x60;}</code>, `x!(a,b,c)`)
* Tail call optimization
* Const expr checking
* Maps with non-string keys
* Closure
* Improve error messages
* More testing!

## 🐞 Bugs

* The `break` and `return` statements do not cast the value.


## 🚀 Getting started

### Execute once

```go
import (
    "fmt"
    "log"
    scripting "github.com/shellyln/dust-lang/scripting"
)

func main() {
    defer func() {
        if err := recover(); err != nil {
            // Catch runtime panic. (e.g. Div0)
            log.Fatal(err)
        }
    }()

    xctx := scripting.NewExecutionContext()

    result, err := xctx.Execute(xctx, "3 + 5")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", result)
}
```

### Compile and execute

```go
import (
    "fmt"
    "log"
    scripting "github.com/shellyln/dust-lang/scripting"
)

func main() {
    defer func() {
        if err := recover(); err != nil {
            // Catch runtime panic. (e.g. Div0)
            log.Fatal(err)
        }
    }()

    xctx := scripting.NewExecutionContext()

    ast, err := xctx.Compile(xctx, "3 + 5")
    if err != nil {
        log.Fatal(err)
    }

    result, err := ExecuteAst(xctx, ast)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", result)
}
```

### Unmarshal

```go
import (
    "fmt"
    "log"
    scripting "github.com/shellyln/dust-lang/scripting"
)

type Foobar struct {
    XFoo string `json:"foo"`
    Bar  string
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            // Catch runtime panic. (e.g. Div0)
            log.Fatal(err)
        }
    }()

    var out []Foobar

    xctx := scripting.NewExecutionContext()

    err := scripting.Unmarshal(xctx, &out, `[{foo:'qwe',Bar:'rty'}]`)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", out)
}
```

### Set host variables and functions

```go
import (
    "fmt"
    "log"
    parser "github.com/shellyln/takenoco/base"
    scripting "github.com/shellyln/dust-lang/scripting"
    xtor "github.com/shellyln/dust-lang/scripting/executor"
)

func main() {
    defer func() {
        if err := recover(); err != nil {
            // Catch runtime panic. (e.g. Div0)
            log.Fatal(err)
        }
    }()

    xctx := scripting.NewExecutionContext(scripting.VariableInfoMap{
        "a": {
            Flags: mnem.ReturnInt | mnem.Lvalue,
            Value: int64(1),
        },
        "sum": {
            Flags: mnem.ReturnFloat | mnem.Callable,
            Value: func(x ...float64) (float64, error) {
                v := 0.0
                for _, w := range x {
                    v += w
                }
                // If an error is set, the script will terminate abnormally.
                return v, nil
            },
        },
    })

    result, err := xctx.Execute(xctx, "sum(a, 3, 5)")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", result)
}
```


## 📚 Syntax

### Comments

```rust
# hash line comment (Allow only on the first line)

/*
 * block comment
 */

// line comment
```

### Identifiers

```rust
let a0_foo_bar = 0;

let 0a = 0; // syntax error!

// raw identifier
let r#try = 0;

// A single `_` is not an identifier.
// The right-hand side is evaluated without being assigned.
let _ = 0;
```

### Variable and constant definitions

```rust
// immutable variable
let a: i64 = 3;
// immutable variable (type inference)
let b = 3_i64;

// mutable variable
let mut c: i64 = 3;
// mutable variable (type inference)
let mut d = 3_i64;

// constant
const E: i64 = 3;
// constant (type inference)
const F = 3_i64;
```

### Types

```rust
// signed integer
let a: int = 3;
let a: i64 = 3;
let a: i32 = 3;
let a: i16 = 3;
let a: i8 = 3;

let b = 3_i64;
let b = 3i64;
let b = 3_i32;
let b = 3_i16;
let b = 3_i8;

// unsigned integer
let c: uint = 3;
let c: u64 = 3;
let c: u32 = 3;
let c: u16 = 3;
let c: u8 = 3;

let d = 3_u64;
let d = 3_u32;
let d = 3_u16;
let d = 3_u8;

// float
let e: float = 3;
let e: f64 = 3;
let e: f32 = 3;

let f = 3.0;
let f = 3_f64;
let f = 3_f32;

// bool
let g: bool = true;

let h = true;

// string
let i: String = "";
let i: string = "";
let i: String = '';
let i: String = ``;

let j = "";
let j = '';
let j = ``;

// array
let k = [1, 2, 3i32];
let k = [1i32;3]; // equivalent to `[1i32, 1i32, 1i32]`
let k = vec![1, 2, 3];
let k = i64![1, 2, 3];
let k = i32![1, 2, 3];
let k = i16![1, 2, 3];
let k = i8![1, 2, 3];
let k = u64![1, 2, 3];
let k = u32![1, 2, 3];
let k = u16![1, 2, 3];
let k = u8![1, 2, 3];
let k = u8![1u8;3]; // equivalent to `u8![1u8, 2u8, 3u8]`
let k = f64![1, 2, 3];
let k = f32![1, 2, 3];
let k = bool![true, false, true];
let k = str!["1", '2', `3`];

// object
let l = {"a": 1, 'b': 2, `c`: 3, ["d" + ""]: 4, e: 5};
let l = hashmap!{"a" => 1, "b" => 2, "c" => 3, "d" => 4, "e" => 5};
let l = map!{"a" => 1, "b" => 2, "c" => 3, "d" => 4, "e" => 5};
let l = collection!{"a" => 1, "b" => 2, "c" => 3, "d" => 4, "e" => 5};

// null
let m: any = None;
let m: any = null;

// unit value
let n: any = ();
let n: any = undefined;
```

### Scope

```rust
{
    let x = 3;
    // x is defined
}
// x is not defined

// scope as an expression
let a = {let b = 3; b + 5};
```


### Control statements / expressions

#### If-else if-else

```rust
let q = doSomething();

// statement
let x = 11, y = 13;
if q < x {
    11
} else if q < y {
    13
} else {
    q / 2
}

// statement with variable definition
if let x = 11; q < x {
    11
} else if let y = 13; q < y {
    13
} else {
    q / 2
}

// expression with variable definition
let z = if let x = 11; q < x {
    11
} else if let y = 13; q < y {
    13
} else {
    q / 2
}
```


#### Infinite loop

```rust
// statement
let mut x = 0, y = 100;
loop {
    y++;
    x = ++x + y
    if x < y {
        break;
    }
}

// statement with variable definition
loop let mut x = 0, y = 100; {
    y++;
    x = ++x + y
    if x < y {
        break;
    }
}

// expression with variable definition
let z = loop let mut x = 0, y = 100; {
    y++;
    x = ++x + y
    if x < y {
        break;
    }
};
```


#### While loop

```rust
// statement
let mut x = 0, y = 100;
while x < y {
    y++;
    x = ++x + y
}

// statement with variable definition
while let mut x = 0, y = 100; x < y {
    y++;
    x = ++x + y
}

// expression with variable definition
let z = while let mut x = 0, y = 100; x < y {
    y++;
    x = ++x + y
};
```


#### Do-while loop

```rust
// statement
let mut x = 0, y = 100;
do {
    y++;
    x = ++x + y
} while x < y; // ";" is required if the statement follows.

// statement with variable definition
do let mut x = 0, y = 100 {
    y++;
    x = ++x + y
} while x < y;

// expression with variable definition
let z = do let mut x = 0, y = 100 {
    y++;
    x = ++x + y
} while x < y;
```


### For-in loop

```rust
// statement
let mut sum = 0;
for n in 0..99 {
    sum = sum + n;
}

let mut sum = 0;
for n in [3 ,5, 7, 11] {
    sum = sum + n;
}

// statement with variable definition
for let mut sum = 0; n in 0..99 {
    sum = sum + n;
}

for let mut sum = 0; n in [3 ,5, 7, 11] {
    sum = sum + n;
}

// expression with variable definition
let z = for let mut sum = 0; n in 0..99 {
    sum = sum + n;
};

let z = for let mut sum = 0; n in [3 ,5, 7, 11] {
    sum = sum + n;
};
```


### Break and continue

```rust
loop {
    break returning 3;
}

'mylabel: loop {
    break 'mylabel returning 3;
}

for x in 0..99 {
    if x == 3 {
        continue;
    }
}

'mylabel: for x in 0..99 {
    if x == 3 {
        continue 'mylabel;
    }
}
```



### Traditional for loop

```rust
// statement
for let mut i = 0, j = 0; i < 10; i++ {
    j++
}

// expression
let z = for let mut i = 0, j = 0; i < 10; i++ {
    j++
};
```


### Function

```rust
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

tarai(8, 6, 0);
```

### Lambda

```rust
let tarai = |x: i64, y: i64, z: i64| -> i64 {
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

tarai(8, 6, 0);

|a: i64, b: i64| -> i64 {a + b}(3, 5)
```


### Operators

```rust
// Member Access `op1 . op2`
hashmap!{abc => 11}.abc;

// Computed Member Access `op1[op2]`
hashmap!{abc => 11}["ab" + "c"];
[1, 3, 5, 7][0];
"abcd"[0]; // u8

// Slicing `op1[start..end]`
[1, 3, 5, 7][0..4];
[1, 3, 5, 7][..4];
[1, 3, 5, 7][0..];
"abcd"[0..4]; // String

// Type cast `op1 as T`
7 as f64;

// Function Call `op1(args)`
my_func1();
my_func2(1, 2, 3);

//-------------------------------

// Postfix Increment `op1 ++`
let mut i = 0;
i--;

// Postfix Decrement `op1 --`
i++;

//-------------------------------

// Logical NOT `! op1`
let mut a = 0;
a = !a;

// Bitwise NOT `~ op1`
a = ~a;

// Unary plus `+ op1`
a = +a;

// Unary negation `- op1`
a = -a;

// Prefix Increment `++ op1`
let mut i = 0;
++i;

// Prefix Decrement `-- op1`
--i;

//-------------------------------

// Exponentiation `op1 ** op2`
let mut a = 3.0, b = 5.0;
a = a ** b;


//-------------------------------

// Multiplication `op1 * op2`
let mut a = 3, b = 1;

// Division `op1 / op2`
a = a / b;

// Remainder `op1 % op2`
a = a % b;

//-------------------------------

// Addition `op1 + op2`
let mut a = 3, b = 5;
a = a + b;

// Subtraction `op1 - op2`
a = a - b;

//-------------------------------

// Bitwise Left Shift `op1 << op2`
let mut a = 3, b = 5;
a = a << 5;

// Bitwise Right Shift `op1 >> op2`
a = a >> 5;

// Bitwise Unsigned Right Shift `op1 >>> op2`
a = a >>> 5;

//-------------------------------

// Less Than `op1 < op2`
let mut a = 3, b = 5;
a < b;

// Less Than Or Equal `op1 <= op2`
a <= b;

// Greater Than `op1 > op2`
a > b;

// Greater Than Or Equal `op1 >= op2`
a >= b;

//-------------------------------

// Equality `op1 == op2`
let mut a = 3, b = 5;
a == b;

// Inequality `op1 != op2`
a != b;

// Strict Equality `op1 === op2`
a === b;

// Strict Inequality `op1 !== op2`
a !== b;

//-------------------------------

// Bitwise AND `op1 & op2`
let mut a = 3, b = 5;
a & b;

//-------------------------------

// Bitwise XOR `op1 ^ op2`
let mut a = 3, b = 5;
a ^ b;

//-------------------------------

// Bitwise OR `op1 | op2`
let mut a = 3, b = 5;
a | b;

//-------------------------------

// Logical AND `op1 && op2`
let mut a = true, b = false;
a && b;

//-------------------------------

// Logical OR `op1 || op2`
let mut a = true, b = false;
a || b;

//-------------------------------

// Conditional (ternary) operator `op1 ? op2 : op3`
let mut a = true;
a ? 3 : 5;

//-------------------------------

// Pipeline Function Call `op1 |> op2(args)`
fn add(a: f64, b: f64) -> f64 {
    a + b
}

let p = add(1, 2) |> add(3) |> add(5); // equivalent to `add(add(add(1,2),3),5)`

//-------------------------------

// Assignment `op1 = op2`
let mut a = 0;
a = 3;

//-------------------------------

// Comma / Sequence `op1, op2`
let mut a = 0;
a = 1 + 2, 3 + 4, 5 + 6;

++a, ++a, ++a;

```
