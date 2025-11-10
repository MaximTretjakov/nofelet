package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NewWebSocket создает настроенный сокет
func NewWebSocket(ctx *gin.Context, logger *slog.Logger) (*websocket.Conn, error) {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
	}

	conn, err := u.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Error("websocket", slog.Any("err", err))
		return nil, err
	}
	conn.SetReadLimit(8192)

	return conn, nil
}
