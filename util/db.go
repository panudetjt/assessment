package util

import (
	"database/sql"
	"errors"
)

type IDB interface {
	QueryRow(query string, args ...interface{}) IRow
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type IRow interface {
	Scan(dest ...any) error
}

type Row struct {
	*sql.Row
}

type DB struct {
	*sql.DB
}

func (d *DB) QueryRow(query string, args ...interface{}) IRow {
	return &Row{d.DB.QueryRow(query, args...)}
}

type MockRow struct {
	Row          any
	LastInsertId int64
	RowsAffected int64
}

func (m *MockRow) Scan(dest ...any) error {
	for _, dp := range dest {
		d, _ := dp.(*int)
		*d = int(m.LastInsertId)
	}
	return nil
}

type MockRowError struct {
	MockRow
}

func (m *MockRowError) Scan(dest ...any) error {
	return errors.New("error")
}

type MockDB struct {
	Query        string
	LastInsertId int64
	RowsAffected int64
	Row          IRow
}

func (m *MockDB) QueryRow(query string, args ...interface{}) IRow {
	m.Query = query
	if m.Row != nil {
		return m.Row
	}
	return &MockRow{LastInsertId: m.LastInsertId, RowsAffected: m.RowsAffected}
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	m.Query = query
	return nil, nil
}
