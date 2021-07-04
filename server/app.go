package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mas2020-golang/goutils/fs"
	"github.com/mas2020-golang/goutils/output"
	"github.com/mas2020-golang/rest-api/handlers"
	"github.com/mas2020-golang/rest-api/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// preliminary stuff
var (
	//l      = log.New(os.Stdout, "", log.LstdFlags) // logger
	config *pgxpool.Config
)

type App struct {
	Router *mux.Router
	DBPool *pgxpool.Pool
}

func (a *App) Initialize(user, password, host, dbname string) {
	var err error
	// load config file
	loadConfig()

	// log settings
	logrus.SetLevel(logrus.Level(utils.Server.Logging.Level))
	logrus.SetFormatter(&output.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	// connection to the database
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	config, err = pgxpool.ParseConfig(connectionString)
	output.CheckErrorAndExitLog("","check the connection string: ", err)

	a.DBPool, err = pgxpool.ConnectConfig(context.Background(), config)
	output.CheckErrorAndExitLog("","unable to connect to database: ", err)
	output.InfoLog("", "connection to the database OK!")

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
		output.InfoLog("", "starting http server...")
		s.ListenAndServe()
		output.InfoLog("", "closing http server...")
	}()

	time.Sleep(time.Millisecond * 100)
	output.InfoLog("", "http server is ready to accept connections")
	// wait for a signal to shutdown the server
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	// read from the channel
	sig := <-sigChan
	output.InfoLog("", fmt.Sprintf("received the %v signal", sig))

	// gracefully shutdown the server (after 10 seconds server is shutdown)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	a.shutdown(ctx, s)
	output.InfoLog("", "shutting down bye!")
	os.Exit(0)
}

// Gracefully shutdown function
func (a *App) shutdown(ctx context.Context, s *http.Server) {
	c := make(chan string)
	go func() {
		output.TraceLog("", "some long running stuff...")
		// db connection: close the pool
		output.DebugLog("", "closing db connections...")
		a.DBPool.Close()
		output.DebugLog("", "closing db connections done!")
		time.Sleep(time.Millisecond * 200)
		c <- "cleanup operations done!"
	}()

	select {
	case done := <-c:
		output.InfoLog("", fmt.Sprintf("normal shutdown: %s", done))
		s.Shutdown(ctx)
	case <-ctx.Done():
		output.InfoLog("", "elapsed timeout")
		s.Shutdown(ctx)
	}
}

// initRoutes inits the routes for the application
func (a *App) initRoutes() {
	// new handler object
	ph := handlers.NewProducts(a.DBPool)
	// common middleware valid for all the calls
	a.Router.Use(commonMiddleware)

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

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func loadConfig(){
	// does the env variable exist?
	if len(os.Getenv("APP_CONFIG")) > 0{
		ok, err := fs.ExistsPath(os.Getenv("APP_CONFIG"))
		output.CheckErrorAndExit("", "", err)
		if ok{
			// load the config file
			err = fs.ReadYaml(os.Getenv("APP_CONFIG"), &utils.Server)
			output.CheckErrorAndExit("", "", err)
			return
		}else{
			output.ErrorLog("", "an error occurred during the load of the configuration file: " +
				"APP_CONFIG points to a wrong path")
			os.Exit(1)
		}
	}

	// read from the default location
	ok, err := fs.ExistsPath("config/server.yml")
	output.CheckErrorAndExit("", "", err)
	if ok{
		// load the config file
		err = fs.ReadYaml("config/server.yml", &utils.Server)
		output.CheckErrorAndExit("", "", err)
		return
	}else{
		output.ErrorLog("", "an error occurred during the load of the configuration file: config/server.yml doesn't exist")
		os.Exit(1)
	}

}