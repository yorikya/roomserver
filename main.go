package main

import (
	"net/http"

	"github.com/yorikya/roomserver/httpserver"
	"github.com/yorikya/roomserver/mqttserver"
)

func main() {

	s := mqttserver.NewServer("office")
	go s.Start()
	defer s.Close()

	httpserver.InitRoutes(s)

	http.ListenAndServe(":3000", nil)

}
