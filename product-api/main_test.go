package main

import (
	"testing"
	"github.tesla.com/chrzhang/go-microservices-restful/product-api/sdk/client"
	"github.tesla.com/chrzhang/go-microservices-restful/product-api/sdk/client/products"
	"fmt"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod)
}