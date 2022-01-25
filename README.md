# go-microservices-restful
Refer to branch without-gorilla-refactor for all comments related to using default server setup.

Gorilla is a web framework that implements a request router and dispatcher for matching incoming requests to their respective handler. We can register more detailed handlers. A router gives us the capability of subrouters, and middleware can be more flexibly added.

The default server mux (in main.go) is replaced with the gorilla/mux router.

On a router, you can add Methods, which registers a new route with a matcher for HTTP methods - you can create a route that is a type "GET" for example and register handlers on that

The Use function allows us to implement Middleware. Middleware is just an HTTP handler, and we can used Middleware pattern to chain multiple handlers together.

Middleware is often for validating requests or authentication. Validate was added to our middleware function.ß
