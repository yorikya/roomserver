package devices

import (
	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/config"
)

const (
	//Custom command
	CUSTOM = "CUSTOM"
	//DHT22 sensor
	dht_Humidity    = "dht_Humidity"
	dht_Temperature = "dht_Temperature"

	//RGB Strip
	rgbstrip = "rgbstrip"

	//AC Air Cool IR
	ir_ac_aircool = "ir_ac_aircool"
)

type Device interface {
	GetID() string
	GetName() string
	GetSensor() string
	GetValueStr() string
	SetValue(string) error
	CreateCMD(string) (string, string, error)
	SendStats(*statsd.Client)
}

func NewDevices(roomName string, roomCfg *config.Room) []Device {
	sens := []Device{}
	for _, device := range roomCfg.Devices {
		switch device.Name {
		case dht_Humidity:
			sens = append(sens, NewDHTHumiditySensor(roomName, device.Sensor))
		case dht_Temperature:
			sens = append(sens, NewDHTTemperatureSensor(roomName, device.Sensor))
		case rgbstrip:
			sens = append(sens, NewRGBStrip(roomName, device.Sensor))
		case ir_ac_aircool:
			sens = append(sens, NewIRACAirCool(roomName, device.Sensor))
		}
	}
	return sens
}
