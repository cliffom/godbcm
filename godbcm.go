package godbcm

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Connection is a representation of a DB connection
type Connection struct {
	ID uuid.UUID
}

// ConnectionManager is used to allocate and deallocate Connections
type ConnectionManager struct {
	maxConnections int
	Connections    map[uuid.UUID]*Connection
}

// GetConnection returns a Connection if there is space in the pool
func (connMgr *ConnectionManager) GetConnection() (*Connection, error) {
	if len(connMgr.Connections) >= connMgr.maxConnections {
		return nil, errors.New("no free connections")
	}

	connection := &Connection{ID: uuid.New()}
	connMgr.Connections[connection.ID] = connection
	return connection, nil
}

// ReleaseConnection returns allocation to the pool
func (connMgr *ConnectionManager) ReleaseConnection(id uuid.UUID) error {
	if _, ok := connMgr.Connections[id]; !ok {
		return fmt.Errorf("invalid connectionID: %v", id)
	}

	delete(connMgr.Connections, id)
	return nil
}

// New returns a new ConnectionManager with a defined pool size
func New(maxConnections int) *ConnectionManager {
	return &ConnectionManager{
		maxConnections: maxConnections,
		Connections:    make(map[uuid.UUID]*Connection, maxConnections),
	}
}
