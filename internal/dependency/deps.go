package dependency

import (
	"fmt"
	"log/slog"

	"nofelet/config"
	"nofelet/internal/dependency/signaling"
)

// Container основной контейнер зависимостей
type Container struct {
	Signaling *signaling.Container
	Logger    *slog.Logger
	Cfg       *config.Config
}

func New(Cfg *config.Config, logger *slog.Logger) (*Container, error) {
	SignalingContainer, err := signaling.New(Cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("создание сигналинг контейнера: %w", err)
	}

	return &Container{
		Signaling: SignalingContainer,
		Logger:    logger,
		Cfg:       Cfg,
	}, nil
}
