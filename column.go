package sqruct

import (
	"regexp"
	"strings"
)

// Column represents database column.
type Column struct {
	parent        *Table
	GoName        string // Column name in Go
	GoStructField string // Struct field definition text in Go
	SQLColumn     string // Column definition text in RDBMS including column constraint statement
}

// SQLName returns column name in RDBMS.
func (c *Column) SQLName() string {
	return strings.ToLower(c.GoName)
}

// GoStructFieldWithTag returns struct field definition text in Go including some tag data, such as
//	ID int64 `mdl:"pk,notnull,uniq,default,autoincr"`
func (c *Column) GoStructFieldWithTag() string {
	tags := []string{}
	if c.PrimaryKey() {
		tags = append(tags, "pk")
	}
	if fk := c.ForeignKey(); fk != nil && !fk.Mirror && len(fk.Column) == 1 {
		tags = append(tags, "fk")
	}
	if c.NotNull() {
		tags = append(tags, "notnull")
	}
	if c.Unique() {
		tags = append(tags, "uniq")
	}
	if c.Default() {
		tags = append(tags, "default")
	}
	if c.AutoIncrement() {
		tags = append(tags, "autoincr")
	}
	s, err := replaceStructTag(c.GoStructField, func(s string) (string, error) {
		// TODO(oov): merge tag data
		if s != "" {
			s += " "
		}
		return s + c.parent.parent.Config.Tag + `:"` + strings.Join(tags, ",") + `"`, nil
	})
	if err != nil {
		panic(err)
	}
	return s
}

// PrimaryKey reports whether this column is primary key.
// if this column is the part of composite primary key, PrimaryKey returns false.
func (c *Column) PrimaryKey() bool {
	return len(c.parent.PrimaryKey.Column) == 1 && c.SQLName() == c.parent.PrimaryKey.Column[0].SQLName()
}

// ForeignKey returns foreign key mapping data.
// if this column isn't a foreign key, returns nil.
func (c *Column) ForeignKey() *ForeignKey {
	for _, fk := range c.parent.ForeignKey {
		for _, cp := range fk.Column {
			if cp.Self == c {
				return &fk
			}
		}
	}
	return nil
}

// GoStructFieldType returns only type from Go struct field definition text.
func (c *Column) GoStructFieldType() string {
	s, err := extractStructFieldType(c.GoStructField)
	if err != nil {
		panic(err)
	}
	return s
}

// NotNull reports whether this column has not null constraint.
func (c *Column) NotNull() bool {
	return regexp.MustCompile(`(?i)NOT\s+NULL|\S+\s+PRIMARY\s+KEY`).MatchString(c.SQLColumn)
}

// Unique reports whether this column has unique constraint.
func (c *Column) Unique() bool {
	return regexp.MustCompile(`(?i)\S+\s+UNIQUE|\S+\s+PRIMARY\s+KEY`).MatchString(c.SQLColumn)
}

// AutoIncrement reports whether this column has auto increment constraint.
func (c *Column) AutoIncrement() bool {
	r := map[string]string{
		"mysql":      `(?i)\S+\s+AUTO_?INCREMENT`,
		"postgresql": `(?i)\S+\s+SERIAL|\S+\s+DEFAULT\s+nextval`,
		"sqlite":     `(?i)\S+\s+AUTO_?INCREMENT|INTEGER\s+PRIMARY\s+KEY`,
	}
	return regexp.MustCompile(r[c.parent.parent.Config.Mode]).MatchString(c.SQLColumn)
}

// Default reports whether this column has default constraint.
func (c *Column) Default() bool {
	return c.AutoIncrement() || regexp.MustCompile(`(?i)\S+\s+DEFAULT\s+\S+`).MatchString(c.SQLColumn)
}
