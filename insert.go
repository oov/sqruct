package sqruct

import "database/sql"

// Ext represents github.com/jmoiron/sqlx.Ext.
type Ext interface {
	BindNamed(string, interface{}) (string, []interface{}, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func insertExec(e Ext, query string, table interface{}, needLastInsertID bool) (int64, error) {
	q, args, err := e.BindNamed(query, table)
	if err != nil {
		return 0, err
	}
	r, err := e.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	if needLastInsertID {
		return r.LastInsertId()
	}
	return 0, nil
}

func insertGet(e Ext, query string, table interface{}, needLastInsertID bool) (int64, error) {
	q, args, err := e.BindNamed(query, table)
	if err != nil {
		return 0, err
	}
	if needLastInsertID {
		r, err := e.Query(q, args...)
		if err != nil {
			return 0, err
		}
		defer r.Close()
		if !r.Next() {
			if err = r.Err(); err != nil {
				return 0, err
			}
			return 0, sql.ErrNoRows
		}
		var i int64
		err = r.Scan(&i)
		return i, err
	}
	_, err = e.Exec(q, args...)
	return 0, err
}

func insert(e Ext, table DBTable, useAutoIncrement bool, defaultVal string, useReturning bool) (int64, error) {
	idx := -1
	if useAutoIncrement {
		idx = table.AutoIncrementColumnIndex()
	}
	q := buildInsert(table.TableName(), table.Columns(), idx, defaultVal, useReturning)
	if useReturning {
		return insertGet(e, q, table, useAutoIncrement)
	}
	return insertExec(e, q, table, useAutoIncrement)
}
