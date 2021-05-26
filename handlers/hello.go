package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	// internal logger
	l *log.Logger
}

// func to build a new Hello struct
func NewHello(l *log.Logger) *Hello{
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("hello handler called")
	if d, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
	} else {
		// return body to the caller
		fmt.Fprintf(w, "Hello %s", d)
	}
}
