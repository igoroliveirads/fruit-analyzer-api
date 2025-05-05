package logger

import (
	"go.uber.org/zap"
)

// NewLogger cria uma nova instância do logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
} 