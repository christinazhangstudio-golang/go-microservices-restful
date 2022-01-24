package handlers

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
)

// implements HTTPHandler interface
type Hello struct {
	l *log.Logger
}

// Dependency injection where log.Logger is being injected
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// h.l. would be better for log since that kind of dependency injection
	// would be useful/fast for unit testing/connecting to db/etc.
	h.l.Println("Hello world")
	// ResponseWriter and Request are used to read/write
	// e.g. curl -v -d 'Christina' localhost:9090 will give "Data: Christina"
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops", http.StatusBadRequest)
		// Alternatively:
		// WriteHeaders allows specifying a status code to caller
		//rw.WriteHeader(http.StatusBadRequest)
		//rw.Write([]byte("Oops"))
		return
	}
	h.l.Printf("Data: %s\n", d)

	// to write back for client, we used ResponseWriter
	fmt.Fprintf(rw, "Hello %s\n", d)
}