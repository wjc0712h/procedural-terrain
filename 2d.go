package main

import (
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/ojrac/opensimplex-go"
	"golang.org/x/exp/rand"
)

type MapGenerator struct {
	MapWidth   int
	MapHeight  int
	NoiseScale float64
	AutoUpdate bool
	NoiseMap   [][]float64
}

func (mg *MapGenerator) generateNoiseMap() {
	mg.NoiseMap = make([][]float64, mg.MapWidth)
	for i := range mg.NoiseMap {
		mg.NoiseMap[i] = make([]float64, mg.MapHeight)
	}

	if mg.NoiseScale <= 0 {
		mg.NoiseScale = 0.0001
	}

	seed := time.Now().UnixNano()
	noise := opensimplex.New(rand.Int63n(seed))

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			sampleX := float64(x) / mg.NoiseScale
			sampleY := float64(y) / mg.NoiseScale

			perlinValue := noise.Eval2(sampleX, sampleY)
			mg.NoiseMap[x][y] = (perlinValue + 1) / 2 // Normalize to [0, 1]
		}
	}
}

func (mg *MapGenerator) drawNoiseMap(win *pixelgl.Window) {
	pic := pixel.MakePictureData((pixel.R(0, 0, float64(mg.MapWidth), float64(mg.MapHeight))))

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			noiseValue := mg.NoiseMap[x][y]
			gray := uint8(math.Floor(noiseValue * 255))
			col := color.RGBA{R: gray, G: gray, B: gray, A: 255}
			pic.Pix[y*mg.MapWidth+x] = col
		}
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	win.Clear(color.RGBA{0, 0, 0, 255})
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "2D Terrain",
		Bounds: pixel.R(0, 0, 800, 800),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	/* Init value */
	mapGen := &MapGenerator{
		MapWidth:   500,
		MapHeight:  500,
		NoiseScale: 100.0,
		AutoUpdate: false,
	}

	mapGen.generateNoiseMap()

	for !win.Closed() {
		if mapGen.AutoUpdate {
			mapGen.generateNoiseMap()
		}
		mapGen.drawNoiseMap(win)
		win.Update()
	}
}
func Run_2D_Terrain() {
	pixelgl.Run(run)
}
