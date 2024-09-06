let inc = 0.01;  
let start = 0;  
let octavesSlider; 
let lacunaritySlider; 

function setup() {
  createCanvas(600, 600);

  
  octavesSlider = createSlider(1, 10, 4, 1);
  octavesSlider.position(10, height + 10);
  octavesSlider.style('width', '200px');
  lacunaritySlider = createSlider(0, 4, 0, 0.1);
  lacunaritySlider.position(10, height + 40);
  lacunaritySlider.style('width', '200px');
  persistenceSlider = createSlider(0, 1, 0.1, 0.1);
  persistenceSlider.position(10, height + 70);
  persistenceSlider.style('width', '200px');
}

function draw() {
  background(0);

  let octaves = octavesSlider.value();
  let lacunarity = lacunaritySlider.value();
  let persistence = persistenceSlider.value();
  stroke(255);
  strokeWeight(4);
  noFill();
  beginShape();

  let xoff = start;
  for (let x = 0; x < width; x++) {
    let y = 0;
    let amplitude = 1;
    let frequency = 1;

    for (let i = 0; i < octaves; i++) {
      y += noise(xoff * frequency) * amplitude;
      amplitude *= persistence;
      frequency *= lacunarity; 
    }

    
    y = map(y, 0, 1, 0, height);
    y = constrain(y, 0, height);

    vertex(x, y);

    xoff += inc;
  }
  endShape();

  start += inc;
}