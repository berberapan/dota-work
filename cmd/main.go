package main

import (
	"log"
	"net/http"

	"github.com/berberapan/dota-work/internal/server"
)

func main() {
	router := http.NewServeMux()
	newServer := server.NewServer(":80", ":443", router, nil)
	if err := newServer.Run(); err != nil {
		log.Fatal(err)
	}
}
