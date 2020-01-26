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
	OnLine                  bool
	shutDownDevices         []string
}

func NewClient(clientID string, shutDownDevices []string, d ...devices.Device) *Client {
	log.Printf("Create a new client '%s'\n", clientID)
	return &Client{
		shutDownDevices: shutDownDevices,
		roomID:          strings.Split(clientID, "_")[0],
		Devices:         d,
		ClientID:        clientID,
		stats:           statsd.NewClient("localhost:8125", statsd.MaxPacketSize(1400), statsd.MetricPrefix(fmt.Sprintf("home.%s.", clientID))),
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

func (c *Client) isShutDownDevice(name string) bool {
	for _, d := range c.shutDownDevices {
		if d == name {
			return true
		}
	}
	return false
}

func (c *Client) RunScenario(scenario string) {
	if scenario == "shutdownall" {
		for _, d := range c.Devices {
			if c.isShutDownDevice(d.GetName()) {
				d.TurnOff()
			}
		}
	}
}

func (c *Client) GetDHTHumidity() devices.Device {
	return c.GetDeviceByName(devices.DHT_Humidity)
}

func (c *Client) GetDHTTemperature() devices.Device {
	return c.GetDeviceByName(devices.DHT_Temperature)
}

func (c *Client) GetRGBstrip() devices.Device {
	return c.GetDeviceByName(devices.RGBstrip)
}

func (c *Client) GetIR_ac_aircool() devices.Device {
	return c.GetDeviceByName(devices.IR_ac_aircool)
}

func (c *Client) GetCamera2MP() devices.Device {
	return c.GetDeviceByName(devices.Camera2MP)
}

func (c *Client) GetDoor() devices.Device {
	return c.GetDeviceByName(devices.Door)
}

func (c *Client) GetLight() devices.Device {
	return c.GetDeviceByName(devices.Light)
}

func (c *Client) Close() {
	log.Println("closing client", c.ClientID)
}
