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

// Analyze implementa a análise específica para bananas
func (a *BananaAnalyzer) Analyze(imageReader io.Reader) (string, error) {
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return "", err
	}

	avgColor := a.getAverageColor(img)
	r, g, b, _ := avgColor.RGBA()

	// Normaliza os valores para 8 bits
	r = r >> 8
	g = g >> 8
	b = b >> 8

	// Converte para HSV para melhor análise de cor
	h, s, v := RGBtoHSV(r, g, b)

	// Log dos valores para debug
	a.logger.Info("Valores de cor",
		zap.Uint32("vermelho", r),
		zap.Uint32("verde", g),
		zap.Uint32("azul", b),
		zap.Float64("hue", h),
		zap.Float64("saturation", s),
		zap.Float64("value", v),
	)

	// Análise baseada em HSV
	// Verde: Hue entre 60-120, alta saturação
	// Madura: Hue entre 40-60, alta saturação
	// Passada: Hue entre 20-40, saturação média-baixa

	if h >= 60 && h <= 120 && s > 0.5 {
		return "verde", nil
	} else if h >= 40 && h <= 60 && s > 0.5 {
		return "madura", nil
	} else if h >= 20 && h <= 40 && s < 0.7 {
		return "passada", nil
	}

	// Fallback para análise RGB se HSV não for conclusivo
	if g > r && g > b && g > 150 {
		return "verde", nil
	} else if r > g && g > b && r > 150 && g > 100 {
		return "madura", nil
	} else if r > g && b > g && r > 150 && b > 100 {
		return "passada", nil
	}

	return "desconhecido", nil
}
