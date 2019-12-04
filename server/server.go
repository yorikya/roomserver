package server

import (
	"log"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/config"
	"github.com/yorikya/roomserver/devices"
)

type Server struct {
	clients map[string]*client.Client
}

func (s *Server) GetClients() map[string]*client.Client {
	return s.clients
}

func NewServer(roomNames ...string) *Server {
	s := &Server{
		clients: make(map[string]*client.Client),
	}
	cfgRooms, err := config.ParseRooms("config/rooms.json")
	if err != nil {
		log.Printf("failed parse file %s", err)
	}
	for _, roomName := range roomNames {
		roomCfg := cfgRooms.GetRoom(roomName)
		if roomCfg == nil {
			continue
		}
		s.clients[roomName] = client.NewClient(roomName, devices.NewDevices(roomCfg)...)
		log.Printf("client: %s was added, data: %+v\n", roomName, s.clients[roomName])
	}

	return s
}

func (s *Server) Close() {

}
