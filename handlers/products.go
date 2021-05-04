package handlers

import (
	"Microservices/Introduction/data"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")

	// prod := &data.Product{}

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
	// }

	// p.l.Printf("Prod: %#v", prod)
	// data.AddProduct(prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT Products", id)

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

//
// MiddlewareVlaidateProduct validates the product in the request and calls next if ok
func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		p.l.Println("Request Body: ", &r.Body)
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to unmarshall json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

// func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.getProducts(rw, r)
// 		return
// 	}
// 	if r.Method == http.MethodPost {
// 		p.addProduct(rw, r)
// 		return
// 	}
// 	// handle an update
// 	if r.Method == http.MethodPut {
// 		//expect the id in the URI
// 		reg := regexp.MustCompile(`/([0-9]+)`)
// 		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

// 		if len(g) != 1 {
// 			p.l.Println("Invalid URI more than one id")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		if len(g[0]) != 2 {
// 			p.l.Println("Invalid URI more than one capture group")
// 			http.Error(rw, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}
// 		idString := g[0][1]
// 		id, err := strconv.Atoi(idString)
// 		if err != nil {
// 			p.l.Println("Invalid URI unable to convert to number")
// 			http.Error(rw, "Invalid URL", http.StatusBadRequest)
// 			return
// 		}

// 		p.updateProducts(id, rw, r)
// 		p.l.Println("got id", id)
// 		return

// 	}
// 	//catch all
// 	rw.WriteHeader(http.StatusMethodNotAllowed)
// }
