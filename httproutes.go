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

	fmt.Fprintln(w, "Terminus Server alpha1\nSee docs on https://github.com/vspaceone/terminus-server/tree/development or ask the devs about the REST API")

	return
}

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json := jsonbuilder.Object()

	//always answer as last action
	defer func() {
		fmt.Fprintf(w, json.MarshalPretty())
	}()

	switch vars["type"] {
	case "user":

		var user User
		var err error

		if r.FormValue("username") != "" {
			user, err = getUserByUsername(r.FormValue("username"))
		} else if r.FormValue("userid") != "" {
			userID, err := strconv.ParseInt(r.FormValue("userid"), 10, 32)

			if err != nil {
				json.Set("status", "wrong userid format")
				return
			}

			user, err = getUserByUserID(int32(userID))
		} else if r.FormValue("uid") != "" {
			user, err = getUserByUID(r.FormValue("uid"))
		}

		if err != nil {
			json.Set("status", "user error")
			return
		}

		json.Set("status", "ok").Begin("data").Set("userid", user.userid).Set("username", user.username).Set("fullname", user.fullname).End()

		if r.FormValue("token") != "" {
			if verifyToken(user.userid, r.FormValue("token")) {
				json.Set("authLevel", user.authlevel)
			} else {
				json.Set("authLevel", "none")
			}
		}
		return

	default:
		json.Set("status", "unknown query")
		return
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
	json := jsonbuilder.Object()
	defer func() {
		fmt.Fprintf(w, json.MarshalPretty())
	}()

	userid, err := strconv.ParseInt(r.FormValue("userid"), 10, 32)

	result := checkPasswordOnUserID(int32(userid), r.FormValue("password"))

	if result && err == nil {
		json.Set("status", "ok").Begin("data").Set("token", newAuthenticatorSession(int32(userid))).End()
	} else {
		json.Set("status", "wrong password or userid")
	}

}
