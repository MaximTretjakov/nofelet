package signaling

import (
	"nofelet/internal/dependency"
	"nofelet/internal/domain/signaling"
)

func New(deps *dependency.Container) error {
	signaling.Register(deps)

	return nil
}
