package three_d

import (
	"image/color"
	"math"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"github.com/ojrac/opensimplex-go"
	"golang.org/x/exp/rand"
)

const (
	mapHeight = 500
	mapWidth  = 500
	title     = "3D Terrain"
)

type TerrainType struct {
	Name   string
	Height float64
	Color  color.RGBA
}

type mapGenerator struct {
	MapWidth  int
	MapHeight int

	NoiseScale  float64
	Octaves     int
	Persistence float64
	Lacunarity  float64
	Offset      math32.Vector2

	HeightMap [][]float64
	Regions   []TerrainType
	Seed      uint64
}

func (mg *mapGenerator) GenerateHeightMap() {
	mg.HeightMap = make([][]float64, mapWidth)
	for i := range mg.HeightMap {
		mg.HeightMap[i] = make([]float64, mapHeight)
	}

	if mg.NoiseScale <= 0 {
		mg.NoiseScale = 0.0001
	}

	seed := uint64(time.Now().UnixNano())
	prng := rand.New(rand.NewSource(seed))
	noise := opensimplex.New(int64(prng.Uint64()))

	octaveOffsets := make([]math32.Vector2, mg.Octaves)
	for i := 0; i < mg.Octaves; i++ {
		offsetX := float64(prng.Intn(200000)-100000) + float64(mg.Offset.X)
		offsetY := float64(prng.Intn(200000)-100000) + float64(mg.Offset.Y)
		octaveOffsets[i] = math32.Vector2{X: float32(offsetX), Y: float32(offsetY)}
	}

	maxNoiseHeight := math.Inf(-1)
	minNoiseHeight := math.Inf(1)

	halfWidth := float64(mapWidth) / 2.0
	halfHeight := float64(mapHeight) / 2.0

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			for i := 0; i < mg.Octaves; i++ {
				sampleX := (float64(x)-halfWidth)/mg.NoiseScale*frequency + float64(octaveOffsets[i].X)
				sampleY := (float64(y)-halfHeight)/mg.NoiseScale*frequency + float64(octaveOffsets[i].Y)

				perlinValue := noise.Eval2(sampleX, sampleY)*2.0 - 1.0
				noiseHeight += perlinValue * amplitude

				amplitude *= mg.Persistence
				frequency *= mg.Lacunarity
			}

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			}
			if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}
			mg.HeightMap[x][y] = noiseHeight
		}
	}

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			mg.HeightMap[x][y] = (mg.HeightMap[x][y] - minNoiseHeight) / (maxNoiseHeight - minNoiseHeight)
		}
	}
}

func GenerateMesh() *graphic.Mesh {
	geom := geometry.NewBox(1, 1, 1)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	return graphic.NewMesh(geom, mat)
}

func Run_3D_Terrain() {
	// 앱 초기 설정
	a := app.App()
	a.IWindow.(*window.GlfwWindow).SetTitle(title)
	a.IWindow.(*window.GlfwWindow).SetSize(800, 800)
	scene := core.NewNode()

	// 카메라 설정
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)
	camera.NewOrbitControl(cam)

	//scene.Add(helper.NewAxes(0.5)) //display axis

	//빛 설정
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	//백그라운드 색
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	mg := &mapGenerator{
		MapWidth:    500,
		MapHeight:   500,
		NoiseScale:  300,
		Octaves:     5,
		Persistence: 0.5,
		Lacunarity:  2.0,
		Offset:      math32.Vector2{X: 0, Y: 0},

		Seed: 0,
		Regions: []TerrainType{
			{"Mountain2", 0.1, color.RGBA{R: 74, G: 89, B: 100, A: 255}},
			{"Mountain", 0.2, color.RGBA{R: 63, G: 73, B: 75, A: 255}},
			{"Ground", 0.3, color.RGBA{R: 79, G: 49, B: 1, A: 255}},
			{"Grass2", 0.4, color.RGBA{R: 102, G: 136, B: 59, A: 255}},
			{"Grass", 0.5, color.RGBA{R: 183, G: 183, B: 71, A: 255}},
			{"Sand2", 0.6, color.RGBA{R: 183, G: 183, B: 71, A: 255}},
			{"Sand", 0.7, color.RGBA{R: 242, G: 223, B: 152, A: 255}},
			{"Water3", 0.8, color.RGBA{R: 104, G: 172, B: 214, A: 255}},
			{"Water2", 0.9, color.RGBA{R: 65, G: 125, B: 201, A: 255}},
			{"Water", 1.0, color.RGBA{R: 56, G: 103, B: 175, A: 255}},
		},
	}
	mg.GenerateHeightMap()
	//mg.generateTerrainMap_island()
	md := GenerateMesh()
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)
	scene.Add(md)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
