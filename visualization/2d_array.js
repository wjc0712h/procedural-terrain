var circles = [];
var circleFadeSpeed = 1; 
var gridFadeSpeed = 5; 
var maxAlpha = 255; 

var gridAlpha = 0; 
var circleFadeInterval = 30; 

function setup() {
  createCanvas(600, 600);

  for (var a = 10; a < width; a += 20) {
    for (var b = 10; b < height; b += 20) {
      circles.push(new Circle(a, b));
    }
  }
  console.log(circles.length);
}

function draw() {
  background(169, 169, 169);

  drawGrid();

  for (var i = 0; i < circles.length; i++) {
    circles[i].update();
    circles[i].show();
  }
}

function drawGrid() {

  stroke(0, 0, 0, gridAlpha);
  strokeWeight(1);

  for (var x = 0; x <= width; x += 20) {
    line(x, 0, x, height);
  }

  for (var y = 0; y <= height; y += 20) {
    line(0, y, width, y);
  }

  if (gridAlpha < maxAlpha) {
    gridAlpha += gridFadeSpeed;
    gridAlpha = constrain(gridAlpha, 0, maxAlpha);
  }
}

function Circle(x, y) {
  this.x = x;
  this.y = y;
  this.alpha = 0; 
  this.visible = false;

  this.update = function() {
      this.visible = true;
      if (this.alpha < maxAlpha) {
        this.alpha += circleFadeSpeed;
        this.alpha = constrain(this.alpha, 0, maxAlpha); 
      }
  }

  this.show = function() {
    if (this.visible) {
      fill(169, 169, 169, this.alpha); 
      stroke(0, 0, 0, this.alpha);
      ellipse(this.x, this.y, 10, 10);
    }
  }
}
