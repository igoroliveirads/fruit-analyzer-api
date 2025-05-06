package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fruit-analyzer-api/internal/config"
	"fruit-analyzer-api/internal/handler"
	"fruit-analyzer-api/internal/server"
	"fruit-analyzer-api/internal/service"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("arquivo .env não encontrado: %v", err)
	}

	// Criar o logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Criar o contexto principal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar o fx
	app := fx.New(
		// Providers
		fx.Provide(
			func() *zap.Logger { return logger },
			config.NewConfig,
			service.NewRoboflowService,
			service.NewFruitAnalyzerService,
			handler.NewFruitAnalyzerHandler,
			server.NewServer,
		),

		// Invoke
		fx.Invoke(func(s *server.Server) error {
			return s.Start()
		}),
	)

	// Configurar o canal de sinais para graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar a aplicação
	if err := app.Start(ctx); err != nil {
		logger.Fatal("Erro ao iniciar aplicação", zap.Error(err))
	}

	// Aguardar sinal de encerramento
	<-sigChan
	logger.Info("Recebido sinal de encerramento")

	// Encerrar a aplicação
	if err := app.Stop(ctx); err != nil {
		logger.Fatal("Erro ao encerrar aplicação", zap.Error(err))
	}
}
