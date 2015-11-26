package gen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/oov/sqruct"

	"golang.org/x/tools/imports"
)

type Name struct {
	Go  string
	SQL string
	m   sqruct.Mode
}

func (n *Name) SQLQuoted() string {
	return n.m.Quote(n.SQL)
}

func (n *Name) SQLForGo() string {
	s := fmt.Sprintf("%q", n.SQLQuoted())
	return s[1 : len(s)-1]
}

// Table represents database table.
type Table struct {
	parent              *Sqruct
	omitMethod          map[string]struct{}
	CreateTableTemplate string
	DropTableTemplate   string
	SourceTemplate      string
	Name                Name
	Column              []*Column
	ColumnAfter         []string
	PrimaryKey          PrimaryKey
	ForeignKey          []ForeignKey
	ManyToMany          []*ManyToMany
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

func (fk *ForeignKey) Match(cols []*Column) bool {
loop:
	for _, c := range cols {
		for _, fc := range fk.Column {
			if fc.Self.Name.SQL == c.Name.SQL {
				continue loop
			}
		}
		return false // not found in fk.Column
	}
	return true
}

type ManyToMany struct {
	s        string
	RelTable *Table
	MyFK     *ForeignKey
	OtherFK  *ForeignKey
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

	for _, v := range t.Column {
		if v.Name.SQL == s {
			return v
		}
	}
	return nil
}

func (t *Table) ForeignKeyByColumns(cols []*Column) *ForeignKey {
	if t == nil {
		return nil
	}

	for i, fk := range t.ForeignKey {
		if fk.Match(cols) {
			return &t.ForeignKey[i]
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

// Mode returns current database mode.
func (t *Table) Mode() sqruct.Mode {
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
