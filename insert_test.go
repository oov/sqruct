package sqruct

import (
	"database/sql"
	"testing"
)

func TestBuildInsert(t *testing.T) {
	datas := []struct {
		Output      string
		Cap         int
		Table       string
		Columns     []string
		AutoIncrCol int
		Mode        Mode
	}{
		// MySQL like
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES (?, ?, ?)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES (?, ?, ?) RETURNING "),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: -1,
			Mode:        MySQL,
		},
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES (?, ?, DEFAULT)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES (?, ?, DEFAULT) RETURNING c3"),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: 2,
			Mode:        MySQL,
		},
		// PostgreSQL like
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES ($1, $2, $3)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES ($1, $2, $3) RETURNING "),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: -1,
			Mode:        PostgreSQL,
		},
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES (DEFAULT, $1, $2)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES (DEFAULT, $1, $2) RETURNING c1"),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: 0,
			Mode:        PostgreSQL,
		},
		// SQLite like
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES (?, ?, ?)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES (?, ?, ?) RETURNING "),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: -1,
			Mode:        SQLite,
		},
		{
			Output:      "INSERT INTO h (c1, c2, c3) VALUES (NULL, ?, ?)",
			Cap:         len("INSERT INTO h (c1, c2, c3) VALUES (NULL, ?, ?) RETURNING c1"),
			Table:       "h",
			Columns:     []string{"c1", "c2", "c3"},
			AutoIncrCol: 0,
			Mode:        SQLite,
		},
	}
	for i, v := range datas {
		q := buildInsert(v.Table, v.Columns, v.AutoIncrCol, v.Mode.DefaultValueKeyword(), v.Mode.Placeholder())
		if string(q) != v.Output {
			t.Errorf("buildInsertQueryData[%d] want %q got %q", i, v.Output, string(q))
		}
		if cap(q) != v.Cap {
			t.Errorf("buildInsertQueryData[%d] cap want %d got %d", i, v.Cap, cap(q))
		}
	}
}

func BenchmarkBuildInsert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buildInsert("hello", []string{"column1", "column2", "column3"}, 2, "DEFAULT", &genericPlaceholder{})
	}
}

