package repo

import "database/sql"

type TxManager interface {
	BeginTx() (*sql.Tx, error)
	GetDB() *sql.DB
}

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}

func (r *txManager) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *txManager) GetDB() *sql.DB {
	return r.db
}
