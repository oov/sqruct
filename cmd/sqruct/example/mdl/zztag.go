// DO NOT EDIT. This file was auto-generated by Sqruct.

package mdl

import (
	"database/sql"

	"github.com/oov/sqruct"
)

// Tag represents the following table.
// 	CREATE TABLE tag(
// 		id INTEGER PRIMARY KEY,
// 		name VARCHAR(30) NOT NULL UNIQUE
// 	);
type Tag struct {
	ID   int64  `mdl:"pk,notnull,uniq,default,autoincr"`
	Name string `mdl:"notnull,uniq"`
}

func GetTag(db sqruct.DB, id int64) (*Tag, error) {

	var t Tag
	err := db.QueryRow(
		"SELECT id, name FROM tag WHERE (id = ?)",
		id,
	).Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}
	return &t, nil

}

func (t *Tag) SelectPostTag(db sqruct.DB) ([]PostTag, error) {

	r, err := db.Query(
		"SELECT postid, tagid FROM posttag WHERE (tagid = ?)",
		t.ID,
	)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ot []PostTag
	for r.Next() {
		var e PostTag
		if err = r.Scan(&e.PostID, &e.TagID); err != nil {
			return nil, err
		}
		ot = append(ot, e)
	}
	if err = r.Err(); err != nil {
		return nil, err
	}
	if ot == nil {
		return nil, sql.ErrNoRows
	}
	return ot, nil

}

func (t *Tag) TableName() string {
	return "tag"
}

func (t *Tag) Columns() []string {
	return []string{"id", "name"}
}

func (t *Tag) Values() []interface{} {
	return []interface{}{t.ID, t.Name}
}

func (t *Tag) AutoIncrementColumnIndex() int {
	return 0
}

func (t *Tag) SqructMode() sqruct.Mode {
	return sqruct.SQLite
}

func (t *Tag) Insert(db sqruct.DB) error {

	i, err := t.SqructMode().Insert(db, t.TableName(), t.Columns(), t.Values(), t.AutoIncrementColumnIndex())
	if err != nil {
		return err
	}
	if i != 0 {
		t.ID = i
	}
	return nil

}

func (t *Tag) Update(db sqruct.DB) error {

	_, err := db.Exec(
		"UPDATE tag SET name = ? WHERE (id = ?)",
		t.Name,
		t.ID,
	)
	return err

}

func (t *Tag) Delete(db sqruct.DB) error {

	_, err := db.Exec(
		"DELETE FROM tag WHERE (id = ?)",
		t.ID,
	)
	return err

}
