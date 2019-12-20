package devices

import (
	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/config"
)

const (
	//DHT22 sensor
	dht = "dht"

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
	CreateCMD(string) (string, error)
	SendStats(*statsd.Client)
}

func NewDevices(roomName string, roomCfg *config.Room) []Device {
	sens := []Device{}
	for _, device := range roomCfg.Devices {
		switch device.Name {
		case dht:
			sens = append(sens, NewHDTSensor(roomName, device.Sensor))
		case rgbstrip:
			sens = append(sens, NewRGBStrip(roomName, device.Sensor))
		case ir_ac_aircool:
			sens = append(sens, NewIRACAirCool(roomName, device.Sensor))
		}
	}
	return sens
}
