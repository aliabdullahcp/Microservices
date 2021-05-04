package main

import (
	"Microservices/Introduction/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ph := handlers.NewProducts(l)
	// gh := handlers.NewGoodbye(l)

	// sm := http.NewServeMux()
	// sm.Handle("/", ph)
	// sm.Handle("/goodbye", gh)

	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareValidateProduct)

	s := &http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sm,                //set the default handler
		IdleTimeout:  120 * time.Second, //set the loggeer for the server
		ReadTimeout:  1 * time.Second,   //max time to read request from the client
		WriteTimeout: 1 * time.Second,   //max time to write request to the client
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//Gracefully Shutdown like doesn't closes server until all the clients have finished their work
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Recieved terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}

// //TUTORIAL # 1
// handlers.TestFunc()
// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
// 	log.Println("Hello World")
// 	d, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		// rw.WriteHeader(http.StatusBadRequest)
// 		// rw.Write([]byte("Ooops"))
// 		http.Error(rw, "Oops", http.StatusBadRequest)

// 		return
// 	}
// 	log.Printf("Data %s\n", d)

// 	fmt.Fprintf(rw, "Hello %s", d)
// })

// http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
// 	log.Println("Goodbye World")
// })

// //hh := handlers.NewHello()

// http.ListenAndServe(":9090", nil)
