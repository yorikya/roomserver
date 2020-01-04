package main

import (
	"log"
	"net/http"

	"github.com/yorikya/roomserver/httphandlers"
	"github.com/yorikya/roomserver/server"
)

func main() {

	s := server.NewServer("config/rooms.json")
	defer s.Close()

	httphandlers.InitRoutes(s)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
