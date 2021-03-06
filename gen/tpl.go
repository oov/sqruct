package gen

const createTableTemplate = `CREATE TABLE {{.Name.SQLQuoted}}(
	{{range $k, $v := .Column}}{{if $k}},
	{{end}}{{$v.Name.SQLQuoted}} {{$v.SQLColumn}}{{end}}{{range $k, $v := .ColumnAfter}},
	{{$v}}{{end}}
);`

const dropTableTemplate = `DROP TABLE {{.Name.SQL}};`

const sourceTemplate = `

// DO NOT EDIT. This file was auto-generated by Sqruct.

package {{.PackageName}}

import (
	"github.com/oov/q"
	"github.com/oov/sqruct"
)

// {{.Name.Go}} represents the following table.
{{.MustCreateTableSQL.AddPrefix "// \t"}}
type {{.Name.Go}} struct {
{{range $k, $v := .Column}}  {{$v.Name.Go}} {{$v.GoStructFieldWithTag}}
{{end}}}

{{$method := print "Get" .Name.Go}}
{{if .OmitMethod $method}}/*{{end}}
func {{$method}}(db sqruct.DB{{range $k, $v := .PrimaryKey.Column}}, {{$v.Name.GoLower}} {{$v.GoStructFieldType}}{{end}}) (*{{.Name.Go}}, error) {
	b, tbl := zz{{.Name.Go}}{}.SelectBuilder()
	sql, args := b.Where(
		{{range $k, $v := .PrimaryKey.Column}}q.Eq(tbl.{{$v.Name.Go}}(), {{$v.Name.GoLower}}),
		{{end}}
	).ToSQL()
	var t {{.Name.Go}}
  err := db.QueryRow(sql, args...).Scan(zz{{.Name.Go}}{}.Pointers(&t)...)
	if err != nil {
  	return nil, err
  }
  return &t, nil
}
{{if .OmitMethod $method}}*/{{end}}

{{$t := .}}
{{range $_, $fk := .ForeignKey}}
{{if and (eq (len $fk.Column) 1) (index $fk.Column 0).Other.PrimaryKey}}
{{$method := print "Get" $fk.Table.Name.Go}}
{{if $t.OmitMethod $method}}/*{{end}}
func (t *{{$t.Name.Go}}) {{$method}}(db sqruct.DB) (*{{$fk.Table.Name.Go}}, error) {
	b, tbl := zz{{$fk.Table.Name.Go}}{}.SelectBuilder()
	sql, args := b.Where(
		{{range $k, $v := $fk.Column}}q.Eq(tbl.{{$v.Other.Name.Go}}(), t.{{$v.Self.Name.Go}}),
		{{end}}
	).ToSQL()
	var ot {{$fk.Table.Name.Go}}
  if err := db.QueryRow(sql, args...).Scan(zz{{$fk.Table.Name.Go}}{}.Pointers(&ot)...); err != nil {
  	return nil, err
  }
  return &ot, nil
}
{{if $t.OmitMethod $method}}*/{{end}}
{{else}}
{{$method := print "Select" $fk.Table.Name.Go}}
{{if $t.OmitMethod $method}}/*{{end}}
func (t *{{$t.Name.Go}}) {{$method}}(db sqruct.DB) ([]{{$fk.Table.Name.Go}}, error) {
	b, tbl := zz{{$fk.Table.Name.Go}}{}.SelectBuilder()
	sql, args := b.Where(
		{{range $k, $v := $fk.Column}}q.Eq(tbl.{{$v.Other.Name.Go}}(), t.{{$v.Self.Name.Go}}),
		{{end}}
	).ToSQL()
	r, err := db.Query(sql, args...)
  if err != nil {
  	return nil, err
  }
	defer r.Close()

	ot := []{{$fk.Table.Name.Go}}{}
	for r.Next() {
		var e {{$fk.Table.Name.Go}}
		if err = r.Scan(zz{{$fk.Table.Name.Go}}{}.Pointers(&e)...); err != nil {
  		return nil, err
  	}
		ot = append(ot, e)
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	return ot, nil
}
{{if $t.OmitMethod $method}}*/{{end}}
{{end}}
{{end}}

{{$t := .}}
{{range $_, $m2m := .ManyToMany}}
{{$relTable := $m2m.RelTable}}
{{$oTable := $m2m.OtherFK.Table}}
{{$method := print "Select" $oTable.Name.Go}}
{{if $t.OmitMethod $method}}/*{{end}}
func (t *{{$t.Name.Go}}) {{$method}}(db sqruct.DB) ([]{{$oTable.Name.Go}}, []{{$relTable.Name.Go}}, error) {
	b, relTbl, _ := zz{{$t.Name.Go}}{}.SelectBuilderFor{{$oTable.Name.Go}}()
	sql, args := b.Where(
		{{range $k, $v := $m2m.MyFK.Column}}q.Eq(relTbl.{{$v.Self.Name.Go}}(), t.{{$v.Other.Name.Go}}),
		{{end}}
	).ToSQL()
	r, err := db.Query(sql, args...)
  if err != nil {
  	return nil, nil, err
  }
	defer r.Close()

	ot, rt := []{{$oTable.Name.Go}}{}, []{{$relTable.Name.Go}}{}
	for r.Next() {
		var oe {{$oTable.Name.Go}}
		var re {{$relTable.Name.Go}}
		if err = r.Scan(append(zz{{$relTable.Name.Go}}{}.Pointers(&re), zz{{$oTable.Name.Go}}{}.Pointers(&oe)...)...); err != nil {
  		return nil, nil, err
  	}
		ot, rt = append(ot, oe), append(rt, re)
	}
	if err = r.Err(); err != nil {
		return nil, nil, err
	}
	return ot, rt, nil
}
{{if $t.OmitMethod $method}}*/{{end}}
{{end}}

{{$method := "Insert"}}
{{if .OmitMethod $method}}/*{{end}}
func (t *{{.Name.Go}}) {{$method}}(db sqruct.DB) error {
	{{$aicol := .AutoIncrementColumn}}
	{{if $aicol}}
		b, tbl := zz{{.Name.Go}}{}.InsertBuilder(t)
		if !sqruct.IsZero(t.{{$aicol.Name.Go}}) {
			sql, args := b.Set(tbl.{{$aicol.Name.Go}}(), t.{{$aicol.Name.Go}}).ToSQL()
			_, err := db.Exec(sql, args...)
			return err
		}

		{{if eq .Mode.String "PostgreSQL"}}
			sql, args := b.Returning(tbl.{{$aicol.Name.Go}}()).ToSQL()
			var i int64
			if err := db.QueryRow(sql, args...).Scan(&i); err != nil {
				return err
			}
		{{else}}
			sql, args := b.ToSQL()
			r, err := db.Exec(sql, args...)
			if err != nil {
				return err
			}
			var i int64
			if i, err = r.LastInsertId(); err != nil {
				return err
			}
		{{end}}
		t.{{$aicol.Name.Go}} = {{if eq $aicol.GoStructFieldType "int64"}}i{{else}}{{$aicol.GoStructFieldType}}(i){{end}}
		return nil
	{{else}}
		b, _ := zz{{.Name.Go}}{}.InsertBuilder(t)
		sql, args := b.ToSQL()
		_, err := db.Exec(sql, args...)
		return err
	{{end}}
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "Update"}}
{{if .OmitMethod $method}}/*{{end}}
{{if .NonPrimaryKeys}}
func (t *{{.Name.Go}}) {{$method}}(db sqruct.DB) error {
	b, tbl := zz{{.Name.Go}}{}.UpdateBuilder(t)
	sql, args := b.Where(
		{{range $_, $v := .PrimaryKey.Column}}q.Eq(tbl.{{$v.Name.Go}}(), t.{{$v.Name.Go}}),
		{{end}}
	).ToSQL()
	_, err := db.Exec(sql, args...)
	return err
}
{{else}}
func (t *{{.Name.Go}}) {{$method}}(db sqruct.DB) error {
	// {{.Name.Go}} has primary key only
	return nil
}
{{end}}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "Delete"}}
{{if .OmitMethod $method}}/*{{end}}
func (t *{{.Name.Go}}) {{$method}}(db sqruct.DB) error {
	b, tbl := zz{{.Name.Go}}{}.DeleteBuilder()
	sql, args := b.Where(
		{{range $_, $v := .PrimaryKey.Column}}q.Eq(tbl.{{$v.Name.Go}}(), t.{{$v.Name.Go}}),
		{{end}}
	).ToSQL()
	_, err := db.Exec(sql, args...)
	return err
}
{{if .OmitMethod $method}}*/{{end}}

// zz{{.Name.Go}} represents {{.Name.Go}} table schema.
type zz{{.Name.Go}} struct {}

{{$method := "T"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}(aliasName ...string) *zz{{.Name.Go}}Table {
	return &zz{{.Name.Go}}Table{q.T({{.Name.SQLQuoted}}, aliasName...)}
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "Columns"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}(b *q.ZSelectBuilder, t *zz{{.Name.Go}}Table) {
	b.Column(
		{{range $_, $v := .Column}}t.{{$v.Name.Go}}(),
		{{end}}
	)
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "Pointers"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}(t *{{.Name.Go}}) []interface{} {
	return []interface{}{ {{range $k, $v := .Column}}{{if $k}},{{end}}&t.{{$v.Name.Go}}{{end}} }
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "InsertBuilder"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}(t *{{.Name.Go}}) (*q.ZInsertBuilder, *zz{{.Name.Go}}Table) {
	tbl := zz{{.Name.Go}}{}.T()
	return q.Insert().Into(tbl).
		{{range $k, $v := .Column}}
			{{if not $v.AutoIncrement}}
				Set(tbl.{{$v.Name.Go}}(), t.{{$v.Name.Go}}).
			{{end}}
		{{end}}
		SetDialect(q.{{.Mode}}), tbl
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "SelectBuilder"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}() (*q.ZSelectBuilder, *zz{{.Name.Go}}Table) {
	tbl := zz{{.Name.Go}}{}.T()
	b := q.Select().From(tbl).SetDialect(q.{{.Mode}})
	zz{{.Name.Go}}{}.Columns(b, tbl)
	return b, tbl
}
{{if .OmitMethod $method}}*/{{end}}

{{$t := .}}
{{range $_, $m2m := .ManyToMany}}
{{$relTable := $m2m.RelTable}}
{{$oTable := $m2m.OtherFK.Table}}
{{$method := print "SelectBuilderFor" $oTable.Name.Go}}
{{if $t.OmitMethod $method}}/*{{end}}
func (zz{{$t.Name.Go}}) {{$method}}() (b *q.ZSelectBuilder, {{$relTable.Name.GoLower}} *zz{{$relTable.Name.Go}}Table, {{$oTable.Name.GoLower}} *zz{{$oTable.Name.Go}}Table) {
	b, relTbl := zz{{$relTable.Name.Go}}{}.SelectBuilder()
	oTbl := zz{{$oTable.Name.Go}}{}.T()
	relTbl.InnerJoin(
		oTbl,
		{{range $_, $v := $m2m.OtherFK.Column}}q.Eq(relTbl.{{$v.Self.Name.Go}}(), oTbl.{{$v.Other.Name.Go}}()),
		{{end}}
	)
	zz{{$oTable.Name.Go}}{}.Columns(b, oTbl)
	return b, relTbl, oTbl
}
{{if $t.OmitMethod $method}}*/{{end}}
{{end}}

{{$method := "UpdateBuilder"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}(t *{{.Name.Go}}) (*q.ZUpdateBuilder, *zz{{.Name.Go}}Table) {
	tbl := zz{{.Name.Go}}{}.T()
	return q.Update(tbl).
		{{range $k, $v := .NonPrimaryKeys}}
			Set(tbl.{{$v.Name.Go}}(), t.{{$v.Name.Go}}).
		{{end}}
		SetDialect(q.{{.Mode}}), tbl
}
{{if .OmitMethod $method}}*/{{end}}

{{$method := "DeleteBuilder"}}
{{if .OmitMethod $method}}/*{{end}}
func (zz{{.Name.Go}}) {{$method}}() (*q.ZDeleteBuilder, *zz{{.Name.Go}}Table) {
	tbl := zz{{.Name.Go}}{}.T()
	return q.Delete().From(tbl).SetDialect(q.{{.Mode}}), tbl
}
{{if .OmitMethod $method}}*/{{end}}

// zz{{.Name.Go}}Table represents {{.Name.Go}} table.
type zz{{.Name.Go}}Table struct { q.Table }

{{$t := .}}
{{range $k, $v := .Column}}
func (t zz{{$t.Name.Go}}Table) {{$v.Name.Go}}(aliasName ...string) q.Column {
	return t.Table.C({{$v.Name.SQLQuoted}}, aliasName...)
}
{{end}}
`
