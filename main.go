package main

import (
	"net/http"
	"log"

	"github.com/yorikya/roomserver/httphandlers"
	"github.com/yorikya/roomserver/server"
)

func main() {

	s := server.NewServer("room1_main", "room1_dht")
	defer s.Close()

	httphandlers.InitRoutes(s)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
