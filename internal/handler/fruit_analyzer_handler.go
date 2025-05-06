package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fruit-analyzer-api/internal/config"
	"fruit-analyzer-api/internal/service"

	"go.uber.org/zap"
)

type FruitAnalyzerHandler struct {
	analyzerService *service.FruitAnalyzerService
	logger          *zap.Logger
}

func NewFruitAnalyzerHandler(cfg *config.Config, logger *zap.Logger) *FruitAnalyzerHandler {
	return &FruitAnalyzerHandler{
		analyzerService: service.NewFruitAnalyzerService(cfg.Roboflow.APIKey),
		logger:          logger,
	}
}

func (h *FruitAnalyzerHandler) AnalyzeFruit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obter o tipo de fruta do parâmetro da URL
	fruitType := r.URL.Query().Get("type")
	if fruitType == "" {
		http.Error(w, "Parâmetro 'type' é obrigatório", http.StatusBadRequest)
		return
	}

	// Limitar o tamanho do arquivo para 10MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Erro ao receber arquivo", zap.Error(err))
		http.Error(w, "Erro ao receber arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Criar diretório temporário se não existir
	tempDir := "temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		h.logger.Error("Erro ao criar diretório temporário", zap.Error(err))
		http.Error(w, "Erro ao criar diretório temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Criar arquivo temporário
	tempFile := filepath.Join(tempDir, handler.Filename)
	dst, err := os.Create(tempFile)
	if err != nil {
		h.logger.Error("Erro ao criar arquivo temporário", zap.Error(err))
		http.Error(w, "Erro ao criar arquivo temporário: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	defer os.Remove(tempFile)

	// Copiar conteúdo do arquivo
	if _, err := io.Copy(dst, file); err != nil {
		h.logger.Error("Erro ao salvar arquivo", zap.Error(err))
		http.Error(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Analisar a fruta
	result, err := h.analyzerService.AnalyzeFruit(tempFile, fruitType)
	if err != nil {
		h.logger.Error("Erro ao analisar fruta", zap.Error(err))
		http.Error(w, "Erro ao analisar fruta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
