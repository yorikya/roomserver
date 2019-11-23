package client

import (
	"fmt"
	"log"
	"strings"

	"sync"
	"time"

	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/devices"
)

const (
	//IncomingMessage index
	device = iota + 1
	sensor
	value
)

type Client struct {
	ClientID, AirCond, LightMain, LightSec string
	LastSeen                               time.Time
	stats                                  *statsd.Client
	mu                                     *sync.Mutex
	sensors                                []devices.Sensor
}

func (c *Client) Update(msg IncomingMessage) {
	//msg.origMsg: hdt/Humidity/30.40
	c.stats.Incr("metrics", 1)
	s := c.getSensor(msg.GetDeviceID(), msg.GetSensorName())
	if s != nil {
		s.SetValue(msg.GetSensorValue())
		s.SendStats(c.stats)
	}
}

func (c *Client) getSensor(id, name string) devices.Sensor {
	for _, s := range c.sensors {
		if s.GetID() == id && s.GetName() == name {
			return s
		}
	}
	return nil
}

func NewClient(clientID string, s ...devices.Sensor) *Client {
	log.Printf("Create a new client '%s'n", clientID)
	return &Client{
		sensors:  s,
		mu:       &sync.Mutex{},
		ClientID: clientID,
		stats:    statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400), statsd.MetricPrefix(fmt.Sprintf("home.%s.", clientID))),
	}
}

type IncomingMessage struct {
	origMsg, delimiter string
	msgList            []string
}

func (m IncomingMessage) GetDeviceID() string {
	return m.msgList[device]
}

func (m IncomingMessage) GetSensorName() string {
	return m.msgList[sensor]
}

func (m IncomingMessage) GetSensorValue() string {
	return m.msgList[value]
}

func NewIncomingMessage(msg string) IncomingMessage {
	delimiter := "/"
	return IncomingMessage{
		origMsg:   msg,
		msgList:   strings.Split(msg, delimiter),
		delimiter: delimiter,
	}
}
