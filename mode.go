package sqruct

import "regexp"

type DBTable interface {
	TableName() string
	Columns() []string
	AutoIncrementColumnIndex() int
}

// Mode represents Sqruct processing mode.
type Mode interface {
	String() string
	IsAutoIncrement(col string) bool
	Insert(e Ext, table DBTable, useAutoIncrement bool) (int64, error)
}

type MySQL struct{}

// String implemenets the Stringer interface.
func (MySQL) String() string { return "MySQL" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (MySQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_INCREMENT`).MatchString(col)
}

// Insert executes insert statement on e.
func (MySQL) Insert(e Ext, table DBTable, useAutoIncrement bool) (int64, error) {
	return insert(e, table, useAutoIncrement, "DEFAULT", false)
}

type PostgreSQL struct{}

// String implemenets the Stringer interface.
func (PostgreSQL) String() string { return "PostgreSQL" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (PostgreSQL) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)(?:BIG)?SERIAL|\S+\s+DEFAULT\s+nextval`).MatchString(col)
}

// Insert executes insert statement on e.
func (PostgreSQL) Insert(e Ext, table DBTable, useAutoIncrement bool) (int64, error) {
	return insert(e, table, useAutoIncrement, "DEFAULT", true)
}

type SQLite struct{}

// String implemenets the Stringer interface.
func (SQLite) String() string { return "SQLite" }

// IsAutoIncrement reports whether this column has auto increment constraint.
func (SQLite) IsAutoIncrement(col string) bool {
	return regexp.MustCompile(`(?i)\S+\s+AUTO_?INCREMENT|INTEGER\s+PRIMARY\s+KEY`).MatchString(col)
}

// Insert executes insert statement on e.
func (SQLite) Insert(e Ext, table DBTable, useAutoIncrement bool) (int64, error) {
	return insert(e, table, useAutoIncrement, "NULL", false)
}
