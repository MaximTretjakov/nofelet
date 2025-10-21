package signaling

import (
	"nofelet/internal/dependency"
	sdpController "nofelet/internal/domain/signaling/controller"
)

func Register(deps *dependency.Container) {
	controller := sdpController.New()

	sdp := deps.Signaling.Routes
	sdp.GET("/connect/:roomId", controller.GetConnection)
}
