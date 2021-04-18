package main

import (
	"Microservices/NotMod/Handlers"
	//"context"
	"os/signal"
	"time"

	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	// hh := Handlers.NewHello(l)
	// gh := Handlers.NewGoodbye(l)
	ph := Handlers.NewProducts(l)


	sm := mux.NewRouter()

	getRouter := sm.Methods("GET").Subrouter()
	//sm.Handle("/products", ph)
	//sm.Handle("/goodbye", gh)
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/{id :[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	ops := middleware.RedocOpts{SpecURL : "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr: ":9090",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	go func () {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)


	tc , _:= context.WithTimeout(context.Background(), 30 * time.Second)
	s.Shutdown(tc)
	
}