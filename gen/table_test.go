package gen

import (
	"testing"
)

func TestMultiLineTextAddPrefix(t *testing.T) {
	sqls := []struct {
		Prefix string
		Before string
		After  string
	}{
		{
			Prefix: `-- `,
			Before: "SOME QUERY;\nSOME QUERY2;",
			After:  "-- SOME QUERY;\n-- SOME QUERY2;",
		},
		{
			Prefix: `-- `,
			Before: "SOME QUERY;\r\nSOME QUERY2;\r\n",
			After:  "-- SOME QUERY;\r\n-- SOME QUERY2;\r\n-- ",
		},
	}
	for i, v := range sqls {
		s := MultiLineText(v.Before).AddPrefix(v.Prefix)
		if string(s) != v.After {
			t.Errorf("sqls[%d]: want %q got %q", i, v.After, s)
		}
	}
}

func TestMultiLineTextAddPostfix(t *testing.T) {
	sqls := []struct {
		Postfix string
		Before  string
		After   string
	}{
		{
			Postfix: `--`,
			Before:  "SOME QUERY;\nSOME QUERY2;",
			After:   "SOME QUERY;--\nSOME QUERY2;--",
		},
		{
			Postfix: `--`,
			Before:  "SOME QUERY;\r\nSOME QUERY2;\r\n",
			After:   "SOME QUERY;--\r\nSOME QUERY2;--\r\n--",
		},
	}
	for i, v := range sqls {
		s := MultiLineText(v.Before).AddPostfix(v.Postfix)
		if string(s) != v.After {
			t.Errorf("sqls[%d]: want %q got %q", i, v.After, s)
		}
	}
}
