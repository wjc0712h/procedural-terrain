package three_d

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/texture"
)

func CreateImageFromColors(mg *MapGenerator) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, mg.MapWidth, mg.MapHeight))
	for y := 0; y < mg.MapHeight; y++ {
		for x := 0; x < mg.MapWidth; x++ {
			height := mg.HeightMap[x][y]
			var col color.RGBA

			for _, region := range mg.Regions {
				if height <= region.Height {
					//fmt.Println(height)
					col = region.Color
					break
				}
			}

			img.SetRGBA(x, y, col)
		}
	}
	return img
}

func CreateTextureFromImage(img *image.RGBA) *texture.Texture2D {
	tex := texture.NewTexture2DFromRGBA(img)
	tex.SetWrapS(gls.REPEAT)
	tex.SetWrapT(gls.REPEAT)
	return tex
}

func SaveImageToFile(img *image.RGBA, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}
