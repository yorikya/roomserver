package client

import (
	"fmt"
	"log"
	"strings"

	"time"

	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/devices"
)

type Client struct {
	roomID, ClientID, IPstr string
	LastSeen                time.Time
	stats                   *statsd.Client
	Devices                 []devices.Device
}

func NewClient(clientID string, d ...devices.Device) *Client {
	log.Printf("Create a new client '%s'\n", clientID)
	return &Client{
		roomID:   strings.Split(clientID, "_")[0],
		Devices:  d,
		ClientID: clientID,
		stats:    statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400), statsd.MetricPrefix(fmt.Sprintf("home.%s.", clientID))),
	}
}

func (c *Client) GetRoomID() string {
	return c.roomID
}

func (c *Client) UpdateIPstr(ip string) {
	if c.IPstr != ip {
		log.Printf("clientID: %s, change IP from: '%s', to: '%s'\n", c.ClientID, c.IPstr, ip)
		c.IPstr = ip
	}
}

func (c *Client) Update(device, value string) {
	//device, sensor, value => hdt/Humidity/30.40
	c.stats.Incr("update", 1)
	c.LastSeen = time.Now()
	if device == "keepalive" {
		log.Println("client: ", c.ClientID, "keepalive message")
		return
	}

	s := c.GetDeviceByName(device)
	if s != nil {
		s.SetValue(value)
		s.SendStats(c.stats)
		return
	}
	log.Printf("clientID: %s, does not has device: %s\n", c.ClientID, device)
}

func (c *Client) GetDeviceByName(name string) devices.Device {
	for _, d := range c.Devices {
		if d.GetName() == name {
			return d
		}
	}
	return nil
}
