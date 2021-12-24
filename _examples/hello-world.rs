// Hello, World!

println("Hello,", "World!");


let localVar1 = 11;
pub let publicVar1 = 13;

fn qwerty(a: i64, b: i64) -> i64 {
    let x: i64 = 13.0 as i64;
    x + a*2 + b*2;
    let f = |a: f64, b: f64| -> f64 {a + b}; // Lambda function
    'label: for let mut i: i64 = 1, j: i64 = 1; i % 100; i++ {
        x + a*2 + b*2 + f(a, b,) as i64;
        if false {
            break returning 9999
        }
        j++
    } as i64
    +
    'label2: loop {
        for let mut i: i64 = 0; i < 100; i++ {
            if 17 <= i {
                break 'label2 returning i * 2
            }
        }
    } as i64 // Use blocks as an expression.
}

fn add(x: i64, y: i64) -> i64 {
    x + y
}

fn higherOrderFunc(f: |i64,i64|->i64, x: i64, y: i64) -> i64 {
    f(x, y) as i64 // BUG: Currently, cast is required.
}

fn asdfg() -> Result<usize, String> {
    let mut ret: Result<usize, String> = "";
    ret = 3_usize;
    ret
}


println(json_stringify([{a:1,b:3,c:"5"},"7",9,[localVar1 as String,publicVar1]]));
println(json_parse("[1,2,3]"));
println(asdfg());

json_stringify(
    vec![
        hashmap!{ // Currently, non-string keys are not supported.
            foo:'12345' as i64,              // JavaScript style
            Bar=>12345 as String,
            baz: higherOrderFunc(add, 3, 5),
            Qux=> 5 |> add(7) |> add(11),    // Pipeline
            quux :7*11,
            Corge =>qwerty(3,5),
            [add(11, 13)] : qwerty(3,5),     // Dynamic property
            [17] => '',                      // Dynamic property
        },
        {}, // JavaScript style
        [], // JavaScript style
    ],
)
