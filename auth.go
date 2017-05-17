package main

import (
	"fmt"
	"time"
)

type Authenticator struct {
	time  int64
	uid   string
	token string
}

var authenticatorSessions = make(map[string]Authenticator)

// Timeout for session in seconds
var sessionTimeout int64 = 60

func newAuthenticatorSession(uid string) string {
	token := genToken()
	authenticatorSessions[uid] = Authenticator{time.Now().Unix(), uid, token}
	return token
}

func sessionExists(uid string) bool {
	_, ok := authenticatorSessions[uid]
	return ok
}

func verifyToken(uid, token string) bool {
	fmt.Println((time.Now().Unix() - authenticatorSessions[uid].time))
	return authenticatorSessions[uid].token == token && (time.Now().Unix()-authenticatorSessions[uid].time) < sessionTimeout
}

func genToken() string {
	str, _ := GenerateRandomString(64)
	return str
}
