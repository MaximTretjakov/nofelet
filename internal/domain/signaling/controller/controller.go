package controller

import (
	"log/slog"
)

type Controller struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) *Controller {
	return &Controller{
		Logger: logger,
	}
}
