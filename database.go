package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user stored in the database
type User struct {
	firstName string // first name
	lastName  string //last name
}

// getUserByUID gets the User entry in the database by his Tag's UID
func getUserByUID(uid string) User {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf, err := getConfiguration()
	conf := inf.(map[string]interface{})
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", conf["dbUser"], conf["dbPassw"], conf["dbHost"], conf["dbName"])

	// Opens the connection
	db, err := sql.Open("mysql", dbStr)
	check(err)
	defer db.Close()

	// Pings to go sure that connection works
	err = db.Ping()
	check(err)

	// Prepares statement
	stmtOut, err := db.Prepare("SELECT fname, lname FROM user, tags WHERE user.userid = tags.userid AND tags.taguid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var fname string
	var lname string
	err = stmtOut.QueryRow(uid).Scan(&fname, &lname)
	check(err)
	fmt.Println("Scanned ", fname, " ", lname)
	var ret = User{fname, lname}

	return ret
}

func compareUserAuthByUID(uid string, password string) bool {
	// get hashed pw and salt
	// hash password parameter with salt by PBKDF2 method
	// compare hashes
}
