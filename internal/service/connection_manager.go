package service

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionInfo struct {
	SessionCode     string
	ParticipantID   string
	ParticipantName string
}

type ConnectionManager struct {
	connections map[*websocket.Conn]*ConnectionInfo
	mutex       sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[*websocket.Conn]*ConnectionInfo),
	}
}

func (cm *ConnectionManager) RegisterConnection(conn *websocket.Conn, sessionCode, participantID, participantName string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.connections[conn] = &ConnectionInfo{
		SessionCode:     sessionCode,
		ParticipantID:   participantID,
		ParticipantName: participantName,
	}
}

func (cm *ConnectionManager) GetConnectionInfo(conn *websocket.Conn) (*ConnectionInfo, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	info, exists := cm.connections[conn]
	return info, exists
}

func (cm *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	delete(cm.connections, conn)
}

func (cm *ConnectionManager) GetConnectionsBySession(sessionCode string) []*websocket.Conn {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var connections []*websocket.Conn
	for conn, info := range cm.connections {
		if info.SessionCode == sessionCode {
			connections = append(connections, conn)
		}
	}
	return connections
}
