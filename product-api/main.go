package main

import (
	"net/http"
	"log"
	"os"
	"github.com/go-openapi/runtime/middleware"
	"time"
	"context"
	"os/signal"
	"github.com/gorilla/mux"

	"github.tesla.com/chrzhang/go-microservices-restful/product-api/handlers"
	"github.tesla.com/chrzhang/go-microservices-restful/product-api/data"

	gohandlers "github.com/gorilla/handlers"

	protos "github.tesla.com/chrzhang/go-microservices-restful/currency/protos"
	"google.golang.org/grpc"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()


	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// create client, we need to connect to particular service, so the gRPC dial is defined above with an address
	cc := protos.NewCurrencyClient(conn)

	// once we have currency client, we want to pass it to handlers
	ph := handlers.NewProducts(l, v, cc)

	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	// can now go to localhost:9090/docs to see docs (JS generated)
	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//CORS - to allow bypass of CORS block for ReactJS app
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))


	//manually creating a server
	s := &http.Server{
		Addr: ":9090",
		Handler: ch(sm), //wrap standard serve mux with CORS handler
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