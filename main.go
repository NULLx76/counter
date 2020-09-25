package main

import (
	"counter/store"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Info("Starting up")

	cfg := getConfig()

	s, err := func() (store.Repository, error) {
		switch cfg.DB {
		case dbMemory:
			return store.NewMemoryStore(), nil
		case dbEtcd3:
			return store.NewEtcdStore([]string{"http://localhost:2379"})
		case dbDisk:
			if err := os.Mkdir(cfg.DiskPath, os.ModePerm); err != nil {
				return nil, err
			}
			return store.NewDiskvStore(cfg.DiskPath), nil
		default:
			return nil, fmt.Errorf("unsupported database: %v", cfg.DB)
		}
	}()
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
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
		Addr:    cfg.Address,
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
