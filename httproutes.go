package main

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/coyove/jsonbuilder"
	"github.com/gorilla/mux"
)

//handler's response informs briefly about this service
func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Terminus Server alpha1\nSee docs or ask the devs about the REST API")

	return
}

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json := jsonbuilder.Object()
	switch vars["type"] {
	case "user":
		user := getUserByUID(r.FormValue("uid"))
		//For some reason jsonbuilder.From(user) does not work
		json.Set("status", "ok").Begin("data").Set("username", string(user.username)).Set("fullname", string(user.fullname)).End()

		if r.FormValue("token") != "" {
			if verifyToken(r.FormValue("uid"), r.FormValue("token")) {
				json.Set("authLevel", user.authlevel)
			} else {
				json.Set("authLevel", "none")
			}
		}

		fmt.Fprintf(w, json.MarshalPretty())
	default:
		json.Set("status", "unknown query")

	}
}

func new(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json := jsonbuilder.Object()
	switch vars["type"] {
	case "user":
		authlevel, _ := strconv.ParseInt(r.FormValue("authlevel"), 10, 64)
		res := createNewUser(r.FormValue("uid"), r.FormValue("username"), r.FormValue("fullname"), r.FormValue("password"), authlevel)
		json.Set("status", res)
		fmt.Fprintf(w, json.MarshalPretty())
	default:
		json.Set("status", "unknown query")

	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	json := jsonbuilder.Object()
	result := checkPasswordOnUserID(getUserIDByUID(r.FormValue("uid")), r.FormValue("password"))
	//For some reason jsonbuilder.From(user) does not work
	if result {
		json.Set("status", "ok").Begin("data").Set("token", newAuthenticatorSession(r.FormValue("uid"))).End()
	} else {
		json.Set("status", "wrong password or userid")
	}

	fmt.Fprintf(w, json.MarshalPretty())
}
