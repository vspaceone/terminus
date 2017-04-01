package main

import (
	"fmt"
	"net/http"

	"github.com/coyove/jsonbuilder"
	"github.com/gorilla/mux"
)

//handler's response informs briefly about this service
func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Terminus Server alpha1\nSee docs or ask the devs about the REST API")

	return
}

//userRequestByUID's response shows the user entry in the database
func userRequestByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := getUserByUID(vars["uid"])
	json := jsonbuilder.Object()

	//For some reason jsonbuilder.From(user) does not work
	json.Set("firstName", string(user.firstName)).Set("lastName", string(user.lastName))

	fmt.Fprintf(w, json.MarshalPretty())

	return
}

// userRequestPermissionToken requests a permission token from the server,
// that autheticates the user for an certain amount of time allowing him not to log in for a while
func userRequestPermissionToken(w http.ResponseWriter, r *http.Request) {
	//verify user password and generate and return a token
	//store token in Authenticator (TODO)
}
