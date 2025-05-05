package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"fruit-analyzer-api/internal/config"
	"fruit-analyzer-api/internal/handler"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Server configura e gerencia o servidor HTTP
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	handler    *handler.FruitHandler
	config     *config.Config
}

// NewServer cria uma nova instância do servidor
func NewServer(logger *zap.Logger, handler *handler.FruitHandler, config *config.Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		logger:  logger,
		handler: handler,
		config:  config,
	}
}

// corsMiddleware adiciona headers CORS
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, origin := range s.config.AllowedOrigins {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// RegisterHooks registra os hooks do ciclo de vida do servidor
func (s *Server) RegisterHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.logger.Info(fmt.Sprintf("Servidor iniciado na porta %d", s.config.Port))
			http.HandleFunc("/analyze", s.corsMiddleware(s.handler.HandleAnalyze))
			go s.httpServer.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.logger.Info("Servidor sendo encerrado")
			return s.httpServer.Shutdown(ctx)
		},
	})
}
