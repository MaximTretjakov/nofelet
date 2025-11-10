package singleton

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"

	"nofelet/internal/domain/signaling/controller/view"
)

const p2pMaxClients = 2

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

// SyncClientsAndBroadcast проверка что все кленты подключены
func (cm *ConnectionManager) SyncClientsAndBroadcast() error {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	if len(cm.clients) == p2pMaxClients {
		for caller := range cm.clients {
			wErr := caller.WriteMessage(websocket.TextMessage, []byte("clientExist"))
			return wErr
		}
	}

	return nil
}

func (cm *ConnectionManager) Broadcast(data view.SDPData, sender *websocket.Conn, logger *slog.Logger) error {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()

	for client, id := range cm.clients {
		if id == cm.clients[sender] {
			if client != sender {
				err := client.WriteJSON(data)
				if err != nil {
					logger.Error("broadcast send", slog.Any("err", err))
					client.Close()
					delete(cm.clients, client)
					return err
				}
			}
		}
	}

	return nil
}
