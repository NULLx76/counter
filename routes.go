package main

import (
	"counter/store"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

// TODO: Tests

func marshal(key string, value store.Value) string {
	return fmt.Sprintf("{ \"%v\": \"%v\" }", key, value.Count)
}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	c := s.Get(r.RequestURI)

	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return false
	}
	token := strings.TrimPrefix(header, "Bearer ")
	parsedToken, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return false
	}

	if c.AccessKey != parsedToken {
		http.Error(w, "Wrong access token", http.StatusUnauthorized)
		return false
	} else {
		return true
	}
}

func getCounter(w http.ResponseWriter, r *http.Request) {
	c := s.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	_, err := fmt.Fprint(w, marshal(r.RequestURI, c))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	c := s.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	if !authenticate(w, r) {
		return
	}

	s.Increment(r.RequestURI)
	getCounter(w, r)
}

func createCounter(w http.ResponseWriter, r *http.Request) {
	c := s.Get(r.RequestURI)
	if c.AccessKey != uuid.Nil {
		http.Error(w, "Counter already exists", http.StatusConflict)
		return
	}

	v := store.Value{Count: 0, AccessKey: uuid.New()}
	s.Set(r.RequestURI, v)

	w.Header().Set("Authorization", "Bearer " + v.AccessKey.String())
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Internal server error encountered when marshalling Value to json", http.StatusInternalServerError)
	}
}

func deleteCounter(w http.ResponseWriter, r *http.Request) {
	c := s.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	if !authenticate(w, r) {
		return
	}

	s.Set(r.RequestURI, store.Value{})
}
