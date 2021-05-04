package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle Hello requests")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// rw.WriteHeader(http.StatusBadRequest)
		// rw.Write([]byte("Ooops"))
		h.l.Println("Error reading body", err)
		http.Error(rw, "Unable to read request body", http.StatusBadRequest)

		return
	}
	log.Printf("Data %s\n", d)

	fmt.Fprintf(rw, "Hello %s", d)
}
