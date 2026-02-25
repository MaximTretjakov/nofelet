package singleton

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/internal/domain/signaling/controller/view"
)

var (
	once     sync.Once
	instance *ConnectionManager
)

// ConnectionManager хранит коннкты клиентов
type ConnectionManager struct {
	Mu      sync.RWMutex
	clients map[*websocket.Conn]string
}

// NewConnectionManager возвращает единственный экземпляр ConnectionManager
func NewConnectionManager() *ConnectionManager {
	once.Do(func() {
		instance = &ConnectionManager{
			clients: make(map[*websocket.Conn]string),
		}
	})
	return instance
}

// Save сохраняет коннекшен и uuid клиента
func (cm *ConnectionManager) Save(conn *websocket.Conn, uuid string) {
	cm.Mu.Lock()
	cm.clients[conn] = uuid
	cm.Mu.Unlock()
}

// DeleteClient удаляем коннекшен клиента
func (cm *ConnectionManager) DeleteClient(conn *websocket.Conn) {
	cm.Mu.Lock()
	delete(cm.clients, conn)
	cm.Mu.Unlock()
}

// Connections возвращает количество клиентов
func (cm *ConnectionManager) Connections() int {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	return len(cm.clients)
}

// Broadcast рассылает сообщения все клиентам доя установления SDP сессии
func (cm *ConnectionManager) Broadcast(data view.SDPData, sender *websocket.Conn, logger *slog.Logger) error {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	for client, id := range cm.clients {
		if id == cm.clients[sender] {
			if client != sender {
				err := client.WriteJSON(data)
				if err != nil {
					logger.Error("Broadcast func", slog.Any("err", err))
					_ = client.Close()
					delete(cm.clients, client)
					return err
				}
			} else {
				logger.Error("Broadcast func", slog.String("err", "client connection not found"))
				_ = client.Close()
				delete(cm.clients, client)
			}
		}
	}

	return nil
}
