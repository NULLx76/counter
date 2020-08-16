package main

import (
	"counter/store"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Routes struct {
	repo store.Repository
}

func NewRoutes(repo store.Repository) Routes {
	return Routes{repo}
}

func marshal(key string, value store.Value) string {
	return fmt.Sprintf("{ \"%v\": %v }", key, value.Count)
}

func (rs *Routes) authenticate(w http.ResponseWriter, r *http.Request) bool {
	c := rs.repo.Get(r.RequestURI)

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

func (rs *Routes) GetCounter(w http.ResponseWriter, r *http.Request) {
	c := rs.repo.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	_, err := fmt.Fprint(w, marshal(r.RequestURI, c))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (rs *Routes) IncrementCounter(w http.ResponseWriter, r *http.Request) {
	c := rs.repo.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	if !rs.authenticate(w, r) {
		return
	}

	rs.repo.Increment(r.RequestURI)
	rs.GetCounter(w, r)
}

func (rs *Routes) CreateCounter(w http.ResponseWriter, r *http.Request) {
	c := rs.repo.Get(r.RequestURI)
	if c.AccessKey != uuid.Nil {
		http.Error(w, "Counter already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)

	v := store.Value{Count: 0, AccessKey: uuid.New()}
	rs.repo.Set(r.RequestURI, v)

	w.Header().Set("Authorization", "Bearer " + v.AccessKey.String())
	_, err := fmt.Fprintf(w, "{ \"%v\": %v, \"AccessKey\": \"%v\" }", r.RequestURI, v.Count, v.AccessKey.String())
	if err != nil {
		http.Error(w, "Internal server error encountered when formatting response", http.StatusInternalServerError)
	}
}

func (rs *Routes) DeleteCounter(w http.ResponseWriter, r *http.Request) {
	c := rs.repo.Get(r.RequestURI)
	if c.AccessKey == uuid.Nil {
		http.Error(w, "Counter not yet created", http.StatusNotFound)
		return
	}

	if !rs.authenticate(w, r) {
		return
	}

	rs.repo.Set(r.RequestURI, store.Value{})
}
