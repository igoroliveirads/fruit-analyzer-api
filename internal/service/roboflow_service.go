package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type RoboflowService struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

type RoboflowResponse struct {
	Predictions []struct {
		Class      string  `json:"class"`
		Confidence float64 `json:"confidence"`
	} `json:"predictions"`
}

// Modelos disponíveis para diferentes frutas
var FruitModels = map[string]string{
	"banana": "banana-ripeness-classification/5",
	// Adicione mais modelos conforme necessário
}

func NewRoboflowService(apiKey string) *RoboflowService {
	return &RoboflowService{
		apiKey:     apiKey,
		apiURL:     "https://detect.roboflow.com",
		httpClient: &http.Client{},
	}
}

func (s *RoboflowService) AnalyzeFruitRipeness(imagePath string, fruitType string) (*RoboflowResponse, error) {
	// Verificar se o tipo de fruta é suportado
	modelID, exists := FruitModels[fruitType]
	if !exists {
		return nil, fmt.Errorf("tipo de fruta não suportado: %s. tipos suportados: %v", fruitType, getSupportedFruits())
	}

	// Abrir o arquivo de imagem
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer file.Close()

	// Criar um buffer para o corpo da requisição
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Adicionar o arquivo ao corpo da requisição
	part, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("erro ao copiar arquivo: %v", err)
	}

	writer.Close()

	// Criar a requisição
	url := fmt.Sprintf("%s/%s?api_key=%s", s.apiURL, modelID, s.apiKey)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Fazer a requisição
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	// Ler a resposta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	// Decodificar a resposta
	var result RoboflowResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	return &result, nil
}

// Função auxiliar para obter a lista de frutas suportadas
func getSupportedFruits() []string {
	fruits := make([]string, 0, len(FruitModels))
	for fruit := range FruitModels {
		fruits = append(fruits, fruit)
	}
	return fruits
}
