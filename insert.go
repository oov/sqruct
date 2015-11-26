package sqruct

import "database/sql"

// DB represents subset of database/sql.DB or database/sql.Tx.
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// buildInsert builds SQL such as "INSERT INTO `table` (`column1`, `column2`) VALUES (?, ?)".
func buildInsert(table string, columns []string, autoIncrCol int, m Mode) []byte {
	defValue := m.DefaultValueKeyword()
	ph := m.Placeholder()
	q := make([]byte, 0, 128)
	q = append(q, "INSERT INTO "...)
	q = append(q, m.Quote(table)...)
	q = append(q, " ("...)
	for i, c := range columns {
		if i != 0 {
			q = append(q, ", "...)
		}
		q = append(q, m.Quote(c)...)
	}
	q = append(q, ") VALUES ("...)
	for i := range columns {
		if i != 0 {
			q = append(q, ", "...)
		}
		if i == autoIncrCol {
			q = append(q, defValue...)
		} else {
			q = append(q, ph.Next()...)
		}
	}
	q = append(q, ')')
	return q
}

func dropColumn(values []interface{}, index int) []interface{} {
	if index == 0 {
		return values[1:]
	}
	if index == len(values)-1 {
		return values[:len(values)-1]
	}
	v := make([]interface{}, len(values)-1)
	copy(v, values[:index])
	copy(v[index:], values[index+1:])
	return v
}

func genericInsert(db DB, table string, columns []string, values []interface{},
	autoIncrColumn int, m Mode) (int64, error) {
	if IsZero(values[autoIncrColumn]) {
		// Drop values[autoIncrColumn] becuase used DEFAULT in this case.
		values = dropColumn(values, autoIncrColumn)
	} else {
		autoIncrColumn = -1
	}
	qb := buildInsert(table, columns, autoIncrColumn, m)
	r, err := db.Exec(string(qb), values...)
	if err != nil {
		return 0, err
	}
	if autoIncrColumn == -1 {
		return 0, nil
	}
	return r.LastInsertId()
}

func postgresInsert(db DB, table string, columns []string, values []interface{},
	autoIncrColumn int, m Mode) (int64, error) {
	if IsZero(values[autoIncrColumn]) {
		// Drop values[autoIncrColumn] because used DEFAULT in this case.
		values = dropColumn(values, autoIncrColumn)
	} else {
		autoIncrColumn = -1
	}
	qb := buildInsert(table, columns, autoIncrColumn, m)
	if autoIncrColumn == -1 {
		_, err := db.Exec(string(qb), values...)
		return 0, err
	}

	qb = append(qb, " RETURNING "...)
	qb = append(qb, columns[autoIncrColumn]...)
	var i int64
	err := db.QueryRow(string(qb), values...).Scan(&i)
	return i, err
}
