package main

import (
	"flag"
	"log"

	"github.com/mhghw/fara-message/api"
	"github.com/mhghw/fara-message/db"
)

// implement this with os args
var port = flag.Int("port", 8080, "Port to run the HTTP server")

func main() {
	db.InitialDatabase()
	flag.Parse()
	err := api.RunWebServer(*port)
	if err != nil {
		log.Print("failed to start HTTP server:", err)
	}

}
