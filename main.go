package main

import (
	"net/http"
	"log"
	"os"
	"github.tesla.com/chrzhang/go-microservices-restful/handlers"
	"github.tesla.com/chrzhang/go-microservices-restful/data"
	"github.com/go-openapi/runtime/middleware"
	"time"
	"context"
	"os/signal"
	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	ph := handlers.NewProducts(l, v)

	sm := mux.NewRouter()

	// gives a route that is filtered specifically for HTTP "GET"
	// then using Subrouter method on route, it is converted into a router
	// once it is a router, you can do things like add a handle on it
	//getRouter := sm.Methods("GET").Subrouter()
	// alternatively:
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.Update)
	// Middleware will get executed before actual handler code
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	// can now go to localhost:9090/docs to see docs (JS generated)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))


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