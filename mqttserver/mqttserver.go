package mqttserver

import (
	"fmt"
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"

	"github.com/yorikya/roomserver/client"
	"github.com/yorikya/roomserver/devices"
)

type Server struct {
	mqttAdaptor *mqtt.Adaptor
	clients     map[string]*client.Client
	mqttbot     *gobot.Robot
}

func (s *Server) GetClients() map[string]*client.Client {
	return s.clients
}

func (s *Server) Publish(topic string, message []byte) bool {
	log.Printf("publish message %s to topic '%s'", string(message), topic)
	return s.mqttAdaptor.Publish(topic, message)
}

func NewServer(roomNames ...string) *Server {
	s := &Server{
		mqttAdaptor: mqtt.NewAdaptor("tcp://0.0.0.0:1883", "serve"),
		clients:     make(map[string]*client.Client),
	}
	for _, roomName := range roomNames {
		s1 := devices.NewHDTSensor("Humidity")
		s2 := devices.NewHDTSensor("Temperature")
		s3 := devices.NewMovementSensor("state")
		s.clients[roomName] = client.NewClient(roomName, s1, s2, s3)
		log.Printf("client: %s was added, data: %+v\n", roomName, s.clients[roomName])
	}

	work := func() {
		for id, c := range s.clients {
			s.mqttAdaptor.On(fmt.Sprintf("%sOutTopic", id), func(msg mqtt.Message) {
				m := client.NewIncomingMessage(string(msg.Payload()))
				c.Update(m)

				log.Printf("client '%s' OutTopic: %s\n", id, m)

			})

			s.mqttAdaptor.On(fmt.Sprintf("%sInTopic", id), func(msg mqtt.Message) {
				log.Printf("client '%s' sInTopic %s", id, string(msg.Payload()))

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
