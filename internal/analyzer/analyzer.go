package analyzer

import (
	"image"
	"image/color"
	"io"
)

// FruitAnalyzer define a interface para análise de frutas
type FruitAnalyzer interface {
	Analyze(imageReader io.Reader) (string, error)
}

// BaseAnalyzer implementa funcionalidades comuns para todos os analisadores
type BaseAnalyzer struct {
	getAverageColor func(image.Image) color.Color
}

// NewBaseAnalyzer cria uma nova instância do analisador base
func NewBaseAnalyzer() *BaseAnalyzer {
	return &BaseAnalyzer{
		getAverageColor: getAverageColor,
	}
}

// getAverageColor calcula a cor média de uma imagem
func getAverageColor(img image.Image) color.Color {
	bounds := img.Bounds()
	var r, g, b uint32
	pixelCount := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.At(x, y)
			pr, pg, pb, _ := pixel.RGBA()
			r += pr
			g += pg
			b += pb
			pixelCount++
		}
	}

	if pixelCount == 0 {
		return color.RGBA{0, 0, 0, 255}
	}

	return color.RGBA{
		R: uint8(r / uint32(pixelCount) >> 8),
		G: uint8(g / uint32(pixelCount) >> 8),
		B: uint8(b / uint32(pixelCount) >> 8),
		A: 255,
	}
}
