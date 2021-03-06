package server

import (
	"fmt"
	"log"
	"time"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/config"
	"github.com/yorikya/roomserver/devices"
	"github.com/yorikya/roomserver/hub"
)

type Server struct {
	clients      map[string]*client.Client
	RoomHub      *hub.Hub
	clientOnline *time.Ticker
	stopTicker   chan bool
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
		clients:      make(map[string]*client.Client),
		RoomHub:      hub.NewHub(),
		stopTicker:   make(chan bool),
		clientOnline: time.NewTicker(time.Second * 10),
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
		s.clients[room.Name] = client.NewClient(room.Name, roomCfg.ShutDownDevices, devices.NewDevices(room.Name, roomCfg)...)

		log.Printf("client: %s was added, data: %+v\n", room.Name, s.clients[room.Name])
	}

	go func(s *Server) {
		log.Println("start server online ticker.")
		for {
			select {
			case <-s.stopTicker:
				log.Println("close server ticker")
				return
			case _ = <-s.clientOnline.C:
				for n, c := range s.clients {
					tmpStatus := c.OnLine
					c.OnLine = time.Since(c.LastSeen) < time.Second*20
					if tmpStatus != c.OnLine {
						log.Println("client ", n, "change online status from", tmpStatus, "to", c.OnLine)
						if c.OnLine {
							s.BrodcastHTMLClients(fmt.Sprintf("%s/status/OnLineGreen.jpeg", c.ClientID))
						} else {
							s.BrodcastHTMLClients(fmt.Sprintf("%s/status/OnLineRed.jpeg", c.ClientID))
						}

					}
				}

			}
		}
	}(s)

	return s
}

func (s *Server) Close() {
	for _, c := range s.clients {
		c.Close()
	}
	s.stopTicker <- true
	s.clientOnline.Stop()
	log.Println("closing server")
}

func (s *Server) BrodcastHTMLClients(msgs ...string) {
	for _, msg := range msgs {
		if err := s.RoomHub.Brodcast(msg); err != nil {
			log.Printf("failed broadcast message: '%s'", msg)
		}
	}
}
