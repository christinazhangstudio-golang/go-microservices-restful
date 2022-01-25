# go-microservices-restful
The Swagger API documentation that was created is used to generate a Go client SDK (HTTP client code). 

Generating the client code is the following command (in the sdk folder):
`swagger generate client -f ../../product-api/swagger.yaml -A product-api`