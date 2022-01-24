package main

import (
	"net/http"
	"log"
	"os"
	"github.tesla.com/chrzhang/go-microservices-restful/handlers"
	"time"
	"context"
	"os/signal"
)

func main() {
	// HandleFunc registers a function to a path on default serve mux (a default http handler)
	//http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		//this code has been moved to package hello, since that kind of compartmentalization is good
	//})

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	//hh := handlers.NewHello(l)
	//gh := handlers.NewGoodbye(l)
	ph := handlers.NewProducts(l)
	//need to register this handler with server
	//first need to create serve mux, which has methods like Handle and Handler
	//when request comes into server, the server has default handler (which is HTTP serve mux),
	//server will call http.ServeMux.serveHTTP, that logic will determine which handler registered to it to call
	sm := http.NewServeMux()
	//sm.Handle("/", hh)
	//sm.Handle("/goodbye", gh)
	sm.Handle("/", ph)

	//manually creating a server
	s := &http.Server{
		Addr: ":9090",
		Handler: sm, 
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// the ListenAndServe will not block since it is now wrapped up in a go function
	go func() {
		// creates the web server - to handle requests, http handlers have to be implemented 
		// first parameter is bind address and second is http handler
		// registers a default handler and uses default serve mux if we specify nil for second parameter
		// a handler in go http is just an interface which has the ServeHTTP(ResponseWriter, *Request) method
		//http.ListenAndServe(":9090", sm)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	// broadcast a notification to channel whenever kill command/interrupt is received
	signal.Notify(sigChan, os.Kill)

	// reading from a channel will block until message is consumed
	// once the message is consumed, server is shutdown
	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	// graceful shutdown waits for requests to be done until shutting down the server
	// useful for finishing up database transactions
	// first create a context with duration of 30 seconds (give graceful shutdown 30 sec until forcefully shutting down)
	tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(tc)
}