func TestAutoIncrementOnInsert(t *testing.T) {
	type Data struct {
		Values       []interface{}
		LastInsertID int64
		Record       []int64
	}
	testSet := []struct {
		AutoIncrementColumnIndex int
		Create                   map[string]string
		Drop                     map[string]string
		Datas                    []Data
	}{
		{
			AutoIncrementColumnIndex: 0,
			Create: map[string]string{
				"MySQL":      `CREATE TABLE tbl(a INTEGER AUTO_INCREMENT PRIMARY KEY, b INT, c INT)`,
				"PostgreSQL": `CREATE TABLE tbl(a SERIAL PRIMARY KEY, b INT, c INT)`,
				"SQLite":     `CREATE TABLE tbl(a INTEGER PRIMARY KEY, b INT, c INT)`,
			},
			Drop: map[string]string{
				"MySQL":      `DROP TABLE tbl`,
				"PostgreSQL": `DROP TABLE tbl`,
				"SQLite":     `DROP TABLE tbl`,
			},
			Datas: []Data{
				{
					Values:       []interface{}{0, 0, 10},
					LastInsertID: 1,
					Record:       []int64{1, 0, 10},
				},
				{
					Values:       []interface{}{0, 100, 200},
					LastInsertID: 2,
					Record:       []int64{2, 100, 200},
				},
				{
					Values:       []interface{}{10, 200, 300},
					LastInsertID: 0,
					Record:       []int64{10, 200, 300},
				},
			},
		},
		{
			AutoIncrementColumnIndex: 1,
			Create: map[string]string{
				"MySQL":      `CREATE TABLE tbl(a INT, b INTEGER AUTO_INCREMENT PRIMARY KEY, c INT)`,
				"PostgreSQL": `CREATE TABLE tbl(a INT, b SERIAL PRIMARY KEY, c INT)`,
				"SQLite":     `CREATE TABLE tbl(a INT, b INTEGER PRIMARY KEY, c INT)`,
			},
			Drop: map[string]string{
				"MySQL":      `DROP TABLE tbl`,
				"PostgreSQL": `DROP TABLE tbl`,
				"SQLite":     `DROP TABLE tbl`,
			},
			Datas: []Data{
				{
					Values:       []interface{}{0, 0, 10},
					LastInsertID: 1,
					Record:       []int64{0, 1, 10},
				},
				{
					Values:       []interface{}{100, 0, 200},
					LastInsertID: 2,
					Record:       []int64{100, 2, 200},
				},
				{
					Values:       []interface{}{10, 200, 300},
					LastInsertID: 0,
					Record:       []int64{10, 200, 300},
				},
			},
		},
		{
			AutoIncrementColumnIndex: 2,
			Create: map[string]string{
				"MySQL":      `CREATE TABLE tbl(a INT, b INT, c INTEGER AUTO_INCREMENT PRIMARY KEY)`,
				"PostgreSQL": `CREATE TABLE tbl(a INT, b INT, c SERIAL PRIMARY KEY)`,
				"SQLite":     `CREATE TABLE tbl(a INT, b INT, c INTEGER PRIMARY KEY)`,
			},
			Drop: map[string]string{
				"MySQL":      `DROP TABLE tbl`,
				"PostgreSQL": `DROP TABLE tbl`,
				"SQLite":     `DROP TABLE tbl`,
			},
			Datas: []Data{
				{
					Values:       []interface{}{0, 10, 0},
					LastInsertID: 1,
					Record:       []int64{0, 10, 1},
				},
				{
					Values:       []interface{}{100, 200, 0},
					LastInsertID: 2,
					Record:       []int64{100, 200, 2},
				},
				{
					Values:       []interface{}{10, 200, 300},
					LastInsertID: 0,
					Record:       []int64{10, 200, 300},
				},
			},
		},
	}
	columns := []string{"a", "b", "c"}
	for _, tester := range []Tester{mySQLTest, postgreSQLTest, sqliteTest} {
		err := tester(func(db *sql.DB, mode Mode) {
			for testSetIdx, test := range testSet {
				func() {
					if _, err := db.Exec(test.Create[mode.String()]); err != nil {
						t.Fatal(err)
					}
					defer func() {
						if _, err := db.Exec(test.Drop[mode.String()]); err != nil {
							t.Fatal(err)
						}
					}()
					for dataIdx, v := range test.Datas {
						r, err := mode.Insert(db, "tbl", columns, v.Values, test.AutoIncrementColumnIndex)
						if err != nil {
							t.Fatalf("testSet[%d] %s.Insert failed: %v", testSetIdx, mode, err)
						}
						if r != v.LastInsertID {
							t.Errorf("testSet[%d] %s datas[%d] LastInsertID want %d got %d", testSetIdx, mode, dataIdx, v.LastInsertID, r)
						}
						rr, rrp := make([]int64, len(v.Record)), make([]interface{}, len(v.Record))
						for i := range rr {
							rrp[i] = &rr[i]
						}
						err = db.QueryRow(
							mode.Placeholder().Rebind(`SELECT * FROM tbl WHERE `+columns[test.AutoIncrementColumnIndex]+` = ?`),
							v.Record[test.AutoIncrementColumnIndex],
						).Scan(rrp...)
						if err != nil {
							t.Fatalf("testSet[%d] %s SELECT failed: %v", testSetIdx, mode, err)
						}
						for i := range rr {
							if rr[i] != v.Record[i] {
								t.Errorf("testSet[%d] %s datas[%d] Record[%d] want %d got %d", testSetIdx, mode, dataIdx, i, v.Record[i], rr[i])
							}
						}
					}
				}()
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
