var xoff = 0;
var yoff = 0;

function setup() {
    createCanvas(800, 800);
    pixelDensity(1);
}

function draw() {
    background(220);
    xoff = parseFloat(document.getElementById("sliderr").value); 
    yoff = parseFloat(document.getElementById("sliderr").value); 
    perlin_noise();
}

function perlin_noise() {
    loadPixels();
    let xoffStart = xoff; 
    for (var x = 0; x < width; x++) {
        var yoffTemp = yoff; 
        for (var y = 0; y < height; y++) {
            var index = (x + y * width) * 4;
            var n = noise(xoffStart, yoffTemp) * 255;
            pixels[index] = n;
            pixels[index + 1] = n;
            pixels[index + 2] = n;
            pixels[index + 3] = 255;
            yoffTemp += 0.01; 
        }
        xoffStart += 0.01; 
    }
    updatePixels();
}
