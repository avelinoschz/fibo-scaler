package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
)

// main abstracts the behavior of the application
// is it running correctly or does it got an error.
// In case of any error the app exits with stdout 1,
// this error can be picked up by the container
// there a restart policy can be managed depending
// on the cases. Max retries on-failure, unless-stopped
func main() {
	if err := run(); err != nil {
		log.Printf("error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// run wraps the execution of the application.
// here is the initialization of the http server
// and the listener for a interrupt signal to
// enable gracefully shutdown
func run() error {
	fh := &fiboHandler{
		prev:    0,
		current: 1,
	}

	errChan := make(chan error, 1)
	go func(c chan error) {
		log.Println("listening on port :8080")
		if err := http.ListenAndServe(":8080", fh); err != nil {
			log.Println("err:", err)
			c <- err
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

// paths used by the handler
var (
	urlRoot    = regexp.MustCompile(`^\/[\/]*$`)
	urlPrev    = regexp.MustCompile(`^\/previous[\/]*$`)
	urlCurrent = regexp.MustCompile(`^\/current[\/]*$`)
	urlNext    = regexp.MustCompile(`^\/next[\/]*$`)
	urlError   = regexp.MustCompile(`^\/error[\/]*$`)
)

// fiboHandler is the handler embedded into de http server.
// Since it is a ver simply application and aiming to guarantee
// a good performance, it was decided to use only the go
// stdlib, instead of an http framework or router.
// fiboHandler is needed to have a shared current state of the
// fibonacci sequence that can be accesed by all the requests.
type fiboHandler struct {
	prev    int
	current int
}

// ServeHTTP is the implementation needed for fiboHandler to satisfy
// the `http.Handler` interface. Multiplexes the received requests
// to it's corresponding handler to process.
// This could be achieved with simple `HandleFunc` in the main()
// and benefit from closures to have access to a shared state by all
// the HandleFunc, but having in mind testability and unit testing,
// it was decided to use separated handlers.
func (fh *fiboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && urlRoot.MatchString(r.URL.Path):
		handlerRoot(w, r)
	case r.Method == http.MethodGet && urlPrev.MatchString(r.URL.Path):
		fh.handlerPrevious(w, r)
	case r.Method == http.MethodGet && urlCurrent.MatchString(r.URL.Path):
		fh.handlerCurrent(w, r)
	case r.Method == http.MethodGet && urlNext.MatchString(r.URL.Path):
		fh.handlerNext(w, r)
	case r.Method == http.MethodGet && urlError.MatchString(r.URL.Path):
		handlerError(w, r)
	default:
		handlerNotFound(w, r)
	}
}

// handler intented to be used as a health check that the app is running
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("/")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is happy!"))
}

// handler that simply return the previous fibonacci number at a given moment
func (fh *fiboHandler) handlerPrevious(w http.ResponseWriter, r *http.Request) {
	log.Println("/previous", fh.prev)
	intByte := []byte(strconv.Itoa(fh.prev))
	w.WriteHeader(http.StatusOK)
	w.Write(intByte)
}

// handler that simply return the current fibonacci number at a given moment
func (fh *fiboHandler) handlerCurrent(w http.ResponseWriter, r *http.Request) {
	log.Println("/current", fh.current)
	intByte := []byte(strconv.Itoa(fh.current))
	w.Write(intByte)
}

// handler that does the actual math of the fibonacci sequence
func (fh *fiboHandler) handlerNext(w http.ResponseWriter, r *http.Request) {
	prev := fh.prev
	fh.prev = fh.current
	fh.current = prev + fh.current
	log.Println("/next", fh.current)
	intByte := []byte(strconv.Itoa(fh.current))
	w.WriteHeader(http.StatusOK)
	w.Write(intByte)
}

// default handler for all the unknown incoming requests
func handlerNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not-found"))
}

// every received handler request is assigned to a new goroutine.
// what this means is that if a single request panics, only the scope
// of the handler crashes, and the http server will continue to handle
// incoming requests
func handlerError(w http.ResponseWriter, r *http.Request) {
	panic("oh no! a panic")
}
