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
	config  *config.Config
	logger  *zap.Logger
	handler *handler.FruitAnalyzerHandler
	server  *http.Server
}

// NewServer cria uma nova instância do servidor
func NewServer(
	cfg *config.Config,
	logger *zap.Logger,
	handler *handler.FruitAnalyzerHandler,
) *Server {
	return &Server{
		config:  cfg,
		logger:  logger,
		handler: handler,
	}
}

func (s *Server) Start() error {
	// Configurar o servidor HTTP
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      s.setupRoutes(),
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
	}

	// Iniciar o servidor em uma goroutine
	go func() {
		s.logger.Info("Iniciando servidor",
			zap.Int("port", s.config.Server.Port),
		)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Erro ao iniciar servidor",
				zap.Error(err),
			)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Encerrando servidor")

	// Criar um contexto com timeout para o shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.Server.ShutdownTimeout)
	defer cancel()

	// Tentar fazer o shutdown gracioso
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("Erro ao encerrar servidor graciosamente",
			zap.Error(err),
		)
		return s.server.Close()
	}

	return nil
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Middleware para logging
	mux.Handle("/analyze", s.loggingMiddleware(s.handler.AnalyzeFruit))

	return mux
}

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Criar um response writer customizado para capturar o status code
		rw := newResponseWriter(w)

		// Chamar o próximo handler
		next.ServeHTTP(rw, r)

		// Log da requisição
		s.logger.Info("Requisição processada",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

// ResponseWriter customizado para capturar o status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Provider para fx
func ProvideServer() fx.Option {
	return fx.Provide(NewServer)
}
