package main

import "github.com/heroku/go-getting-started/internal/server"

func main() {
	//port := os.Getenv("PORT")
	//
	//if port == "" {
	//	log.Fatal("$PORT must be set")
	//}
	//
	//server.Start(":" + port)
	server.Start(":" + "8000")
}
