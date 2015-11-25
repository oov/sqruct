package sqruct

import (
	"regexp"
	"strings"
)

// Mode represents Sqruct processing mode.
type Mode interface {
	// String implemenets the Stringer interface. it return such as "MySQL".
	String() string
	// DefaultValueKeyword returns default value keyword that is used in insert statements.
	DefaultValueKeyword() string
	// IsAutoIncrement reports whether this column has auto increment constraint.
	IsAutoIncrement(col string) bool
	// Insert executes insert statement on db.
	Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error)
	// Placeholder creates placeholder generator that is used in SQL statements.
	Placeholder() Placeholder
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

func (mySQL) DefaultValueKeyword() string { return "DEFAULT" }

func (mySQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_INCREMENT`).MatchString(col)
}

func (m mySQL) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return genericInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.Placeholder())
}

func (mySQL) Placeholder() Placeholder {
	return &genericPlaceholder{}
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

func (postgreSQL) DefaultValueKeyword() string { return "DEFAULT" }

func (postgreSQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)(?:BIG)?SERIAL|\S+\s+DEFAULT\s+nextval`).MatchString(col)
}

func (m postgreSQL) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return postgresInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.Placeholder())
}

func (postgreSQL) Placeholder() Placeholder {
	return &postgresPlaceholder{}
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

func (sqlite) DefaultValueKeyword() string { return "NULL" }

func (sqlite) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_?INCREMENT|INTEGER\s+PRIMARY\s+KEY`).MatchString(col)
}

func (m sqlite) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return genericInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.Placeholder())
}

func (sqlite) Placeholder() Placeholder {
	return &genericPlaceholder{}
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
