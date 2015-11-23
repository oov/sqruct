package sqruct

import (
	"database/sql"
	"testing"
)

func TestAutoIncrementOnInsert(t *testing.T) {
	tests := []struct {
		Create string
		Drop   string
		Tester Tester
	}{
		{
			Create: `CREATE TABLE tbl(a INTEGER AUTO_INCREMENT PRIMARY KEY, b INT)`,
			Drop:   `DROP TABLE tbl`,
			Tester: mySQLTest,
		},
		{
			Create: `CREATE TABLE tbl(a SERIAL PRIMARY KEY, b INT)`,
			Drop:   `DROP TABLE tbl`,
			Tester: postgreSQLTest,
		},
		{
			Create: `CREATE TABLE tbl(a INTEGER AUTO_INCREMENT PRIMARY KEY, b INT)`,
			Drop:   `DROP TABLE tbl`,
			Tester: sqliteTest,
		},
	}

	for _, v := range tests {
		err := v.Tester(func(db *sql.DB, mode Mode) {
			if _, err := db.Exec(v.Create); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if _, err := db.Exec(v.Drop); err != nil {
					t.Fatal(err)
				}
			}()
			r, err := mode.Insert(db, "tbl", []string{"a", "b"}, []interface{}{0, 0}, 0)
			if err != nil {
				t.Error(err)
			}
			if r != 1 {
				t.Errorf("want 1 got %d", r)
			}
			r, err = mode.Insert(db, "tbl", []string{"a", "b"}, []interface{}{0, 0}, 0)
			if err != nil {
				t.Error(err)
			}
			if r != 2 {
				t.Errorf("want 2 got %d", r)
			}
			r, err = mode.Insert(db, "tbl", []string{"a", "b"}, []interface{}{10, 0}, 0)
			if err != nil {
				t.Error(err)
			}
			if r != 0 {
				t.Errorf("want 0 got %d", r)
			}
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
