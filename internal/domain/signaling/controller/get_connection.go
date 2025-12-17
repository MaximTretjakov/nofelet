package controller

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"nofelet/config"
	"nofelet/internal/domain/signaling/controller/view"
	"nofelet/pkg/singleton"
)

const (
	p2pLimitConnections = 1
	ice                 = "ice-candidate"
)

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

		if data.Type == ice {
			data.SDP = ""
		}

		if brErr := cm.Broadcast(data, conn, logger); brErr != nil {
			logger.Error("broadcast", slog.Any("err", brErr))
		}

		if config.Current().Debug {
			printSocketData(data, logger, conn)
		}
	}
}

func printSocketData(data view.SDPData, logger *slog.Logger, conn *websocket.Conn) {
	fmt.Println()
	message := fmt.Sprintf("from=%s | data=%+v\n",
		conn.RemoteAddr().String(),
		data,
	)
	logger.Info("wss", slog.String(":", message))
}
