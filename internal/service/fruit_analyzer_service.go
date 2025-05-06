package service

import (
	"fmt"
)

type FruitAnalyzerService struct {
	roboflowService *RoboflowService
}

type AnalysisResult struct {
	FruitType     string  `json:"fruit_type"`
	RipenessLevel string  `json:"ripeness_level"`
	Confidence    float64 `json:"confidence"`
}

func NewFruitAnalyzerService(roboflowAPIKey string) *FruitAnalyzerService {
	return &FruitAnalyzerService{
		roboflowService: NewRoboflowService(roboflowAPIKey),
	}
}

func (s *FruitAnalyzerService) AnalyzeFruit(imagePath string, fruitType string) (*AnalysisResult, error) {
	// Análise do Roboflow
	roboflowResult, err := s.roboflowService.AnalyzeFruitRipeness(imagePath, fruitType)
	if err != nil {
		return nil, fmt.Errorf("erro na análise do roboflow: %v", err)
	}

	// Retornar resultado
	result := &AnalysisResult{
		FruitType:     fruitType,
		RipenessLevel: roboflowResult.Predictions[0].Class,
		Confidence:    roboflowResult.Predictions[0].Confidence,
	}

	return result, nil
}
