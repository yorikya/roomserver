package client

import (
	"fmt"
	"log"
	"strconv"
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

	//DHT22 sensor
	dht         = "dht"
	humidity    = "Humidity"
	temperature = "Temperature"

	//Movement sensor
	movesensor = "movesensor"
)

type Client struct {
	ClientID, AirCond, LightMain, LightSec string
	humidity, temperature                  float64
	MovementSen                            int64
	LastSeen                               time.Time
	stats                                  *statsd.Client
	mu                                     *sync.Mutex
	sensors                                []devices.Sensor
}

func (c *Client) update(msg IncomingMessage) {
	switch msg.GetDeviceID() {
	case dht:
		val, err := strconv.ParseFloat(msg.GetSensorValue(), 64)
		if err != nil {
			log.Printf("error parse float in update client %s\n", err)
			return
		}
		switch msg.GetSensorName() {
		case humidity:
			c.mu.Lock()
			c.humidity = val
			c.mu.Unlock()
		case temperature:
			c.mu.Lock()
			c.temperature = val
			c.mu.Unlock()
		}
	case movesensor:
		n, err := strconv.ParseInt(msg.GetSensorValue(), 10, 64)
		if err != nil {
			log.Printf("failed parse int in update client function, %s", err)
			return
		}
		c.mu.Lock()
		c.MovementSen = n
		c.mu.Unlock()
	}
}

func (c *Client) sendDeviceStats(msg IncomingMessage) {
	switch msg.GetDeviceID() {
	case dht:
		val, err := strconv.ParseFloat(msg.GetSensorValue(), 64)
		if err != nil {
			log.Printf("error parse float in OutTopic %s\n", err)
			return
		}
		c.stats.FGauge(fmt.Sprintf("%s.%s", msg.GetDeviceID(), msg.GetSensorName()), val)
	case movesensor:
		n, err := strconv.ParseInt(msg.GetSensorValue(), 10, 64)
		if err != nil {
			log.Printf("failed parse int in send stats client function, %s", err)
			return
		}
		c.stats.Gauge(msg.GetDeviceID(), n)
	}
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
	log.Printf("Create a new client '%s', sensors %s \n", clientID, s)
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
