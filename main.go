package main

import (
	"log"
	"os"

	"github.com/heroku/go-getting-started/server"
	"github.com/heroku/go-getting-started/simulator"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	go simulator.Start()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server.Start(":" + port)
	// server.Start(":" + "8000")
}
