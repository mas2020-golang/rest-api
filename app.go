package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mas2020-golang/rest-api/handlers"
	"github.com/mas2020-golang/rest-api/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// preliminary stuff
var (
	l      = log.New(os.Stdout, "", log.LstdFlags) // logger
	config *pgxpool.Config
)

type App struct {
	Router *mux.Router
	DBPool *pgxpool.Pool
}

func (a *App) Initialize(user, password, host, dbname string) {
	var err error
	// connection to the database
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	config, err = pgxpool.ParseConfig(connectionString)
	if err != nil {
		l.Fatal("check the connection string:", err.Error())
	}

	a.DBPool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		l.Fatal("unable to connect to database: ", err)
	} else {
		l.Println("connection to the database OK!")
	}

	a.Router = mux.NewRouter()
	// create a pwd for the JWT signing algorithm
	utils.Server.GeneratePwd()
	// init the routes
	a.initRoutes()
}

func (a *App) Run(addr string) {
	//var err error
	// http server parameters
	s := &http.Server{
		Addr:         addr,              // configure the bind address
		Handler:      a.Router,          // set the default handler
		IdleTimeout:  120 * time.Second, // max time for connections using the TCP Keep Alive
		ReadTimeout:  1 * time.Second,   // max time to read request from the client
		WriteTimeout: 1 * time.Second,   // max time to write response to the client
	}

	// start the server in a separate go routine
	go func() {
		l.Println("starting http server...")
		s.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 200)
	l.Println("http server is ready to accept connections")
	// wait for a signal to shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// read from the channel
	sig := <-sigChan
	l.Printf("received the %v signal", sig)

	// gracefully shutdown the server (after 10 seconds server is shutdown)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//s.Shutdown(tc)
	a.shutdown(ctx, s)
	l.Println("shutting down, bye")
	os.Exit(0)
}

// Gracefully shutdown function
func (a *App) shutdown(ctx context.Context, s *http.Server) {
	c := make(chan string)
	go func() {
		l.Println("some long running stuff...")
		// db connection: close the pool
		l.Printf("closing db connections...")
		a.DBPool.Close()
		l.Println("closing db connections done!")
		time.Sleep(time.Millisecond * 200)
		c <- "cleanup operations done!"
	}()

	select {
	case done := <-c:
		l.Println("normal shutdown:", done)
		s.Shutdown(ctx)
	case <-ctx.Done():
		l.Println("elapsed timeout:", ctx.Err())
		s.Shutdown(ctx)
	}
}

// initRoutes inits the routes for the application
func (a *App) initRoutes() {
	// new handler object
	ph := handlers.NewProducts(l, a.DBPool)
	//gb := handlers.NewGoodBye(l)
	//a.Router.Use(handlers.AuthMiddleware)

	// products sub router (for every call is checked the Token, for POST and PUT is also used the validation middleware
	prodRouter := a.Router.PathPrefix("/products").Subrouter()
	prodRouter.HandleFunc("", ph.GetProducts).Methods(http.MethodGet)
	prodRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct).Methods(http.MethodGet)
	prodRouter.Use(handlers.AuthMiddleware)

	putPostRouter := prodRouter.Methods(http.MethodPost, http.MethodPut).Subrouter()
	putPostRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct).Methods(http.MethodPut)
	putPostRouter.HandleFunc("", ph.AddProduct).Methods(http.MethodPost)
	putPostRouter.Use(ph.MiddlewareProductValidation)

	// login handler
	a.Router.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)

}
