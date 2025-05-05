package main

import (
	"fruit-analyzer-api/internal/analyzer"
	"fruit-analyzer-api/internal/handler"
	"fruit-analyzer-api/internal/server"
	"fruit-analyzer-api/pkg/logger"

	"go.uber.org/fx"
)

// AsAnalyzer converte BananaAnalyzer para a interface FruitAnalyzer
func AsAnalyzer(banana *analyzer.BananaAnalyzer) analyzer.FruitAnalyzer {
	return banana
}

func main() {
	app := fx.New(
		// Providers
		fx.Provide(
			logger.NewLogger,
			analyzer.NewBananaAnalyzer,
			fx.Annotate(
				AsAnalyzer,
				fx.As(new(analyzer.FruitAnalyzer)),
			),
			handler.NewFruitHandler,
			server.NewServer,
		),
		// Invoke
		fx.Invoke(func(s *server.Server, l fx.Lifecycle) {
			s.RegisterHooks(l)
		}),
	)

	app.Run()
}
