
let localVar1 = 1;
pub let publicVar1 = 3;

fn qwerty(a: i64, b: i64) -> i64 {
    let x: i64 = 13.0 as i64;
    x + a*2 + b*2;
    let f = |a: f64, b: f64| -> f64 {a + b};
    'label: for let mut i: i64 = 1, j: i64 = 1; i % 100; i++ {
        x + a*2 + b*2 + f(a, b,) as i64;
        j++
    }
}

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

fn add(x: i64, y: i64) -> i64 {
    x + y
}

fn zzz(f: |i64,i64|->i64, x: i64, y: i64) -> i64 {
    f(x, y) as i64
}

fn zzzz() -> Result<usize, String> {
    let mut ret: Result<usize, String> = "";
    ret = 3_usize;
    ret
}

// References:
//   https://www.suzu6.net/posts/290-rust-mandelbrot/
//   https://github.com/ProgrammingRust/mandelbrot
fn det_mandelbrot(re: f64, im: f64, limit: usize) -> Option<usize> {
    // z_(n+1) = z_n**2 + c
    //         = (a + bi)**2 + c
    //         = a**2 - b**2 + 2abi + c

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

pub fn plot_mandelbrot() -> [u8] {
    let pixcel_w = 100, pixcel_h = 75;
    let nw_re = -1.2, nw_im = 0.35; // top-left     (north-west)
    let se_re = -1.0, se_im = 0.20; // bottom-right (south-east)
    let width = se_re - nw_re, height = nw_im - se_im;

    let mut re = 0.0, im = 0.0;
    let mut v: Option<usize> = None;
    let buf = u8![0_u8; (pixcel_w * pixcel_h) as usize];

    for x in 0..(pixcel_w - 1) {
        for y in 0..(pixcel_h - 1) {
            re = nw_re + x as f64 * width  / pixcel_w as f64;
            im = nw_im - y as f64 * height / pixcel_h as f64;

            //let Some(v) = det_mandelbrot(re, im, 255);
            v = det_mandelbrot(re, im, 255);
            if v !== None {
                buf[pixcel_w * y + x] = (255 - v as usize) as u8;
            } else {
                buf[pixcel_w * y + x] = 0_u8;
            }
        }
    }
    buf
}

/*
let html =
'<!DOCTYPE html>
<html>
<head></head>
<body>
<canvas id="myCanvas" width="100" height="75"></canvas>
<script>
const buf = [...];
const canvas = document.getElementById("myCanvas");
const ctx = canvas.getContext("2d");
const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
const width = imageData.width, height = imageData.height;
const pixels = imageData.data;
for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
        const offset = y * width + x;
        pixels[offset * 4 + 0] = buf[offset]; // r
        pixels[offset * 4 + 1] = buf[offset]; // g
        pixels[offset * 4 + 2] = buf[offset]; // b
        pixels[offset * 4 + 3] = 255; // a
    }
}
ctx.putImageData(imageData, 0, 0);
</script>
</body>
</html>';
*/

// println("Hello,", "World!");
// println(json_stringify([{a:1,b:3,c:"5"},"7",9]));
// println(json_parse("[1,2,3]"));

zzzz();

//println(json_stringify(
vec![hashmap!{
    foo:'qwe',
    Bar:plot_mandelbrot() as String,
    baz:zzz(add, 3, 5),
    Qux:5 |> add(7) |> add(11),
    Quux:7*11,test:qwerty(3,5)
},{baz=>tarai(7, 6, 0),Qux=>tarai2(8,6,0)}]
//))