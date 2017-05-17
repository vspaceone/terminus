package main

import (
	"fmt"
	"log"
	"net/http"

	"sync"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello startin'")
	r := mux.NewRouter()

	http.Handle("/", r)

	r.HandleFunc("/", handler)
	r.HandleFunc("/get/{type}", get)
	//TODO use auth tokens
	r.HandleFunc("/new/{type}", new)
	r.HandleFunc("/auth", auth)

	// Starting servers in goroutines and waiting for their return
	var wg sync.WaitGroup
	wg.Add(2)
	go startHTTP(&wg, r)  // Unsecure server - use should be discouraged
	go startHTTPS(&wg, r) //secure server
	wg.Wait()
}

func startHTTP(wg *sync.WaitGroup, r *mux.Router) {
	defer wg.Done()
	log.Fatal(http.ListenAndServe(":8000", r))
}

func startHTTPS(wg *sync.WaitGroup, r *mux.Router) {
	defer wg.Done()
	// TODO use letsencrypt
	log.Fatal(http.ListenAndServeTLS(":8001", "server.crt", "server.key", r))
}
