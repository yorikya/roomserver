package devices

import (
	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/config"
)

const (
	//Custom command
	CUSTOM = "CUSTOM"
	//DHT22 sensor
	DHT_Humidity    = "dht_Humidity"
	DHT_Temperature = "dht_Temperature"

	//RGB Strip
	RGBstrip = "rgbstrip"

	//AC Air Cool IR
	IR_ac_aircool = "ir_ac_aircool"

	//Video camera
	Camera2MP = "camera"
)

type Device interface {
	GetID() string
	GetName() string
	GetSensor() string
	GetValueStr() string
	GetOptions(string) []string
	SetValue(string) error
	CreateCMD(string) (string, string, error)
	SendStats(*statsd.Client)
}

func NewDevices(roomName string, roomCfg *config.Room) []Device {
	sens := []Device{}
	for _, device := range roomCfg.Devices {
		switch device.Name {
		case DHT_Humidity:
			sens = append(sens, NewDHTHumiditySensor(roomName, device.Sensor))
		case DHT_Temperature:
			sens = append(sens, NewDHTTemperatureSensor(roomName, device.Sensor))
		case RGBstrip:
			sens = append(sens, NewRGBStrip(roomName, device.Sensor))
		case IR_ac_aircool:
			sens = append(sens, NewIRACAirCool(roomName, device.Sensor))
		case Camera2MP:
			sens = append(sens, NewCamera(roomName, device.Sensor))
		}
	}
	return sens
}
