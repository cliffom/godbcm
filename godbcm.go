package godbcm

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Connection struct {
	ID uuid.UUID
}

type ConnectionManager struct {
	maxConnections int
	Connections    map[uuid.UUID]*Connection
}

func (connMgr *ConnectionManager) GetConnection() (*Connection, error) {
	if len(connMgr.Connections) >= connMgr.maxConnections {
		return nil, errors.New("no free connections")
	}

	connection := &Connection{ID: uuid.New()}
	connMgr.Connections[connection.ID] = connection
	return connection, nil
}

func (connMgr *ConnectionManager) ReleaseConnection(id uuid.UUID) error {
	if _, ok := connMgr.Connections[id]; !ok {
		return fmt.Errorf("invalid connectionID: %v", id)
	}

	delete(connMgr.Connections, id)
	return nil
}

func New(maxConnections int) *ConnectionManager {
	return &ConnectionManager{
		maxConnections: maxConnections,
		Connections:    make(map[uuid.UUID]*Connection, maxConnections),
	}
}
