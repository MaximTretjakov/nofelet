package signaling

import (
	"nofelet/internal/dependency"
	sdpController "nofelet/internal/domain/signaling/controller"
)

func Register(deps *dependency.Container) {
	controller := sdpController.New(deps.Logger, deps.Cfg)

	sdp := deps.Signaling.Routes
	sdp.GET("/connect/:uuid", controller.GetConnection)
	sdp.GET("/turn-credentials/generate", controller.GetCoTURNCredentials)
}
