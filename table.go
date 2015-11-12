package sqruct

import (
	"bytes"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"
)

// Table represents database table.
type Table struct {
	parent              *Sqruct
	omitMethod          map[string]struct{}
	CreateTableTemplate string
	DropTableTemplate   string
	SourceTemplate      string
	GoName              string
	Column              []*Column
	ColumnAfter         []string
	PrimaryKey          PrimaryKey
	ForeignKey          []ForeignKey
}

// PrimarkyKey represents primary key constraint in database table.
type PrimaryKey struct {
	Column []*Column
}

// Composite reports whether pk is composite primary key.
func (pk *PrimaryKey) Composite() bool {
	return len(pk.Column) > 1
}

// ForeignKey represents foreign key constraint in database table.
type ForeignKey struct {
	Table  *Table
	Column []ColumnPair
	Mirror bool
}

// ColumnPair represents column pair. it is used in ForeignKey.
type ColumnPair struct {
	Self  *Column
	Other *Column
}

// ColumnByName finds column by column name.
func (t *Table) ColumnByName(s string) *Column {
	if t == nil {
		return nil
	}

	s = strings.ToLower(s)
	for _, v := range t.Column {
		if v.SQLName() == s {
			return v
		}
	}
	return nil
}

// CompositePrimaryKey reports whether this table has composite primary key constraint.
func (t *Table) CompositePrimaryKey() bool {
	return t.PrimaryKey.Composite()
}

// NonPrimaryKeys returns columns which is non-primary key.
func (t *Table) NonPrimaryKeys() []*Column {
	r := []*Column{}
loop:
	for _, col := range t.Column {
		for _, v := range t.PrimaryKey.Column {
			if v == col {
				continue loop
			}
		}
		r = append(r, col)
	}
	return r
}

// AutoIncrementColumn returns column which has auto increment constraint.
func (t *Table) AutoIncrementColumn() *Column {
	for _, col := range t.Column {
		if col.AutoIncrement() {
			return col
		}
	}
	return nil
}

// PackageName returns package name in Go.
func (t *Table) PackageName() string {
	return t.parent.Config.Package
}

// SQLName returns table name in RDBMS.
func (t *Table) SQLName() string {
	return strings.ToLower(t.GoName)
}

// Mode returns current database mode.
func (t *Table) Mode() Mode {
	return t.parent.Config.Mode
}

// OmitMethod reports whether methodName is omitted in Go source.
func (t *Table) OmitMethod(methodName string) bool {
	_, ok := t.omitMethod[methodName]
	return ok
}

func render(tpl string, tpl2 string, tpl3 string, data interface{}) ([]byte, error) {
	if tpl == "" {
		if tpl2 != "" {
			tpl = tpl2
		} else {
			tpl = tpl3
		}
	}
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBufferString("")
	err = t.Execute(b, data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// CreateTableSQL returns create table SQL.
func (t *Table) CreateTableSQL() (MultiLineText, error) {
	b, err := render(t.CreateTableTemplate, t.parent.CreateTableTemplate, createTableTemplate, t)
	if err != nil {
		return "", err
	}
	return MultiLineText(b), nil
}

// MustCreateTableSQL is like CreateTableSQL but panics if an error occurred.
func (t *Table) MustCreateTableSQL() MultiLineText {
	s, err := t.CreateTableSQL()
	if err != nil {
		panic(err)
	}
	return s
}

// DropTableSQL returns drop table SQL.
func (t *Table) DropTableSQL() (MultiLineText, error) {
	b, err := render(t.DropTableTemplate, t.parent.DropTableTemplate, dropTableTemplate, t)
	if err != nil {
		return "", err
	}
	return MultiLineText(b), nil
}

// MustDropTableSQL is like DropTableSQL but panics if an error occurred.
func (t *Table) MustDropTableSQL() MultiLineText {
	s, err := t.DropTableSQL()
	if err != nil {
		panic(err)
	}
	return s
}

// Source returns source code in Go.
func (t *Table) Source() (string, error) {
	b, err := render(t.SourceTemplate, t.parent.SourceTemplate, sourceTemplate, t)
	if err != nil {
		return "", err
	}

	b, err = imports.Process("", b, nil)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
