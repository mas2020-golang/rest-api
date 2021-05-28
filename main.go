package main

import (
	"context"
	"github.com/mas2020-golang/rest-api/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// preliminary stuff
	l := log.New(os.Stdout, "", log.LstdFlags) // logger
	// new handler object
	products := handlers.NewProducts(l)
	gb := handlers.NewGoodBye(l)

	sm := http.NewServeMux()
	sm.Handle("/", products)
	sm.Handle("/goodbye", gb)

	// http server parameters
	s := &http.Server{
		Addr:         ":9090",           // configure the bind address
		Handler:      sm,                // set the default handler
		IdleTimeout:  120 * time.Second, // max time for connections using the TCP Keep Alive
		ReadTimeout:  1 * time.Second,   // max time to read request from the client
		WriteTimeout: 1 * time.Second,   // max time to writeß response to the client
	}

	// start the server in a separate go routine
	go func() {
		l.Println("starting http server...")
		s.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 200)
	l.Println("http server is ready")
	// wait for a signal to shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// read from the channel
	sig := <-sigChan
	l.Printf("received the %v signal", sig)

	// gracefully shutdown the server
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
