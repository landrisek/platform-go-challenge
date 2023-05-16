package models

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MySQLPersister struct {
	Conn          *sqlx.DB
	ch            CryptoHasher
}

// Crypto defines abstraction for
// persister related cryptography helper
type Crypto interface {
	EncryptBytes([]byte) ([]byte, error)
	DecryptBytes([]byte) ([]byte, error)
	EncryptString(string) (string, error)
	DecryptString(string) (string, error)
}

// Hasher defines abstraction for
// persister related hashing helper
type Hasher interface {
	HashBytes([]byte) ([]byte, error)
	HashString(string) (string, error)
}

// CryptoHasher defines abstraction that
// is capable of carring persister Crypto and Hasher
type CryptoHasher struct {
	Crypto
	Hasher
}

// NewMySQLPersister will create a new mysql persistor
func NewMySQLPersister(conn *sqlx.DB, ch CryptoHasher) (*MySQLPersister, error) {
	// check that provided connection has appropriate driver
	driver := strings.ToLower(conn.DriverName())
	if !strings.Contains(driver, "mysql") && !strings.Contains(driver, "mariadb") {
		return nil, fmt.Errorf("unsupported DB driver %q", driver)
	}

	return &MySQLPersister{
		Conn: conn,
		ch:   ch,
	}, nil
}