package singleton

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	once     sync.Once
	instance *ConnectionManager
)

// ConnectionManager хранит коннкты клиентов
type ConnectionManager struct {
	Clients map[*websocket.Conn]bool
	Mu      sync.RWMutex
}

// NewConnectionManager возвращает единственный экземпляр ConnectionManager
func NewConnectionManager() *ConnectionManager {
	once.Do(func() {
		instance = &ConnectionManager{
			Clients: make(map[*websocket.Conn]bool),
		}
	})
	return instance
}
