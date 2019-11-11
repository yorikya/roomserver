package mqttserver

import (
	"encoding/json"
	"fmt"
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"

	"github.com/yorikya/roomserver/client"
)

type Server struct {
	mqttAdaptor *mqtt.Adaptor
	clients     map[string]*client.Client
	mqttbot     *gobot.Robot
}

func (s *Server) GetClients() map[string]*client.Client {
	return s.clients
}

func (s *Server) Publish(roomID string, message []byte) bool {
	return s.mqttAdaptor.Publish(fmt.Sprintf("%sInTopic", roomID), message)
}

func NewServer(roomNames ...string) *Server {
	s := &Server{
		mqttAdaptor: mqtt.NewAdaptor("tcp://0.0.0.0:1883", "roomserve"),
		clients:     make(map[string]*client.Client),
	}
	for _, roomName := range roomNames {
		s.clients[roomName] = client.NewClient(roomName)
		log.Printf("client: %s was added, data: %+v\n", roomName, s.clients[roomName])
	}

	work := func() {
		for id, c := range s.clients {
			s.mqttAdaptor.On(fmt.Sprintf("%sOutTopic", id), func(msg mqtt.Message) {
				m := client.OutMessage{}
				err := json.Unmarshal(msg.Payload(), &m)
				if err != nil {
					log.Println("error:", err)
					return
				}
				switch m.Action {
				case "update":
					log.Printf("client '%s' update state: %+v", id, m)
					c.UpdateState(m)
				default:
					log.Printf("'%s' is unknown action", m.Action)
				}
			})
		}
	}

	s.mqttbot = gobot.NewRobot("mqttBot",
		[]gobot.Connection{s.mqttAdaptor},
		work,
	)

	return s
}

func (s *Server) Start() {
	log.Println("starting mqtt server...")
	s.mqttbot.Start()
}

func (s *Server) Close() {
	s.mqttbot.Stop()
}
