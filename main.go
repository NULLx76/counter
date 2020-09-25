package main

import (
	"counter/store"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	log.Info("Starting up")
	s, err := store.NewEtcdStore([]string{"http://localhost:2379"})
	if err != nil {
		log.Fatalf("Couldn't connect to etcd")
	}

	defer s.Close()
	rs := NewRoutes(s)

	// Router
	r := mux.NewRouter()
	r.PathPrefix("/").Methods("GET").HandlerFunc(rs.GetCounter)
	r.PathPrefix("/").Methods("PUT").HandlerFunc(rs.IncrementCounter)
	r.PathPrefix("/").Methods("POST").HandlerFunc(rs.CreateCounter)
	r.PathPrefix("/").Methods("DELETE").HandlerFunc(rs.DeleteCounter)
	r.Use(rootMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Started listing on %v", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func rootMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/" {
			_, _ = fmt.Fprint(w, "Hello World!")
		} else {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}
	})
}
