package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GoodBye struct {
	// internal logger
	l *log.Logger
}

// func to build a new GoodBye struct
func NewGoodBye(l *log.Logger) *GoodBye{
	return &GoodBye{l}
}

func (h *GoodBye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("GoodBye handler called")
	if d, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, "an error occurred", http.StatusInternalServerError)
	} else {
		// return body to the caller
		fmt.Fprintf(w, "Goodbye %s", d)
	}
}
