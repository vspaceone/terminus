package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	"crypto/rand"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user stored in the database
type User struct {
	userid   int32
	username string // first name
	fullname string //last name
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
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT user.userid, username, fullname FROM user, tags WHERE user.userid = tags.userid AND tags.taguid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var userid int32
	var username string
	var fullname string
	err = stmtOut.QueryRow(uid).Scan(&userid, &username, &fullname)
	check(err)
	fmt.Println("Scanned ", userid, "", username, " ", fullname)
	var ret = User{userid, username, fullname}

	return ret
}

func getUserIDByUID(uid string) int32 {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT tags.userid FROM tags WHERE tags.taguid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var userid int32
	err = stmtOut.QueryRow(uid).Scan(&userid)
	check(err)
	fmt.Println("Scanned ", userid)

	return userid
}

func getUserByUserID(userid int32) User {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT userid, username, fullname FROM user WHERE userid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	//var userid int32
	var username string
	var fullname string
	err = stmtOut.QueryRow(userid).Scan(&userid, &username, &fullname)
	check(err)
	fmt.Println("Scanned ", userid, "", username, " ", fullname)

	var ret = User{userid, username, fullname}

	return ret
}

func getUserByUsername(username string) User {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT userid, username, fullname FROM user WHERE userid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var userid int32
	//var username string
	var fullname string
	err = stmtOut.QueryRow(userid).Scan(&userid, &username, &fullname)
	check(err)
	fmt.Println("Scanned ", userid, " ", username, " ", fullname)

	var ret = User{userid, username, fullname}

	return ret
}

func compareUserAuthByUID(uid string, password string) bool {
	// get hashed pw and salt
	// hash password parameter with salt by PBKDF2 method
	// compare hashes
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT password, salt FROM passwords, user, tags WHERE user.userid = passwords.userid AND user.userid = tags.userid AND tags.taguid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var dbPassword [256]byte
	var salt [256]byte
	err = stmtOut.QueryRow(uid).Scan(&dbPassword, &salt)
	check(err)
	fmt.Println("Scanned ", dbPassword, " ", salt)

	newHashedPw := pbkdf2.Key([]byte(password), salt[:], 4096, 256, sha1.New)

	return reflect.DeepEqual(dbPassword[:], newHashedPw)
}

func doesUIDExist(uid string) bool {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT userid FROM tags WHERE uid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	rows, err := stmtOut.Query(uid)
	check(err)

	var i int32
	for i = 0; rows.Next(); i++ {
	}

	return i > 0
}

func doesUserIDExist(userid int32) bool {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT * FROM user WHERE userid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	rows, err := stmtOut.Query(userid)
	check(err)

	var i int32
	for i = 0; rows.Next(); i++ {
	}

	return i > 0
}

func doesUsernameExist(username string) bool {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT * FROM user WHERE username = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	rows, err := stmtOut.Query(username)
	check(err)

	var i int32
	for i = 0; rows.Next(); i++ {
	}

	return i > 0
}

func createNewUser(uid, username, fullname, password string) bool {
	//assign tag uid to user, and set auth
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	if doesUsernameExist(username) {
		return false
	}

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("INSERT INTO user (username, fullname) VALUES(?, ?)")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	res, err := stmtOut.Exec(username, fullname)
	check(err)
	liid, err := res.LastInsertId()
	ra, err := res.RowsAffected()
	fmt.Println("Created user ", username, " ", fullname, " ", liid, " ", ra)

	assignUIDToUser(uid, int32(liid))
	assignPasswordToUser(password, int32(liid))

	if err != nil {
		return true
	}
	return false
}

func assignUIDToUser(uid string, userid int32) bool {
	//assign tag uid to user, and set auth
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	if doesUIDExist(uid) || !doesUserIDExist(userid) {
		return false
	}

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("INSERT INTO tags (taguid, userid) VALUES(?, ?)")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	res, err := stmtOut.Exec(uid, userid)
	check(err)
	liid, err := res.LastInsertId()
	ra, err := res.RowsAffected()
	fmt.Println("Assigned uid ", uid, " to ", userid, " ", liid, " ", ra)

	if err != nil {
		return true
	}
	return false

}

func assignPasswordToUser(password string, userid int32) bool {
	//assign tag uid to user, and set auth
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	if !doesUserIDExist(userid) {
		return false
	}

	var salt = make([]byte, 256)
	i, err := rand.Read(salt)
	newHashedPw := pbkdf2.Key([]byte(password), salt[:], 4096, 256, sha1.New)

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("REPLACE INTO passwords (userid, password, salt) VALUES (?, ?, ?)")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	res, err := stmtOut.Exec(userid, newHashedPw, salt)
	check(err)
	liid, err := res.LastInsertId()
	ra, err := res.RowsAffected()
	fmt.Println("Assigned pw ", password, " as ", newHashedPw, " to ", userid, " with salt ", salt, " len ", i, " ", liid, " ", ra)

	if err != nil {
		return true
	}
	return false

}

func checkPasswordOnUserID(userid int32, password string) bool {
	//assign tag uid to user, and set auth
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
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
	stmtOut, err := db.Prepare("SELECT userid, password, salt FROM passwords WHERE userid = ?")
	check(err)
	defer stmtOut.Close()

	// Queries the statement and returns the answer
	var dbuserid int32
	var dbpassword = make([]byte, 256)
	var salt = make([]byte, 256)
	err = stmtOut.QueryRow(userid).Scan(&dbuserid, &dbpassword, &salt)
	check(err)
	fmt.Println("Scanned ", dbuserid, " ", dbpassword, " ", salt)

	newHashedPw := pbkdf2.Key([]byte(password), salt, 4096, 256, sha1.New)

	return reflect.DeepEqual(newHashedPw, dbpassword)

}
