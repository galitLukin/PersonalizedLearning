package main

import (
	"database/sql"
	"fmt"
)

type userSchema struct {
	FName    string
	LName    string
	Email    string
	Userid   string
	Password string
}

func dbGetUsers(db *sql.DB) string {
	rows, err := db.Query(`SELECT * FROM users;`)
	dbCheck(err)
	defer rows.Close()

	var s string

	// query
	for rows.Next() {
		u := user{}
		err = rows.Scan(&u.UserName, &u.First, &u.Last, &u.UserName, &u.Password)
		dbCheck(err)
		s += fmt.Sprintf(`email: "%s" firstName: "%s", lastName: "%s", passWord: "%s"`, u.UserName, u.First, u.Last, u.Password)
		s += "\n"
	}
	return s
}

func dbGetUser(db *sql.DB, email string) user {
	q := fmt.Sprintf(`SELECT * FROM users WHERE email="%s";`, email)
	fmt.Println(q)
	rows, err := db.Query(q)
	dbCheck(err)
	defer rows.Close()

	// data to be used in query
	var u user
	for rows.Next() {
		err = rows.Scan(&u.UserName, &u.First, &u.Last, &u.UserName, &u.Password)
		dbCheck(err)
		s := fmt.Sprintf(`email: "%s" firstName: "%s", lastName: "%s", passWord: "%s"`, u.UserName, u.First, u.Last, u.Password)
		fmt.Printf(`RETRIEVED USER: %#v`, s)
	}
	return u
}

func dbCreateUser(db *sql.DB, newUser user) string {
	q := fmt.Sprintf(`insert into test02.users (fName, lName, email, password) values ("%s", "%s", "%s", "%s");`, newUser.First, newUser.Last, newUser.UserName, newUser.Password)
	stmt, err := db.Prepare(q)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "INSERTED RECORD ", n)
}

func dbDeleteUser(db *sql.DB) string {
	stmt, err := db.Prepare(`DELETE FROM test02.users WHERE fName="Omer";`)
	dbCheck(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	dbCheck(err)

	n, err := r.RowsAffected()
	dbCheck(err)

	return fmt.Sprintf("%s%d", "DELETED RECORD ", n)
}

func dbCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
