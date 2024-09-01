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
	Offset      pixel.Vec

	Seed       uint64
	AutoUpdate bool
	Regions    []TerrainType
	NoiseMap   [][]float64
}

func (mg *MapGenerator) GenerateNoiseMap() {
	// 노이즈맵 초기화
	mg.NoiseMap = make([][]float64, mg.MapWidth)
	for i := range mg.NoiseMap {
		mg.NoiseMap[i] = make([]float64, mg.MapHeight)
	}
	// NoiseScale 초기값 0 방지
	if mg.NoiseScale <= 0 {
		mg.NoiseScale = 0.0001
	}

	// opensimplex 이용해서 노이즈 생성
	seed := uint64(time.Now().UnixNano())
	prng := rand.New(rand.NewSource(seed))
	noise := opensimplex.New(int64(rand.Uint64()))

	octaveOffsets := make([]pixel.Vec, mg.Octaves)
	for i := 0; i < mg.Octaves; i++ {
		offsetX := float64(prng.Intn(200000)-100000) + mg.Offset.X
		offsetY := float64(prng.Intn(200000)-100000) + mg.Offset.Y
		octaveOffsets[i] = pixel.V(offsetX, offsetY)
	}

	maxNoiseHeight := math.Inf(-1)
	minNoiseHeight := math.Inf(1)

	halfWidth := float64(mg.MapWidth) / 2.0
	halfHeight := float64(mg.MapHeight) / 2.0

	//각각의 좌표를 순회하며 높이 계산.
	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0
			//각 옥타브는 최종 노이즈 높이에 기여, 진폭과 주파수에 의해 조정됨. 이 값들은 옥타브에 따라 Persistence와 Lacunarity에 따라 달라짐.
			for i := 0; i < mg.Octaves; i++ {
				sampleX := (float64(x)-halfWidth)/mg.NoiseScale*frequency + octaveOffsets[i].X
				sampleY := (float64(y)-halfHeight)/mg.NoiseScale*frequency + octaveOffsets[i].Y

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
			mg.NoiseMap[x][y] = noiseHeight
		}
	}
	//Normalization 0 ~ 1
	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			mg.NoiseMap[x][y] = (mg.NoiseMap[x][y] - minNoiseHeight) / (maxNoiseHeight - minNoiseHeight)
		}
	}
}

func (mg *MapGenerator) GenerateNoiseMap_island() {
	mg.NoiseMap = make([][]float64, mg.MapWidth)
	for i := range mg.NoiseMap {
		mg.NoiseMap[i] = make([]float64, mg.MapHeight)
	}

	if mg.NoiseScale <= 0 {
		mg.NoiseScale = 0.0001
	}

	seed := uint64(time.Now().UnixNano())
	prng := rand.New(rand.NewSource(seed))
	noise := opensimplex.New(int64(rand.Uint64()))

	octaveOffsets := make([]pixel.Vec, mg.Octaves)
	for i := 0; i < mg.Octaves; i++ {
		offsetX := float64(prng.Intn(200000)-100000) + mg.Offset.X
		offsetY := float64(prng.Intn(200000)-100000) + mg.Offset.Y
		octaveOffsets[i] = pixel.V(offsetX, offsetY)
	}

	maxNoiseHeight := math.Inf(-1)
	minNoiseHeight := math.Inf(1)

	//중앙 좌표와 최대 거리 계산
	halfWidth := float64(mg.MapWidth) / 2.0
	halfHeight := float64(mg.MapHeight) / 2.0
	maxDistance := math.Sqrt(halfWidth*halfWidth + halfHeight*halfHeight)

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			amplitude := 1.0
			frequency := 1.0
			noiseHeight := 0.0

			for i := 0; i < mg.Octaves; i++ {
				sampleX := (float64(x)-halfWidth)/mg.NoiseScale*frequency + octaveOffsets[i].X
				sampleY := (float64(y)-halfHeight)/mg.NoiseScale*frequency + octaveOffsets[i].Y

				perlinValue := noise.Eval2(sampleX, sampleY)*2.0 - 1.0
				noiseHeight += perlinValue * amplitude

				amplitude *= mg.Persistence
				frequency *= mg.Lacunarity
			}

			distance := math.Sqrt(math.Pow(float64(x)-halfWidth, 2) + math.Pow(float64(y)-halfHeight, 2))
			gradient := 1.0 - (distance / maxDistance) /* 그래디언트 계산. 중앙이 높고 멀어질수록 낮아짐 */
			noiseHeight *= gradient

			if noiseHeight > maxNoiseHeight {
				maxNoiseHeight = noiseHeight
			}
			if noiseHeight < minNoiseHeight {
				minNoiseHeight = noiseHeight
			}
			mg.NoiseMap[x][y] = noiseHeight
		}
	}

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			mg.NoiseMap[x][y] = (mg.NoiseMap[x][y] - minNoiseHeight) / (maxNoiseHeight - minNoiseHeight)
		}
	}
}

/* 화면 출력 담당 함수 */
func (mg *MapGenerator) DrawNoiseMap(win *pixelgl.Window) {
	pic := pixel.MakePictureData(pixel.R(0, 0, float64(mg.MapWidth), float64(mg.MapHeight)))

	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			noiseValue := mg.NoiseMap[x][y]
			var col color.RGBA
			for _, region := range mg.Regions {
				if noiseValue <= region.Height {
					col = region.Color
					break
				}
			}

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

	mapGen := &MapGenerator{
		MapWidth:    500,
		MapHeight:   500,
		NoiseScale:  300,
		Octaves:     5,
		Persistence: 0.5,
		Lacunarity:  2.0,
		Seed:        0, //uint64(time.Now().UnixNano())
		Offset:      pixel.V(0, 0),
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
		AutoUpdate: true,
	}

	//mapGen.GenerateNoiseMap_island()
	mapGen.GenerateNoiseMap()
	for !win.Closed() {
		if mapGen.AutoUpdate {
			mapGen.GenerateNoiseMap_island()
			//mapGen.GenerateNoiseMap()

		}
		mapGen.DrawNoiseMap(win)
		win.Update()
		time.Sleep(time.Millisecond * 1000)
	}
}

func Run_2D_Terrain() {
	pixelgl.Run(run)
}
