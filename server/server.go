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

func (s *Server) GetClient(name string) *client.Client {
	c, ok := s.clients[name]
	if !ok {
		log.Printf("clientID '%s' does not exists\n", name)
		return nil
	}
	return c
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
			log.Printf("dosen't have config for '%s'", roomName)
			continue
		}
		s.clients[roomName] = client.NewClient(roomName, devices.NewDevices(roomName, roomCfg)...)
		log.Printf("client: %s was added, data: %+v\n", roomName, s.clients[roomName])
	}

	return s
}

func (s *Server) Close() {

}
