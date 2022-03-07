package main

import (
	"log"
	"net/http"
)

func main() {
	//server is assigned to a concrete store
	server := &PlayerServer{NewInMemoryPlayerStore()}
	//ListenAndServe takes in a Handler as second parameter meaning that the second parameter has to implement the Handler interface
	//To implement the Handler interface, you have to have the ServeHTTP method
	//server is of type *PlayerServer which has the ServeHTTP method
	log.Fatal(http.ListenAndServe(":5000", server))
}
