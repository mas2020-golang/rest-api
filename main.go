package main

import (
	"github.com/mas2020-golang/rest-api/server"
	"os"
)

func main() {
	// start the application
	a := server.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":9090")
}
