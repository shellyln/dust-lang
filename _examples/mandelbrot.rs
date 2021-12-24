// Mandelbrot set
//
// ./dust ./_examples/mandelbrot.rs > ./_examples/index.html && xdg-open ./_examples/index.html

const CANVAS_WIDTH  = 100;
const CANVAS_HEIGHT =  75;


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

pub fn plot_mandelbrot() -> [u16] {
    let pixcel_w = CANVAS_WIDTH, pixcel_h = CANVAS_HEIGHT;
    let nw_re = -1.2, nw_im = 0.35; // top-left     (north-west)
    let se_re = -1.0, se_im = 0.20; // bottom-right (south-east)
    let width = se_re - nw_re, height = nw_im - se_im;

    let mut re = 0.0, im = 0.0;
    let mut v: Option<usize> = None;
    let buf = u16![0_u16; (pixcel_w * pixcel_h) as usize];

    for x in 0..(pixcel_w - 1) {
        for y in 0..(pixcel_h - 1) {
            re = nw_re + x as f64 * width  / pixcel_w as f64;
            im = nw_im - y as f64 * height / pixcel_h as f64;

            //let Some(v) = det_mandelbrot(re, im, 255);
            v = det_mandelbrot(re, im, 255);
            if v !== None {
                buf[pixcel_w * y + x] = (255 - v as usize) as u16;
            } else {
                buf[pixcel_w * y + x] = 0_u16;
            }
        }
    }
    buf
}


'<!DOCTYPE html>
<html>
<head></head>
<body>
  <canvas id="myCanvas" width="' + CANVAS_WIDTH + '" height="' + CANVAS_HEIGHT + '"></canvas>
  <script>
    const buf = ' + json_stringify(plot_mandelbrot()) + ';
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
            pixels[offset * 4 + 3] = 255;         // a
        }
    }
    ctx.putImageData(imageData, 0, 0);
  </script>
</body>
</html>'
