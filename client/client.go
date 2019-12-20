package client

import (
	"fmt"
	"log"

	"sync"
	"time"

	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/devices"
)

type Client struct {
	ClientID, IPstr string
	LastSeen        time.Time
	stats           *statsd.Client
	mu              *sync.Mutex
	Devices         []devices.Device
}

func NewClient(clientID string, d ...devices.Device) *Client {
	log.Printf("Create a new client '%s'\n", clientID)
	return &Client{
		Devices:  d,
		mu:       &sync.Mutex{},
		ClientID: clientID,
		stats:    statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400), statsd.MetricPrefix(fmt.Sprintf("home.%s.", clientID))),
	}
}

func (c *Client) UpdateIPstr(ip string) {
	if c.IPstr != ip {
		log.Printf("clientID: %s, change IP from: '%s', to: '%s'\n", c.ClientID, c.IPstr, ip)
		c.mu.Lock()
		c.IPstr = ip
		c.mu.Unlock()
	}
}

func (c *Client) Update(device, sensor, value string) {
	//device, sensor, value => hdt/Humidity/30.40
	c.stats.Incr("update", 1)
	c.mu.Lock()
	c.LastSeen = time.Now()
	c.mu.Unlock()
	if device == "keepalive" {
		log.Println("client: ", c.ClientID, "keepalive message")
		return
	}

	s := c.getSensor(device, sensor)
	if s != nil {
		s.SetValue(value)
		s.SendStats(c.stats)
		return
	}
	log.Printf("clientID: %s, does not has device: %s, sensor: %s\n", c.ClientID, device, sensor)
}

func (c *Client) getSensor(name, sensor string) devices.Device {
	for _, d := range c.Devices {
		if d.GetName() == name && d.GetSensor() == sensor {
			return d
		}
	}
	return nil
}

func (c *Client) GetSensorByName(name string) devices.Device {
	for _, d := range c.Devices {
		if d.GetName() == name {
			return d
		}
	}
	return nil
}
