package main

//go:generate sed -i "s/WIFI_ACCESS_POINT_NAME/$WIFI_ACCESS_POINT_NAME/g" arduino/room1/esp8622main.ino & sed -i "s/WIFI_ACCESS_POINT_PASSWORD/$WIFI_ACCESS_POINT_PASSWORD/g" arduino/room1/esp8622main.ino

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
