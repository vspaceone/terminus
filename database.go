package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/pbkdf2"

	"crypto/rand"
	"reflect"

	"strings"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user stored in the database
type User struct {
	userid    int32
	username  string // first name
	fullname  string //last name
	authlevel int32
}

// getUserByUID gets the User entry in the database by his Tag's UID
func getUserByUID(uid string) User {
	var userid int32
	var username string
	var fullname string
	var authlevel int32

	getS("user, tags", "user.userid, username, fullname, authlevel", "user.userid = tags.userid AND tags.taguid = "+uid, &userid, &username, &fullname, &authlevel)

	return User{userid, username, fullname, authlevel}

}

func getUserIDByUID(uid string) int32 {
	var userid int32
	err := getS("tags", "userid", "taguid = "+uid, &userid)

	if err != nil {
		return -1
	}

	return userid
}

func getUserByUserID(userid int32) User {
	var username string
	var fullname string
	var authlevel int32
	getS("user", "userid, username, fullname, authlevel", "userid = "+strconv.FormatInt(int64(userid), 10), &userid, &username, &fullname, &authlevel)

	var ret = User{userid, username, fullname, authlevel}

	return ret
}

func getUserByUsername(username string) User {
	var userid int32
	var fullname string
	var authlevel int32
	getS("user", "userid, username, fullname, authlevel", "username = "+username, &userid, &username, &fullname, &authlevel)

	var ret = User{userid, username, fullname, authlevel}

	return ret
}

func compareUserAuthByUID(uid string, password string) bool {
	var dbPassword [256]byte
	var salt [256]byte
	getS("passwords, user, tags", "password, salt", "user.userid = passwords.userid AND user.userid = tags.userid AND tags.taguid = "+uid, &dbPassword, &salt)

	newHashedPw := pbkdf2.Key([]byte(password), salt[:], 4096, 256, sha1.New)

	return reflect.DeepEqual(dbPassword[:], newHashedPw)
}

func doesUIDExist(uid string) bool {
	var i int32
	getS("tags", "userid", "taguid = "+uid, &i)

	return i > 0
}

func doesUserIDExist(userid int32) bool {
	var i int32
	getS("user", "COUNT(*)", "userid = "+strconv.FormatInt(int64(userid), 10), &i)

	return i > 0
}

func doesUsernameExist(username string) bool {
	var i int32
	getS("user", "COUNT(*)", "username = "+username, &i)

	return i > 0
}

func createNewUser(uid, username, fullname, password string, authlevel int64) bool {
	if doesUsernameExist(username) {
		return false
	}

	vals := make([]interface{}, 3)
	vals[0] = username
	vals[1] = fullname
	vals[2] = authlevel

	lastID, _, err := put("user", []string{"username", "fullname", "authlevel"}, vals)

	if err != nil {
		return false
	}

	assgUID := assignUIDToUser(uid, int32(lastID))
	assgPW := assignPasswordToUser(password, int32(lastID))

	return assgPW && assgUID
}

func assignUIDToUser(uid string, userid int32) bool {
	if doesUIDExist(uid) || !doesUserIDExist(userid) {
		return false
	}

	data := make([]interface{}, 0)
	data = append(data, uid)
	data = append(data, userid)

	_, rowsAffected, err := put("tags", []string{"taguid", "userid"}, data)

	if err != nil {
		return false
	}
	return rowsAffected > 0
}

func assignPasswordToUser(password string, userid int32) bool {
	if !doesUserIDExist(userid) {
		return false
	}

	var salt = make([]byte, 256)
	_, err := rand.Read(salt)
	newHashedPw := pbkdf2.Key([]byte(password), salt[:], 4096, 256, sha1.New)

	data := make([]interface{}, 0)
	data = append(data, userid)
	data = append(data, newHashedPw)
	data = append(data, salt)

	_, rowsAffected, err := put("passwords", []string{"userid", "password", "salt"}, data)

	if err != nil {
		return false
	}
	return rowsAffected > 0
}

