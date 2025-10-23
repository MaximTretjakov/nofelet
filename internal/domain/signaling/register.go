package signaling

import (
	"nofelet/internal/dependency"
	sdpController "nofelet/internal/domain/signaling/controller"
)

func Register(deps *dependency.Container) {
	controller := sdpController.New(deps.Logger)

	sdp := deps.Signaling.Routes
	sdp.GET("/connect", controller.GetConnection)
}
