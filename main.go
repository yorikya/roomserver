package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yorikya/roomserver/httphandlers"
	"github.com/yorikya/roomserver/server"
)

func main() {

	s := server.NewServer("config/rooms.json")
	defer s.Close()

	//Catch Ctrl + c interupt
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		s.Close()
		os.Exit(0)
	}()

	httphandlers.InitRoutes(s)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
