package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MySQLPersister struct {
	Conn          *sqlx.DB
	ch            CryptoHasher
	userTxMethods map[string]func(User, ...interface{}) (*User, error)
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

	mp := &MySQLPersister{
		Conn: conn,
		ch:   ch,
	}
	// define list of possible ExecuteTx operations
	mp.userTxMethods = map[string]func(User, ...interface{}) (*User, error){
		"INSERT": mp.InsertUser,
		"UPDATE": mp.UpdateUser,
	}
	return mp, nil
}





