package main

import (
	"github.com/mas2020-golang/rest-api/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// preliminary stuff
	l := log.New(os.Stdout, "rest-api ", log.LstdFlags) // logger
	// new handler object
	hh := handlers.NewHello(l)
	gb := handlers.NewGoodBye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gb)

	http.ListenAndServe(":9090", sm)
}
