# go-microservices-restful

Multi-part requests is somewhat deprecated. It's technically valid, but most people are not using plain HTTP to POST data. It's not RESTful either. For that reason, it's not that relevant to include multi-part support here.

Moving on to gRPC services! JSON and REST are not maximally optimized. gRPC will use binary-based message protocol using protobufs. In a proto file, services, methods for services, and messages (input/output) for methods are defined.

https://grpc.io/docs/languages/go/basics/#client-side-streaming-rpc

An interface (i.e. `CurrencyServer`) is generated to create a gRPC server using the method `GetRate()`. A struct needs to implement this interface. To create the server, we need the `RegisterCurrencyServer`, which maps the implementation of the interface (how you want to handle getRate()) to gRPC server. In JSON terms, this is kinda like `CurrencyServer` is the handlers, and gRPC is HTTP server - matching routes to a server. The struct `Currency` that implements `CurrencyServer` is defined in `currency.go` under server.

`grpccurl` can be used for testing.

```
$ grpcurl --plaintext localhost:9092 list
Currency
grpc.reflection.v1alpha.ServerReflection
```

```
$ grpcurl --plaintext localhost:9092 describe Currency.GetRate
Currency.GetRate is a method:
rpc GetRate ( .RateRequest ) returns ( .RateResponse );
```

```
$ grpcurl --plaintext localhost:9092 describe .RateRequest
RateRequest is a message:
message RateRequest {
  string Base = 1;
  string Destination = 2;
}
```

```
$ grpcurl --plaintext -d '{"Base":"GBP", "Destination":"USD"}' localhost:9092 Currency.GetRate
{
  "rate": 0.5
}
```


Now that a gRPC service is created (i.e. `CurrencyService`), need to look at how CurrencyService can be integrated with API service or how there can be client-side calls from go into CurrencyService.

First is adding enumerations for allowed currencies. Also `Base/Destination` are changed from string type to this enumeration type.

We want to be able to call the service (the one that defines `Currency`) from `product-api`; product API has upstream dependencies, such as conversion of currency. So we need to construct a client that calls the `Currency` service.

To do this, we just use the client creation method generated from the proto file - `CurrencyClient`. We create a `currencyClient` using the `NewCurrencyClient` method. Protobufs are good since clients are interfaces and testing would be easy!

In `main.go` in `product-api`, we create an instance of this client.

Recap: Call `getRate()` method from `get.go`, which is from the `CurrencyClient` which has been auto-gen'ed from proto file. We constructed a new `CurrencyClient` using the address of the service and created a gRPC connection, which we then use in `get.go` (`ListSingle`)