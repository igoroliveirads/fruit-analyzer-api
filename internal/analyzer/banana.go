package analyzer

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"

	"go.uber.org/zap"
)

// BananaAnalyzer implementa a análise específica para bananas
type BananaAnalyzer struct {
	*BaseAnalyzer
	logger *zap.Logger
}

// NewBananaAnalyzer cria uma nova instância do analisador de bananas
func NewBananaAnalyzer(logger *zap.Logger) *BananaAnalyzer {
	return &BananaAnalyzer{
		BaseAnalyzer: NewBaseAnalyzer(),
		logger:       logger,
	}
}

// RGBtoHSV converte RGB para HSV
func RGBtoHSV(r, g, b uint32) (h, s, v float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(math.Max(rf, gf), bf)
	min := math.Min(math.Min(rf, gf), bf)
	delta := max - min

	// Hue
	h = 0
	if delta == 0 {
		h = 0
	} else if max == rf {
		h = 60 * math.Mod(((gf-bf)/delta), 6)
	} else if max == gf {
		h = 60 * (((bf - rf) / delta) + 2)
	} else if max == bf {
		h = 60 * (((rf - gf) / delta) + 4)
	}

	if h < 0 {
		h += 360
	}

	// Saturation
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	// Value
	v = max

	return h, s, v
}

// getAverageHSV obtém a média HSV da imagem
func (a *BananaAnalyzer) getAverageHSV(img image.Image) (float64, float64, float64) {
	bounds := img.Bounds()
	var totalH, totalS, totalV float64
	pixelCount := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8

			h, s, v := RGBtoHSV(r, g, b)
			totalH += h
			totalS += s
			totalV += v
			pixelCount++
		}
	}

	if pixelCount == 0 {
		return 0, 0, 0
	}

	return totalH / float64(pixelCount),
		totalS / float64(pixelCount),
		totalV / float64(pixelCount)
}

// Analyze implementa a análise específica para bananas
func (a *BananaAnalyzer) Analyze(imageReader io.Reader) (string, error) {
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return "", err
	}

	// Obtém a média HSV da imagem
	h, s, v := a.getAverageHSV(img)

	// Log dos valores para debug
	a.logger.Info("Análise de cor",
		zap.Float64("hue", h),
		zap.Float64("saturation", s),
		zap.Float64("value", v),
	)

	// Análise baseada apenas no Hue (matiz)
	if h >= 50 && h <= 80 {
		return "verde", nil
	} else if h >= 35 && h <= 55 {
		return "madura", nil
	} else if h >= 15 && h <= 35 {
		return "passada", nil
	}

	return "desconhecido", nil
}
