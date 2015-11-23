package sqruct

import (
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
		q := buildInsert(v.Table, v.Columns, v.AutoIncrCol, v.Mode.DefaultValueKeyword(), v.Mode.PlaceholderGenerator())
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
		buildInsert("hello", []string{"column1", "column2", "column3"}, 2, "DEFAULT", genericPlaceholderGenerator{})
	}
}
