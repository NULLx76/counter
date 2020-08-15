package main

import (
	"counter/store"
	"fmt"
	"log"
	"net/http"
)

var s store.Repository

func main() {
	s = store.NewMemoryStore()
	http.HandleFunc("/", handleHttp)

	fmt.Print("Starting server...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		_, _ = fmt.Fprint(w, "Yeet!")
		return
	}

	// TODO: Authentication
	// 		 POST to create a counter and receive a unique token to update that counter?

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
		case "GET":
			_, err := fmt.Fprint(w, marshal(r.RequestURI, s.Get(r.RequestURI)))
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
			}
		case "PUT":
			s.Increment(r.RequestURI)
			_, err := fmt.Fprint(w, marshal(r.RequestURI, s.Get(r.RequestURI)))
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
			}
		case "POST":
			fallthrough
		default:
			http.Error(w, "Unsupported method", 400)
			return
	}
}

func marshal(key string, value int) string {
	return fmt.Sprintf("{ \"%v\": \"%v\" }", key, value)
}
