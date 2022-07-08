package main

import (
	"log"
	"os"

	"github.com/heroku/go-getting-started/server"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server.Start(":" + port)
	//server.Start(":" + "8000")
}
