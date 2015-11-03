// Package sqruct is database table to struct mapping tool.
package sqruct

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents Sqruct configuration.
type Config struct {
	Mode    string // Processing mode, such as "sqlite".
	Package string // Package name in source code in Go.
	Tag     string // Tag name in struct definition in Go.
	Dir     string // Source code output directory.
}

// Sqruct...
type Sqruct struct {
	CreateTableTemplate string
	DropTableTemplate   string
	SourceTemplate      string
	Table               []*Table
	Config              Config
}

// Export writes sources to files.
func (sq *Sqruct) Export(baseDir string) error {
	d := path.Join(baseDir, sq.Config.Dir)
	for _, t := range sq.Table {
		src, err := t.Source()
		if err != nil {
			return err
		}
		f, err := os.Create(path.Join(d, "zz"+strings.ToLower(t.GoName)+".go"))
		if err != nil {
			return err
		}
		if _, err = f.WriteString(src); err != nil {
			f.Close()
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (sq *Sqruct) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var err error
	var ms yaml.MapSlice
	if err = unmarshal(&ms); err != nil {
		return err
	}

	// load config
	{
		var t struct {
			C Config `yaml:".config"`
		}
		if err = unmarshal(&t); err != nil {
			return err
		}
		switch strings.ToLower(strings.TrimSpace(t.C.Mode)) {
		case "mysql":
			sq.Config.Mode = "mysql"
		case "postgresql", "postgres":
			sq.Config.Mode = "mysql"
		case "sqlite", "sqlite2", "sqlite3":
			sq.Config.Mode = "sqlite"
		default:
			sq.Config.Mode = "sqlite"
		}
		sq.Config.Package = strings.ToLower(strings.TrimSpace(t.C.Package))
		if sq.Config.Package == "" {
			sq.Config.Package = "mdl"
		}
		sq.Config.Tag = strings.ToLower(strings.TrimSpace(t.C.Tag))
		if sq.Config.Tag == "" {
			sq.Config.Tag = "mdl"
		}
		sq.Config.Dir = strings.TrimSpace(t.C.Dir)
		if sq.Config.Dir == "" {
			sq.Config.Dir = "mdl"
		}
	}

	sq.Table = []*Table{}
	for _, v := range ms {
		switch key, _ := v.Key.(string); key {
		case ".createTable":
			sq.CreateTableTemplate, _ = v.Value.(string)
		case ".dropTable":
			sq.DropTableTemplate, _ = v.Value.(string)
		case ".source":
			sq.SourceTemplate, _ = v.Value.(string)
		default:
			if key[0] == '.' {
				continue // ignored
			}

			ms2, ok := v.Value.(yaml.MapSlice)
			if !ok {
				return fmt.Errorf("an error occurred during parsing %q in yaml", key)
			}
			t, err := sq.parseTable(key, ms2)
			if err != nil {
				return err
			}
			sq.Table = append(sq.Table, t)
		}
	}

	for k := range sq.Table {
		if err = sq.markPrimaryKeys(sq.Table[k]); err != nil {
			return err
		}
		if err = sq.markForeignKeys(sq.Table[k]); err != nil {
			return err
		}
	}

	return nil
}

func (sq *Sqruct) parseTable(name string, ms yaml.MapSlice) (*Table, error) {
	t := &Table{
		parent:      sq,
		omitMethod:  map[string]struct{}{},
		GoName:      name,
		ColumnAfter: []string{},
	}
	for _, v := range ms {
		switch key, _ := v.Key.(string); key {
		case ".createTable":
			t.CreateTableTemplate, _ = v.Value.(string)
		case ".dropTable":
			t.DropTableTemplate, _ = v.Value.(string)
		case ".source":
			t.SourceTemplate, _ = v.Value.(string)
		case ".after":
			switch v := v.Value.(type) {
			case string:
				t.ColumnAfter = append(t.ColumnAfter, v)
			case []interface{}:
				for _, is := range v {
					s, _ := is.(string)
					t.ColumnAfter = append(t.ColumnAfter, s)
				}
			}
		case ".omit":
			switch v := v.Value.(type) {
			case string:
				for _, s := range strings.Split(v, ",") {
					t.omitMethod[strings.TrimSpace(s)] = struct{}{}
				}
			case []interface{}:
				for _, is := range v {
					s, _ := is.(string)
					t.omitMethod[strings.TrimSpace(s)] = struct{}{}
				}
			}
		default:
			if key[0] == '.' {
				continue // ignored
			}

			s, _ := v.Value.(string)
			idx := strings.Index(s, "|")
			t.Column = append(t.Column, &Column{
				parent:        t,
				GoName:        key,
				GoStructField: strings.TrimSpace(s[:idx]),
				SQLColumn:     strings.TrimSpace(s[idx+1:]),
			})
		}
	}
	return t, nil
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if s[0] == '`' || s[0] == '"' || s[0] == '\'' {
		// TODO(oov): need some intelligent processing
		return s[1 : len(s)-1]
	}
	return s
}

func splitUnquote(s string) []string {
	r := strings.Split(s, ",")
	for i, s := range r {
		r[i] = unquote(s)
	}
	return r
}

// TableByName finds table by name.
func (sq *Sqruct) TableByName(s string) *Table {
	s = strings.ToLower(s)
	for _, v := range sq.Table {
		if v.SQLName() == s {
			return v
		}
	}
	return nil
}

const namePattern = `([^'"()\x60]+|"[^"]+"|'[^']+'|\x60[^\x60]+\x60)` // \x60 = `

func (sq *Sqruct) markPrimaryKeys(t *Table) error {
	t.PrimaryKey.Column = []*Column{}

	// find from column constraint
	columnRE := regexp.MustCompile(`(?i)PRIMARY\s+KEY`)
	for k, c := range t.Column {
		if columnRE.MatchString(c.SQLColumn) {
			t.PrimaryKey.Column = []*Column{t.Column[k]}
			return nil
		}
	}

	// find from table constraint
	afterRE := regexp.MustCompile(`(?i)PRIMARY\s+KEY\s*\(` + namePattern + `\)`)
	for _, a := range t.ColumnAfter {
		m := afterRE.FindStringSubmatch(a)
		if len(m) == 0 {
			continue
		}
		for _, c := range strings.Split(m[1], ",") {
			c = unquote(c)
			col := t.ColumnByName(c)
			if col == nil {
				return fmt.Errorf(`primary key column %q is not found in table %q`, c, t.SQLName())
			}
			t.PrimaryKey.Column = append(t.PrimaryKey.Column, col)
		}
	}
	return nil
}

func (sq *Sqruct) markForeignKeys(t *Table) error {
	const ref = `REFERENCES\s+` + namePattern + `\s*\(` + namePattern + `\)`

	t.ForeignKey = []ForeignKey{}

	// find from column constraint
	columnRE := regexp.MustCompile(`(?i)\s+` + ref)
	for _, c := range t.Column {
		m := columnRE.FindStringSubmatch(c.SQLColumn)
		if len(m) == 0 {
			continue
		}

		oTable, oCol := unquote(m[1]), unquote(m[2])
		err := sq.registerForeignKey(t, oTable, []string{c.SQLName()}, []string{oCol})
		if err != nil {
			return err
		}
	}

	// find from table constraint
	afterRE := regexp.MustCompile(`(?i)FOREIGN\s+KEY\s*\(` + namePattern + `\)\s*` + ref)
	for _, a := range t.ColumnAfter {
		m := afterRE.FindStringSubmatch(a)
		if len(m) == 0 {
			continue
		}
		cols, oTable, oCols := splitUnquote(m[1]), unquote(m[2]), splitUnquote(m[3])
		err := sq.registerForeignKey(t, oTable, cols, oCols)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sq *Sqruct) registerForeignKey(t *Table, otherTable string, cols []string, otherCols []string) error {
	ot := sq.TableByName(otherTable)
	if ot == nil {
		return fmt.Errorf(`foreign key table %q is not found`, otherTable)
	}
	if len(cols) != len(otherCols) {
		return fmt.Errorf(`number of columns are mismatched`)
	}

	fk := ForeignKey{Table: ot}
	var p ColumnPair
	for i := range cols {
		p = ColumnPair{
			Self:  t.ColumnByName(cols[i]),
			Other: fk.Table.ColumnByName(otherCols[i]),
		}
		if p.Self == nil {
			return fmt.Errorf(`column %q is not found in table %q`, cols[i], t.SQLName())
		}
		if p.Other == nil {
			return fmt.Errorf(`column %q is not found in table %q`, otherCols[i], fk.Table.SQLName())
		}
		fk.Column = append(fk.Column, p)
	}

	t.ForeignKey = append(t.ForeignKey, fk)

	rfk := ForeignKey{Table: t, Mirror: true}
	for _, c := range fk.Column {
		rfk.Column = append(rfk.Column, ColumnPair{Self: c.Other, Other: c.Self})
	}
	ot.ForeignKey = append(ot.ForeignKey, rfk)
	return nil
}
