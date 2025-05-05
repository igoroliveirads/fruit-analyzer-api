package handler

import (
	"net/http"

	"fruit-analyzer-api/internal/analyzer"

	"go.uber.org/zap"
)

// FruitHandler gerencia as requisições relacionadas à análise de frutas
type FruitHandler struct {
	analyzer analyzer.FruitAnalyzer
	logger   *zap.Logger
}

// NewFruitHandler cria uma nova instância do handler de frutas
func NewFruitHandler(analyzer analyzer.FruitAnalyzer, logger *zap.Logger) *FruitHandler {
	return &FruitHandler{
		analyzer: analyzer,
		logger:   logger,
	}
}

// HandleAnalyze processa a requisição de análise de fruta
func (h *FruitHandler) HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Erro ao receber imagem", zap.Error(err))
		http.Error(w, "Erro ao processar imagem", http.StatusBadRequest)
		return
	}
	defer file.Close()

	status, err := h.analyzer.Analyze(file)
	if err != nil {
		h.logger.Error("Erro ao analisar fruta", zap.Error(err))
		http.Error(w, "Erro ao analisar fruta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "` + status + `"}`))
}
