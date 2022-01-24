package handlers

import (
	"log"
	"github.tesla.com/chrzhang/go-microservices-restful/data"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//handle get
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	//handle add
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// handle update
	if r.Method == http.MethodPut {
		//expect the id in the URI (e.g. localhost:9090/1) and parse it out
		regx := regexp.MustCompile(`/([0-9]+)`)
		// will return groups as string array
		g := regx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		//p.l.Println("got id", id)

		p.updateProduct(id, rw, r)
		return
	}

	//catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	// create a new product object
	prod := &data.Product{}
	// reader is response body from HTTP Request
	// go actually hasn't read everything from the client at the point of
	// receiving an HTTP request - will buffer some stuff but not everything
	// so reading progressively (instead of storing it in one data slice) is
	// better with a buffered writer (response writer body)
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	// # allows to view fields
	p.l.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT products")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err != data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}