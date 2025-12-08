package signaling

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"

	"nofelet/config"
	"nofelet/middleware"
)

type Container struct {
	Routes *gin.Engine
	Logger *slog.Logger
	Cfg    *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	routes, err := newRoutes(logger)
	if err != nil {
		return nil, fmt.Errorf("инициализация роутера: %w", err)
	}

	return &Container{
		Routes: routes,
		Logger: logger,
		Cfg:    cfg,
	}, nil
}

func newRoutes(logger *slog.Logger) (*gin.Engine, error) {
	router := gin.New()
	router.ContextWithFallback = true
	router.HandleMethodNotAllowed = true
	router.Use(
		gin.Recovery(),
		middleware.DurationLoggerMiddleware(),
	)
	return router, nil
}
