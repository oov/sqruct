package sqruct

import (
	"testing"
)

func TestBuildInsertQuery(t *testing.T) {
	datas := []struct {
		Output  string
		Table   string
		Columns []string
		Use     []bool
	}{
		{
			Output:  "INSERT INTO hello (column1, column2, column3) VALUES (:column1, :column2, NULL)",
			Table:   "hello",
			Columns: []string{"column1", "column2", "column3"},
			Use:     []bool{true, true, false},
		},
		{
			Output:  "INSERT INTO h (c1, c2, c3) VALUES (:c1, :c2, NULL)",
			Table:   "h",
			Columns: []string{"c1", "c2", "c3"},
			Use:     []bool{true, true, false},
		},
	}
	for i, v := range datas {
		q := BuildInsertQuery(v.Table, v.Columns, v.Use)
		if q != v.Output {
			t.Errorf("buildInsertQueryData[%d] want %q got %q", i, v.Output, q)
		}
	}
}

func BenchmarkBuildInsertQuery(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BuildInsertQuery("hello", []string{"column1", "column2", "column3"}, []bool{true, true, false})
	}
}