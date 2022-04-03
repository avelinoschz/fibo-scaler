package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("configuring application")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is happy!"))
	})

	http.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
	})

	log.Println("starting server on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal()
	}
}
