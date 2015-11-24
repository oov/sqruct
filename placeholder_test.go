package sqruct

import "testing"

func TestPlaceholderGeneratorRebind(t *testing.T) {
	tests := []struct {
		q string
		o string
		m Mode
	}{
		{
			q: "SELECT * FROM table WHERE (hello = ?)AND(world = ?)",
			o: "SELECT * FROM table WHERE (hello = ?)AND(world = ?)",
			m: MySQL,
		},
		{
			q: "SELECT * FROM table WHERE (hello = ?)AND(world = ?)",
			o: "SELECT * FROM table WHERE (hello = $1)AND(world = $2)",
			m: PostgreSQL,
		},
		{
			q: "SELECT * FROM table WHERE (hello = ?)AND(world = ?)",
			o: "SELECT * FROM table WHERE (hello = ?)AND(world = ?)",
			m: SQLite,
		},
		{
			q: "INSERT INTO table (a, b, c, d, e, f, g, h, i, j) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			o: "INSERT INTO table (a, b, c, d, e, f, g, h, i, j) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			m: PostgreSQL,
		},
	}
	for i, v := range tests {
		if r := v.m.Placeholder().Rebind(v.q); r != v.o {
			t.Errorf("tests[%d] %s want %q got %q", i, v.m, v.o, r)
		}
	}
}
