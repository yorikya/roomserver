package server

import (
	"log"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/config"
	"github.com/yorikya/roomserver/devices"
	"github.com/yorikya/roomserver/hub"
)

type Server struct {
	clients map[string]*client.Client
	RoomHub *hub.Hub
}

func (s *Server) GetClients() map[string]*client.Client {
	return s.clients
}

func (s *Server) GetRooms() []string {
	set := make(map[string]bool)
	for _, c := range s.clients {
		set[c.GetRoomID()] = true
	}
	rooms := []string{}
	for r := range set {
		rooms = append(rooms, r)
	}
	return rooms
}

func (s *Server) GetClient(name string) *client.Client {
	c, ok := s.clients[name]
	if !ok {
		log.Printf("clientID '%s' does not exists\n", name)
		return nil
	}
	return c
}

func NewServer(configPath string) *Server {
	s := &Server{
		clients: make(map[string]*client.Client),
		RoomHub: hub.NewHub(),
	}
	cfgRooms, err := config.ParseRooms(configPath)
	if err != nil {
		log.Printf("failed parse file %s", err)
	}
	for _, room := range cfgRooms.Rooms {
		roomCfg := cfgRooms.GetRoom(room.Name)
		if roomCfg == nil {
			log.Printf("dosen't have config for '%s'", room.Name)
			continue
		}
		s.clients[room.Name] = client.NewClient(room.Name, devices.NewDevices(room.Name, roomCfg)...)
		log.Printf("client: %s was added, data: %+v\n", room.Name, s.clients[room.Name])
	}

	return s
}

func (s *Server) Close() {

}
