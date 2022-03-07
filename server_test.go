package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//mock store, so that you can define the store according to the needs of your test
//consider that the store is defined in the TestGETPlayer and the TestPOSTWin tests respectively
//you need a store interface (!) in the main code to have different implementations of it and to use a mock in the test file

type mockStore struct {
	playerScores map[string]int
	winCalls     []string
}

func (s *mockStore) GetPlayerScore(name string) int {
	score := s.playerScores[name]
	return score
}

func (s *mockStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayer(t *testing.T) {
	mockStore := mockStore{
		playerScores: map[string]int{
			"pepper": 20,
			"floyd":  10,
		},
		winCalls: nil,
	}
	server := &PlayerServer{
		store: &mockStore,
	}
	t.Run("returns Peppers score", func(t *testing.T) {
		//nil because we don't need to set a response body in this GET request:
		req, _ := http.NewRequest("GET", "/players/pepper", nil)
		//spy for the response:
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		got := res.Body.String()
		want := "20"
		assertScore(t, got, want)
	})
	t.Run("returns Floyd score", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/players/floyd", nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assertScore(t, res.Body.String(), "10")
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/players/borbe", nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusNotFound)
	})
	t.Run("returns 200 on existing player pepper", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/players/pepper", nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusOK)
	})
}

func TestPOSTWins(t *testing.T) {
	mockStore := mockStore{
		playerScores: map[string]int{},
		winCalls:     nil,
	}
	server := &PlayerServer{store: &mockStore}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/players/pepper", nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusAccepted)
	})
	t.Run("it records wins when POST", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/players/pepper", nil)
		res := httptest.NewRecorder()
		server.ServeHTTP(res, req)
		assertCode(t, res.Code, http.StatusAccepted)
		if mockStore.winCalls[0] != "pepper" {
			t.Errorf("did not store correct winner. got: %s, want: %s", mockStore.winCalls[0], "pepper")

		}
	})
}

func assertScore(t testing.TB, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("The response body is wrong. Got: %s, want: %s", got, want)
	}
}

func assertCode(t testing.TB, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("The response code is wrong. Got: %d, want: %d", got, want)
	}
}
