package gen

import (
	"database/sql"
	"regexp"
	"strings"
)

// DB represents subset of database/sql.DB or database/sql.Tx.
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Mode represents Sqruct processing mode.
type Mode interface {
	// String implemenets the Stringer interface. it return such as "MySQL".
	String() string
	// IsAutoIncrement reports whether this column has auto increment constraint.
	IsAutoIncrement(col string) bool
	Quote(string) string
	Unquote(string) string
}

var (
	MySQL      Mode = mySQL{}
	PostgreSQL Mode = postgreSQL{}
	SQLite     Mode = sqlite{}
)

type mySQL struct{}

func (mySQL) String() string { return "MySQL" }

func (mySQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_INCREMENT`).MatchString(col)
}

func (mySQL) Quote(s string) string {
	return "`" + strings.Replace(s, "`", "``", -1) + "`"
}

func (mySQL) Unquote(s string) string {
	if len(s) < 2 || s[0] != '`' || s[len(s)-1] != '`' {
		return s
	}
	return strings.Replace(s[1:len(s)-1], "``", "`", -1)
}

type postgreSQL struct{}

func (postgreSQL) String() string { return "PostgreSQL" }

func (postgreSQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)(?:BIG)?SERIAL|\S+\s+DEFAULT\s+nextval`).MatchString(col)
}

func (postgreSQL) Quote(s string) string {
	return `"` + strings.Replace(s, `"`, `""`, -1) + `"`
}

func (postgreSQL) Unquote(s string) string {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return s
	}
	return strings.Replace(s[1:len(s)-1], `""`, `"`, -1)
}

type sqlite struct{}

func (sqlite) String() string { return "SQLite" }

func (sqlite) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_?INCREMENT|INTEGER\s+PRIMARY\s+KEY`).MatchString(col)
}

func (sqlite) Quote(s string) string {
	return `"` + strings.Replace(s, `"`, `""`, -1) + `"`
}

func (sqlite) Unquote(s string) string {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return s
	}
	return strings.Replace(s[1:len(s)-1], `""`, `"`, -1)
}
