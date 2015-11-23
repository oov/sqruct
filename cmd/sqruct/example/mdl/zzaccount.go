// DO NOT EDIT. This file was auto-generated by Sqruct.

package mdl

import (
	"database/sql"

	"github.com/oov/sqruct"
)

// Account represents the following table.
// 	CREATE TABLE account(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name VARCHAR(30) NOT NULL UNIQUE
// 	);
type Account struct {
	ID   int64  `mdl:"pk,notnull,uniq,default,autoincr"`
	Name string `mdl:"notnull,uniq"`
}

func GetAccount(db sqruct.DB, id int64) (*Account, error) {

	r, err := db.Query(
		"SELECT id, name FROM account WHERE (id = ?)",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if !r.Next() {
		if err = r.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var t Account
	if err = r.Scan(&t.ID, &t.Name); err != nil {
		return nil, err
	}

	return &t, nil

}

func (t *Account) SelectPost(db sqruct.DB) ([]Post, error) {

	r, err := db.Query(
		"SELECT id, accountid, at, message FROM post WHERE (accountid = ?)",
		t.ID,
	)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ot []Post
	for r.Next() {
		var e Post
		if err = r.Scan(&e.ID, &e.AccountID, &e.At, &e.Message); err != nil {
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

func (t *Account) TableName() string {
	return "account"
}

func (t *Account) Columns() []string {
	return []string{"id", "name"}
}

func (t *Account) Values() []interface{} {
	return []interface{}{t.ID, t.Name}
}

func (t *Account) AutoIncrementColumnIndex() int {
	return 0
}

func (t *Account) Insert(db sqruct.DB) error {

	i, err := sqruct.SQLite.Insert(db, t.TableName(), t.Columns(), t.Values(), t.AutoIncrementColumnIndex())
	if err != nil {
		return err
	}
	if i != 0 {
		t.ID = i
	}
	return nil

}

func (t *Account) Update(db sqruct.DB) error {

	_, err := db.Exec(
		"UPDATE account SET name = ? WHERE (id = ?)",
		t.Name,
		t.ID,
	)
	return err

}

func (t *Account) Delete(db sqruct.DB) error {

	_, err := db.Exec(
		"DELETE FROM account WHERE (id = ?)",
		t.ID,
	)
	return err

}
