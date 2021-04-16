package Handlers

import (
	"net/http"
	"fmt"
	"log"
	"io/ioutil"
)

type GoodBye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

func (g *GoodBye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Helo")
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Data %s\n", d)
}

