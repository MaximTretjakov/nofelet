package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"nofelet/internal/domain/signaling/controller/view"
	"nofelet/pkg/singleton"
)

// GetConnection /connect/:uuid установка sdp сессии
func (c *Controller) GetConnection(ctx *gin.Context) {
	cm := singleton.NewConnectionManager()

	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
	}

	conn, err := u.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	conn.SetReadLimit(8192)

	cm.Mu.Lock()
	cm.Clients[conn] = ctx.Param("uuid")
	cm.Mu.Unlock()

	go handleClient(conn, cm)
}

func handleClient(conn *websocket.Conn, cm *singleton.ConnectionManager) {
	defer func() {
		cm.Mu.Lock()
		delete(cm.Clients, conn)
		cm.Mu.Unlock()
		conn.Close()
	}()

	var message view.Message
	for {
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}
		broadcast(message, conn, cm)
	}
}

func broadcast(message view.Message, sender *websocket.Conn, cm *singleton.ConnectionManager) {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	for client, id := range cm.Clients {
		if id == cm.Clients[sender] {
			if client != sender {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(cm.Clients, client)
				}
			}
		}
	}
}
