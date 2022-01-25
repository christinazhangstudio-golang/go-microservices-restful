package handlers

import (
	"log"
	"github.tesla.com/chrzhang/go-microservices-restful/data"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"context"
	"fmt"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// with using gorilla/mux this can be converted to a public method
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	// converting product struct to JSON representation using encoding/json Marshal
	// GetProducts is defined in products.go in data package
	lp := data.GetProducts()
	// we can use this below implementation, but an encoder would be even better
	// because writing direct does not have to allocate memory,
	// and the encoder is marginally faster than marshal (especially when multiple threads are involved)
	// so the ToJSON method in products.go in the data package is used
	//d, err := json.Marshal(lp)
	//if err != nil {
	//	http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	//}

	//rw.Write(d)
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	// MOVED THIS TO MIDDLEWARE
	// prod := &data.Product{}

	// err = prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	// }

	prod := r.Context().Value(KeyProduct{}).(data.Product)	//to cast

	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	//when variables are passed in to the URL, gorilla uses mux.Vars to extract them
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT products")

	// MOVED THIS TO MIDDLEWARE
	// prod := &data.Product{}

	// err = prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	// }
	prod := r.Context().Value(KeyProduct{}).(data.Product)	//to cast

	err = data.UpdateProduct(id, &prod)
	if err != data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}

// to use context in go, we define a key
type KeyProduct struct {}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("Error validating product", err)
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("Error validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}


		// put the product onto the request, since it has context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod) 	
		// make a copy of that context and call next
		req := r.WithContext(ctx)	
		next.ServeHTTP(rw, req)
	})
}