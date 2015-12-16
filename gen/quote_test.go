package gen

import "testing"

func TestQuote(t *testing.T) {
	datas := []struct {
		in  string
		out string
		m   Mode
	}{
		{
			in:  "hello",
			out: "`hello`",
			m:   MySQL,
		},
		{
			in:  `heello`,
			out: `"heello"`,
			m:   PostgreSQL,
		},
		{
			in:  `heello`,
			out: `"heello"`,
			m:   SQLite,
		},
	}
	for _, v := range datas {
		if r := v.m.Quote(v.in); r != v.out {
			t.Errorf("%s.Quote want %q got %q", v.m, v.out, r)
		}
		if r := v.m.Unquote(v.out); r != v.in {
			t.Errorf("%s.Unquote want %q got %q", v.m, v.in, r)
		}
	}
}
