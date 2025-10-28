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
	}

	conn, err := u.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	cm.Mu.Lock()
	if len(cm.Clients) >= 2 {
		conn.Close()
	}
	cm.Clients[conn] = true
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

	for client := range cm.Clients {
		if client != sender {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(cm.Clients, client)
			}
		}
	}
}
