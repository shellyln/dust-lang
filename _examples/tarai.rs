// Tarai function

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

fn tarai2(x: i64, y: i64, z: i64) -> i64 {
    if x <= y {
        y
    } else {
        tarai2(
            tarai2(x - 1, y, z),
            tarai2(y - 1, z, x),
            tarai2(z - 1, x, y),
        )
    } as i64
}

vec![hashmap!{
    tarai => tarai(7, 6, 0),
}, {
    tarai => tarai2(8, 6, 0),
}, {
    tarai => tarai(9, 6, 0),
}, {
    tarai => tarai2(10, 6, 0),
}, {
    tarai => tarai(11, 6, 0),
}, {
    tarai => tarai(12, 6, 0),
}]
