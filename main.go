package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
func run() error {
	log.Print("configuring application")

	v1, v2 := 0, 1
	http.HandleFunc("/previous", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/previous", v1)
		intByte := []byte(strconv.Itoa(v1))
		w.Write(intByte)
	})

	http.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/current", v2)
		intByte := []byte(strconv.Itoa(v2))
		w.Write(intByte)
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		f := v1
		v1, v2 = v2, v2+f
		log.Println("/next", v2)
		intByte := []byte(strconv.Itoa(v2))
		w.Write(intByte)
	})

	errChan := make(chan error, 1)
	go func(c chan error) {
		log.Println("starting server on port :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal()
		}
	}(errChan)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	select {
	case err := <-errChan:
		return err

	case <-signalChan:
		log.Println("shutting down gracefully")
		return nil
	}
}
