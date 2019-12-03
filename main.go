package main

import (
	"net/http"

	"github.com/yorikya/roomserver/httpserver"
	"github.com/yorikya/roomserver/mqttserver"
)

func main() {

	s := mqttserver.NewServer("room1")
	// go s.Start()
	defer s.Close()

	httpserver.InitRoutes(s)

	http.ListenAndServe(":80", nil)

}