func checkPasswordOnUserID(userid int32, password string) bool {
	var dbuserid int32
	var dbpassword = make([]byte, 256)
	var salt = make([]byte, 256)

	getS("passwords", "userid, password, salt", "userid = "+strconv.FormatInt(int64(userid), 10), &dbuserid, &dbpassword, &salt)

	newHashedPw := pbkdf2.Key([]byte(password), salt, 4096, 256, sha1.New)

	return reflect.DeepEqual(newHashedPw, dbpassword)
}

// put() inserts a new row or replaces an existing one if the primary key is already used
func put(table string, keys []string, vals []interface{}) (lastID int64, rowsAffected int64, err error) {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	// Checking args
	if table == "" {
		return 0, 0, fmt.Errorf("table should be specified")
	}
	if keys == nil {
		return 0, 0, fmt.Errorf("keys should not be null")
	}
	if vals == nil {
		return 0, 0, fmt.Errorf("vals should not be null")
	}
	if len(vals) != len(keys) {
		return 0, 0, fmt.Errorf("vals and keys should have the same length")
	}

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
	conf := inf.(map[string]interface{})
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", conf["dbUser"], conf["dbPassw"], conf["dbHost"], conf["dbName"])

	// Opening the connection
	db, err := sql.Open("mysql", dbStr)
	check(err)
	defer db.Close()

	// Pings to go sure that connection works
	err = db.Ping()
	check(err)

	// Prepares statement
	stmt := "REPLACE INTO " + table
	stmt += "(" + strings.Join(keys, ",") + ") VALUES "
	qM := make([]string, 0)
	for i := 0; i < len(vals); i++ {
		qM = append(qM, "?")
	}
	stmt += "(" + strings.Join(qM, ",") + ")"
	fmt.Println(stmt)
	stmtOut, err := db.Prepare(stmt)
	check(err)
	defer stmtOut.Close()

	// Execute
	res, err := stmtOut.Exec(vals...)
	check(err)
	liid, err := res.LastInsertId()
	ra, err := res.RowsAffected()
	fmt.Println("Last insert id: " + strconv.FormatInt(liid, 10) + ", rows affected: " + strconv.FormatInt(ra, 10))

	return liid, ra, err
}

// getS() gets only the first row of a query
func getS(table, keys, where string, scanVars ...interface{}) error {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	// Checking args
	if table == "" {
		return fmt.Errorf("table should be specified")
	}
	if keys == "" {
		return fmt.Errorf("keys should not be null")
	}

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
	conf := inf.(map[string]interface{})
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", conf["dbUser"], conf["dbPassw"], conf["dbHost"], conf["dbName"])

	// Opening the connection
	db, err := sql.Open("mysql", dbStr)
	check(err)
	defer db.Close()

	// Pings to go sure that connection works
	err = db.Ping()
	check(err)

	// Prepares statement
	stmt := "SELECT " + keys + " FROM " + table + " WHERE " + where

	stmtOut, err := db.Prepare(stmt)
	check(err)
	defer stmtOut.Close()

	// Execute
	err = stmtOut.QueryRow().Scan(scanVars...)
	check(err)

	return err
}

// getM() gets multiple rows with a query
func getM(table, keys, where string) (*sql.Rows, error) {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println(rec)
		}
	}()

	// Checking args
	if table == "" {
		return nil, fmt.Errorf("table should be specified")
	}
	if keys == "" {
		return nil, fmt.Errorf("keys should not be null")
	}

	//Getting configuration for the MySQL connection
	inf := getConfiguration()
	conf := inf.(map[string]interface{})
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", conf["dbUser"], conf["dbPassw"], conf["dbHost"], conf["dbName"])

	// Opening the connection
	db, err := sql.Open("mysql", dbStr)
	check(err)
	defer db.Close()

	// Pings to go sure that connection works
	err = db.Ping()
	check(err)

	// Prepares statement
	stmt := "SELECT " + keys + " FROM " + table + " WHERE " + where

	stmtOut, err := db.Prepare(stmt)
	check(err)
	defer stmtOut.Close()

	// Execute
	row, err := stmtOut.Query()
	check(err)

	return row, err
}
