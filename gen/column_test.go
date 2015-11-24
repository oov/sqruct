package gen

import (
	"testing"

	"github.com/oov/sqruct"
)

func TestColumnNotNull(t *testing.T) {
	tests := []struct {
		c string
		o bool
	}{
		{c: "INTEGER PRIMARY KEY", o: true},
		{c: "INTEGER UNIQUE NOT NULL", o: true},
		{c: "INTEGER UNIQUE", o: false},
		{c: "INTEGER", o: false},
	}
	for i, v := range tests {
		if r := (&Column{SQLColumn: v.c}).NotNull(); r != v.o {
			t.Errorf("tests[%d] %q want %t got %t", i, v.c, v.o, r)
		}
	}
}

func TestColumnUnique(t *testing.T) {
	tests := []struct {
		c string
		o bool
	}{
		{c: "INTEGER UNIQUE", o: true},
		{c: "INTEGER PRIMARY KEY", o: true},
		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: true},
		{c: "INTEGER", o: false},
	}
	for i, v := range tests {
		if r := (&Column{SQLColumn: v.c}).Unique(); r != v.o {
			t.Errorf("tests[%d] %q want %t got %t", i, v.c, v.o, r)
		}
	}
}

func TestColumnAutoIncrement(t *testing.T) {
	tests := []struct {
		c string
		o bool
		m sqruct.Mode
	}{
		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: true, m: sqruct.MySQL},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: false, m: sqruct.MySQL},
		{c: "INTEGER PRIMARY KEY", o: false, m: sqruct.MySQL},
		{c: "SERIAL", o: false, m: sqruct.MySQL},
		{c: "BIGSERIAL", o: false, m: sqruct.MySQL},

		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: false, m: sqruct.PostgreSQL},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: false, m: sqruct.PostgreSQL},
		{c: "INTEGER PRIMARY KEY", o: false, m: sqruct.PostgreSQL},
		{c: "SERIAL", o: true, m: sqruct.PostgreSQL},
		{c: "BIGSERIAL", o: true, m: sqruct.PostgreSQL},

		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: true, m: sqruct.SQLite},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: true, m: sqruct.SQLite},
		{c: "INTEGER PRIMARY KEY", o: true, m: sqruct.SQLite},
		{c: "INT PRIMARY KEY", o: false, m: sqruct.SQLite},
		{c: "SERIAL", o: false, m: sqruct.SQLite},
		{c: "BIGSERIAL", o: false, m: sqruct.SQLite},
	}
	for i, v := range tests {
		if r := (&Column{parent: &Table{parent: &Sqruct{Config: Config{Mode: v.m}}}, SQLColumn: v.c}).AutoIncrement(); r != v.o {
			t.Errorf("tests[%d] %q want %t got %t", i, v.c, v.o, r)
		}
	}
}

func TestColumnDefault(t *testing.T) {
	tests := []struct {
		c string
		o bool
		m sqruct.Mode
	}{
		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: true, m: sqruct.MySQL},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: false, m: sqruct.MySQL},
		{c: "INTEGER PRIMARY KEY", o: false, m: sqruct.MySQL},
		{c: "SERIAL", o: false, m: sqruct.MySQL},
		{c: "BIGSERIAL", o: false, m: sqruct.MySQL},
		{c: "INTEGER NOT NULL", o: false, m: sqruct.MySQL},
		{c: "INTEGER NOT NULL DEFAULT 0", o: true, m: sqruct.MySQL},

		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: false, m: sqruct.PostgreSQL},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: false, m: sqruct.PostgreSQL},
		{c: "INTEGER PRIMARY KEY", o: false, m: sqruct.PostgreSQL},
		{c: "SERIAL", o: true, m: sqruct.PostgreSQL},
		{c: "BIGSERIAL", o: true, m: sqruct.PostgreSQL},
		{c: "INTEGER NOT NULL", o: false, m: sqruct.PostgreSQL},
		{c: "INTEGER NOT NULL DEFAULT 0", o: true, m: sqruct.PostgreSQL},

		{c: "INTEGER PRIMARY KEY AUTO_INCREMENT", o: true, m: sqruct.SQLite},
		{c: "INTEGER PRIMARY KEY AUTOINCREMENT", o: true, m: sqruct.SQLite},
		{c: "INTEGER PRIMARY KEY", o: true, m: sqruct.SQLite},
		{c: "INT PRIMARY KEY", o: false, m: sqruct.SQLite},
		{c: "SERIAL", o: false, m: sqruct.SQLite},
		{c: "BIGSERIAL", o: false, m: sqruct.SQLite},
		{c: "INTEGER NOT NULL", o: false, m: sqruct.SQLite},
		{c: "INTEGER NOT NULL DEFAULT 0", o: true, m: sqruct.SQLite},
	}
	for i, v := range tests {
		if r := (&Column{parent: &Table{parent: &Sqruct{Config: Config{Mode: v.m}}}, SQLColumn: v.c}).Default(); r != v.o {
			t.Errorf("tests[%d] %q want %t got %t", i, v.c, v.o, r)
		}
	}
}
