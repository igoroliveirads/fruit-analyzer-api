package server

import (
	"context"
	"net/http"

	"fruit-analyzer-api/internal/handler"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server configura e gerencia o servidor HTTP
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	handler    *handler.FruitHandler
}

// NewServer cria uma nova instância do servidor
func NewServer(logger *zap.Logger, handler *handler.FruitHandler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: ":8080",
		},
		logger:  logger,
		handler: handler,
	}
}

// RegisterHooks registra os hooks do ciclo de vida do servidor
func (s *Server) RegisterHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.logger.Info("Servidor iniciado na porta 8080")
			http.HandleFunc("/analyze", s.handler.HandleAnalyze)
			go s.httpServer.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Info("Servidor sendo encerrado")
			return s.httpServer.Shutdown(ctx)
		},
	})
}
