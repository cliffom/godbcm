package godbcm_test

import (
	"log"
	"testing"

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

	for i := 0; i < limit; i++ {
		var err error
		connections[i], err = mgr.GetConnection()

		if err != nil {
			log.Fatalf("unexpected err: %v", err)
		}
	}

	if _, err := mgr.GetConnection(); err == nil {
		log.Fatalf("expected error but got none")
	}
}

func TestReleaseConnection(t *testing.T) {
	mgr := godbcm.New(1)

	connection, _ := mgr.GetConnection()
	mgr.ReleaseConnection(connection.ID)

	if _, err := mgr.GetConnection(); err != nil {
		log.Fatalf("could not get new connection: %v", err)
	}
}

func TestReleaseInvalidConnection(t *testing.T) {
	mgr := godbcm.New(1)
	if err := mgr.ReleaseConnection(uuid.New()); err == nil {
		log.Fatalf("expected error but got none")
	}
}
