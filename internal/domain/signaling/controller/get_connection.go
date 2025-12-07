package controller

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"nofelet/internal/domain/signaling/controller/view"
	"nofelet/pkg/singleton"
)

const p2pLimitConnections = 1

// GetConnection /connect/:uuid установка sdp сессии
func (c *Controller) GetConnection(ctx *gin.Context) {
	conn, sErr := NewWebSocket(ctx, c.Logger)
	if sErr != nil {
		c.Logger.Error("socket creation", slog.Any("err", sErr))
	}

	cm := singleton.NewConnectionManager()
	if cm.Connections() <= p2pLimitConnections {
		cm.Save(conn, ctx.Param("uuid"))
		if syncErr := cm.SyncClientsAndBroadcast(); syncErr != nil {
			c.Logger.Error("sync", slog.Any("err", syncErr))
		}
		go handleClient(conn, cm, c.Logger)
	}
}

func handleClient(conn *websocket.Conn, cm *singleton.ConnectionManager, logger *slog.Logger) {
	defer func() {
		cm.DeleteClient(conn)
		conn.Close()
	}()

	var data view.SDPData
	for {
		readErr := conn.ReadJSON(&data)
		if readErr != nil {
			logger.Error("socket read", slog.Any("err", readErr))
			break
		}
		if brErr := cm.Broadcast(data, conn, logger); brErr != nil {
			logger.Error("broadcast", slog.Any("err", brErr))
		}
	}
}
