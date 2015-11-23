package sqruct

import (
	"regexp"
	"strings"
)

// Mode represents Sqruct processing mode.
type Mode interface {
	String() string
	DefaultValueKeyword() string
	IsAutoIncrement(col string) bool
	Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error)
	PlaceholderGenerator() PlaceholderGenerator
	Rebind(string) string
}

var (
	MySQL      Mode = mySQL{}
	PostgreSQL Mode = postgreSQL{}
	SQLite     Mode = sqlite{}
)

type mySQL struct{}

// String implemenets the Stringer interface.
func (mySQL) String() string { return "MySQL" }

func (mySQL) DefaultValueKeyword() string { return "DEFAULT" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (mySQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_INCREMENT`).MatchString(col)
}

// Insert executes insert statement on db.
func (m mySQL) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return genericInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.PlaceholderGenerator())
}

func (mySQL) PlaceholderGenerator() PlaceholderGenerator {
	return genericPlaceholderGenerator{}
}

func (mySQL) Rebind(s string) string { return s }

type postgreSQL struct{}

// String implemenets the Stringer interface.
func (postgreSQL) String() string { return "PostgreSQL" }

func (postgreSQL) DefaultValueKeyword() string { return "DEFAULT" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (postgreSQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)(?:BIG)?SERIAL|\S+\s+DEFAULT\s+nextval`).MatchString(col)
}

// Insert executes insert statement on db.
func (m postgreSQL) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return postgresInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.PlaceholderGenerator())
}

func (postgreSQL) PlaceholderGenerator() PlaceholderGenerator {
	var r postgresPlaceholderGenerator
	return &r
}

func (m postgreSQL) Rebind(s string) string {
	g := m.PlaceholderGenerator()
	r := make([]byte, 0, len(s)+8)
	var p int
	for {
		p = strings.IndexByte(s, '?')
		if p == -1 {
			r = append(r, s...)
			break
		}
		r = append(r, s[:p]...)
		r = append(r, g.Placeholder()...)
		s = s[p+1:]
	}
	return string(r)
}

type sqlite struct{}

// String implemenets the Stringer interface.
func (sqlite) String() string { return "SQLite" }

func (sqlite) DefaultValueKeyword() string { return "NULL" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (sqlite) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_?INCREMENT|INTEGER\s+PRIMARY\s+KEY`).MatchString(col)
}

// Insert executes insert statement on db.
func (m sqlite) Insert(db DB, table string, columns []string, values []interface{}, autoIncrColumn int) (int64, error) {
	return genericInsert(db, table, columns, values, autoIncrColumn, m.DefaultValueKeyword(), m.PlaceholderGenerator())
}

func (sqlite) PlaceholderGenerator() PlaceholderGenerator {
	return genericPlaceholderGenerator{}
}

func (sqlite) Rebind(s string) string { return s }
