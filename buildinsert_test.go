package sqruct

import (
	"testing"
)

func TestBuildInsert(t *testing.T) {
	datas := []struct {
		Output       string
		Table        string
		Columns      []string
		AutoIncrCol  int
		DefValue     string
		UseReturning bool
	}{
		// MySQL like
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (:c1, :c2, :c3)",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  -1,
			DefValue:     "DEFAULT",
			UseReturning: false,
		},
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (:c1, :c2, DEFAULT)",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  2,
			DefValue:     "DEFAULT",
			UseReturning: false,
		},
		// PostgreSQL like
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (:c1, :c2, :c3)",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  -1,
			UseReturning: true,
		},
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (DEFAULT, :c2, :c3) RETURNING c1",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  0,
			DefValue:     "DEFAULT",
			UseReturning: true,
		},
		// SQLite like
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (:c1, :c2, :c3)",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  -1,
			DefValue:     "NULL",
			UseReturning: false,
		},
		{
			Output:       "INSERT INTO h (c1, c2, c3) VALUES (NULL, :c2, :c3)",
			Table:        "h",
			Columns:      []string{"c1", "c2", "c3"},
			AutoIncrCol:  0,
			DefValue:     "NULL",
			UseReturning: false,
		},
	}
	for i, v := range datas {
		q := buildInsert(v.Table, v.Columns, v.AutoIncrCol, v.DefValue, v.UseReturning)
		if q != v.Output {
			t.Errorf("buildInsertQueryData[%d] want %q got %q", i, v.Output, q)
		}
	}
}

func BenchmarkBuildInsert(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buildInsert("hello", []string{"column1", "column2", "column3"}, 2, "DEFAULT", true)
	}
}
