package godbcm_test

import (
	"log"
	"testing"
	"time"

	"github.com/cliffom/godbcm"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	mgr := godbcm.New(1)

	if mgr == nil {
		log.Printf("expected a connection manager but got nil")
	}
}

func TestGetConnection(t *testing.T) {
	limit := 3

	connections := make(map[int]*godbcm.Connection)
	mgr := godbcm.New(limit)

	// Get all the connections in the pool
	for i := 0; i < limit; i++ {
		var err error
		connections[i], err = mgr.GetConnection()

		if err != nil {
			log.Fatalf("unexpected err: %v", err)
		}
	}

	// Attempt to get a new connection when the pool is empty
	if _, err := mgr.GetConnection(); err == nil {
		log.Fatalf("expected error but got none")
	}
}

func TestWaitForConnection(t *testing.T) {
	delay := 500 * time.Millisecond

	// Create the initial connection
	mgr := godbcm.New(1)
	connection, _ := mgr.GetConnection()

	// Setup a routine to release the initial connection
	// after the delay (delay/2). Doing this will allow the polling
	// mechanism defined after this routine to execute
	go func(mgr *godbcm.ConnectionManager, connID uuid.UUID, delay time.Duration) {
		time.Sleep(delay)
		mgr.ReleaseConnection(connID)
	}(mgr, connection.ID, delay/2)

	// Poll for a new connection, waiting for the initial connection
	// to be released
	if _, err := mgr.WaitForConnection(delay); err != nil {
		log.Fatalf("could not get a connection")
	}

	// Try to get a new connection without releasing the
	// previous connection
	if _, err := mgr.WaitForConnection(delay); err == nil {
		log.Fatalf("expected an error getting a connection")
	}
}

func TestReleaseConnection(t *testing.T) {
	mgr := godbcm.New(1)

	// Get a connection and release it immediately
	connection, _ := mgr.GetConnection()
	mgr.ReleaseConnection(connection.ID)

	// Get a new connection
	if _, err := mgr.GetConnection(); err != nil {
		log.Fatalf("could not get new connection: %v", err)
	}

	// Attempt to release non-existant connection
	if err := mgr.ReleaseConnection(uuid.New()); err == nil {
		log.Fatalf("expected error but got none")
	}
}
