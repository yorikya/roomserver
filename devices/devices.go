package devices

import (
	"fmt"
	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/config"
)

const (
	//Update message format
	updateMsgFmt = "%s/update/%s/%s/%s" // ex: room1/update/rgbstrip/text-success/LUM_1900

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

	//Main Door
	Door = "door"

	//Main Light
	Light = "light"
)

type Device interface {
	GetRoomName() string
	GetName() string
	GetSensor() string
	GetValueStr() string
	InRangeThreshold() bool
	GetOptions(string) []string
	SetValue(string) error
	CreateCMD(string) (string, string, []string, error)
	SendStats(*statsd.Client)
	TurnOff()
}

func NewDevices(roomName string, roomCfg *config.Room) []Device {
	sens := []Device{}
	for _, device := range roomCfg.Devices {
		switch device.Name {
		case DHT_Humidity:
			sens = append(sens, NewDHTHumiditySensor(roomName, device.Sensor, device.Threshold))
		case DHT_Temperature:
			sens = append(sens, NewDHTTemperatureSensor(roomName, device.Sensor, device.Threshold))
		case RGBstrip:
			sens = append(sens, NewRGBStrip(roomName, device.Sensor))
		case IR_ac_aircool:
			sens = append(sens, NewIRACAirCool(roomName, device.Sensor))
		case Camera2MP:
			sens = append(sens, NewCamera(roomName, device.Sensor))
		case Door:
			sens = append(sens, NewDoor(roomName, device.Sensor))
		case Light:
			sens = append(sens, NewLight(roomName, device.Sensor))
		}
	}
	return sens
}

func UpdateMsg(roomName, deviceID, textStyle, value string) string {
	return fmt.Sprintf(updateMsgFmt, roomName, deviceID, textStyle, value)
}
