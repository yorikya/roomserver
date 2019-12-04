package main

import (
	"net/http"

	"github.com/yorikya/roomserver/httpserver"
	"github.com/yorikya/roomserver/server"
)

func main() {

	s := server.NewServer("room1")
	defer s.Close()

	httpserver.InitRoutes(s)

	http.ListenAndServe(":80", nil)

}
