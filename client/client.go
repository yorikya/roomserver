package client

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/smira/go-statsd"
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
)

type Client struct {
	ClientID, MovementSen, AirCond, LightMain, LightSec string
	Humidity, Temperature                               float64
	LastSeen                                            time.Time
	stats                                               *statsd.Client
}

func (c *Client) Update(msg IncomingMessage) {
	//msg.origMsg: hdt/Humidity/30.40
	c.stats.Incr("metrics", 1)
	switch msg.GetDeviceID() {
	case dht:
		val, err := strconv.ParseFloat(msg.GetSensorValue(), 64)
		if err != nil {
			log.Printf("error parse float in OutTopic %s\n", err)
			return
		}
		c.stats.FGauge(fmt.Sprintf("%s.%s", msg.GetDeviceID(), msg.GetSensorName()), val)
		switch msg.GetSensorName() {
		case humidity:
			c.Humidity = val
		case temperature:
			c.Temperature = val
		default:
			log.Printf("'%s' is uknown sensor name\n", msg.GetSensorName())
		}
	default:
		c.stats.Incr("device.unknown", 1)
		log.Printf("'%s' is unknown device", msg.GetDeviceID())

	}

}

func NewClient(clientID string) *Client {
	return &Client{
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
