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
	"github.com/g3n/engine/util/helper"
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

type MapGenerator struct {
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

func (mg *MapGenerator) GenerateHeightMap() {
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
		offsetX := float64(mg.Offset.X)
		offsetY := float64(mg.Offset.Y)
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
func (mg *MapGenerator) GenerateHeightMap_island() {
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
		offsetX := float64(mg.Offset.X)
		offsetY := float64(mg.Offset.Y)
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

			distanceFromCenter := math.Sqrt(math.Pow(float64(x)-halfWidth, 2) + math.Pow(float64(y)-halfHeight, 2))
			gradient := 1 - (distanceFromCenter / math.Max(halfWidth, halfHeight))
			gradient = math.Max(gradient, 0)

			noiseHeight *= gradient

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			}
			if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}
			mg.HeightMap[x][y] = -1 * noiseHeight
		}
	}

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			mg.HeightMap[x][y] = (mg.HeightMap[x][y] - minNoiseHeight) / (maxNoiseHeight - minNoiseHeight)
		}
	}
}

func GenerateTerrainMesh(mg *MapGenerator) *graphic.Mesh {

	geom := geometry.NewGeometry()

	vertices := math32.NewArrayF32(0, 3*mg.MapWidth*mg.MapHeight)
	indices := math32.NewArrayU32(0, 6*(mg.MapWidth-1)*(mg.MapHeight-1))
	colors := math32.NewArrayF32(0, 3*mg.MapWidth*mg.MapHeight)

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {

			height := mg.HeightMap[x][y]
			posX := float32(x)
			posY := float32(height * 100) // 높이 스케일
			posZ := float32(y)

			vertices.Append(posX, posY, posZ)

			var color math32.Color
			// with:
			for _, region := range mg.Regions {
				if height <= region.Height {
					color = math32.Color{
						R: float32(region.Color.R) / 255,
						G: float32(region.Color.G) / 255,
						B: float32(region.Color.B) / 255,
					}
					break
				}
			}
			colors.Append(color.R, color.G, color.B)
		}
	}

	for y := 0; y < mg.MapHeight-1; y++ {
		for x := 0; x < mg.MapWidth-1; x++ {
			topLeft := uint32(y*mg.MapWidth + x)
			topRight := topLeft + 1
			bottomLeft := topLeft + uint32(mg.MapWidth)
			bottomRight := bottomLeft + 1

			indices.Append(topLeft, bottomLeft, topRight)
			indices.Append(topRight, bottomLeft, bottomRight)
		}
	}
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mat.SetShininess(30.0) // Increase shininess for specular highlights

	normals := math32.NewArrayF32(len(vertices), cap(vertices))
	normals = geometry.CalculateNormals(indices, vertices, normals)

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	//geom.AddVBO(gls.NewVBO(colors).AddAttrib(gls.VertexColor))
	geom.AddVBO(gls.NewVBO(normals).AddAttrib(gls.VertexNormal))

	mesh := graphic.NewMesh(geom, mat)
	return mesh
}

func Run_3D_Terrain() {
	// 앱 초기 설정
	a := app.App()
	a.IWindow.(*window.GlfwWindow).SetTitle(title)
	a.IWindow.(*window.GlfwWindow).SetSize(800, 800)
	scene := core.NewNode()
	scene.Add(helper.NewAxes(1000)) //display axis

	//백그라운드 색
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	mg := &MapGenerator{
		MapWidth:    500,
		MapHeight:   500,
		NoiseScale:  100,
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
	//mg.GenerateHeightMap()
	mg.GenerateHeightMap_island()
	terrainMesh := GenerateTerrainMesh(mg)
	scene.Add(terrainMesh)

	//빛 설정
	scene.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 1))
	pointLight := light.NewDirectional(&math32.Color{R: 255, G: 255, B: 255}, 100.0)
	pointLight.SetPosition(float32(mg.MapWidth)/2, 400, float32(mg.MapHeight)/2)
	scene.Add(pointLight)

	// 카메라 설정
	cam := camera.New(1)
	cam.SetFar(3000)
	cam.SetPosition(600, 400, 600)
	cam.LookAt(&math32.Vector3{X: float32(mg.MapWidth) / 2, Y: 50, Z: float32(mg.MapHeight) / 2}, &math32.Vector3{X: 0, Y: 1, Z: 0})
	camera.NewOrbitControl(cam).SetTarget(math32.Vector3{X: float32(mg.MapWidth) / 2, Y: 50, Z: float32(mg.MapHeight) / 2})
	scene.Add(cam)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
