package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

//we have server with link to store
//PlayerServer has ServeHTTP method, thus implements the Handler interface
type PlayerServer struct {
	store PlayerStore
}

//we have a store interface
//so the methods of the store interface will be received by a concrete implementation of the store interface / by a conrete store
//this store interface enables you to have different store implementations (one in the main code and a mock in the test code)
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(player string)
}

//concrete type of store which implements the store interface
type InMemoryStore struct {
	store map[string]int
	mu    sync.Mutex
}

//function to intialize the concrete store and to use this initialized store in the integration test
func NewInMemoryPlayerStore() *InMemoryStore {
	return &InMemoryStore{store: map[string]int{}} //initializes the map
}

//concrete type of store implementing the store interface
func (s *InMemoryStore) GetPlayerScore(name string) int {
	player := s.store[name]
	return player
}

//concrete type of store implementing the store interface
func (s *InMemoryStore) RecordWin(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[name]++
}

//server needs to implement the ServeHTTP method to implement the Handler interface
//so that the server can be used in the ListenAndServe method
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	if r.Method == "GET" {
		p.showScore(player, w)
	}
	if r.Method == "POST" {
		p.processWin(player, w)
	}
}

//function to deal with GET requests
func (p *PlayerServer) showScore(player string, w http.ResponseWriter) {

	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprint(w, score)
}

//function to deal with POST requests
func (p *PlayerServer) processWin(player string, w http.ResponseWriter) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
