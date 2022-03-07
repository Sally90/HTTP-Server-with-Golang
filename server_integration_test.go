package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

//integration test tests everything at once: server, store, POST and GET requests

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, NewGetRequest(player))
	assertCode(t, response.Code, http.StatusOK)
	assertScore(t, response.Body.String(), "3")
}

func NewPostRequest(player string) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf("/players/%s", player), nil)
	return req
}

func NewGetRequest(player string) *http.Request {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/players/%s", player), nil)
	return req
}
