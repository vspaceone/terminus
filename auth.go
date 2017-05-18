package main

import (
	"time"
)

type authenticator struct {
	time  int64
	uid   int32
	token string
}

var authenticatorSessions = make(map[int32]authenticator)

// Timeout for session in seconds
var sessionTimeout int64 = 60

func newAuthenticatorSession(userid int32) string {
	token := genToken()
	authenticatorSessions[userid] = authenticator{time.Now().Unix(), userid, token}
	return token
}

func sessionExists(userid int32) bool {
	_, ok := authenticatorSessions[userid]
	return ok
}

func verifyToken(userid int32, token string) bool {
	//fmt.Println((time.Now().Unix() - authenticatorSessions[uid].time))
	return authenticatorSessions[userid].token == token && (time.Now().Unix()-authenticatorSessions[userid].time) < sessionTimeout
}

func genToken() string {
	str, _ := GenerateRandomString(64)
	return str
}
