package sqruct

import "testing"

func TestColumns(t *testing.T) {
	tests := []struct {
		o       string
		prefix  string
		columns []string
	}{
		{
			o:       "id, name, age",
			prefix:  "",
			columns: []string{"id", "name", "age"},
		},
		{
			o:       "table.id, table.name, table.age",
			prefix:  "table",
			columns: []string{"id", "name", "age"},
		},
	}
	for i, v := range tests {
		if r := Columns(v.prefix, v.columns); r != v.o {
			t.Errorf("tests[%d] want %q got %q", i, v.o, r)
		}
	}
}
