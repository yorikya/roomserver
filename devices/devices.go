package devices

import (
	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/config"
)

const (
	//DHT22 sensor
	dht = "dht"

	//Movement sensor
	movesensor = "movesensor"
)

type Sensor interface {
	GetID() string
	GetName() string
	GetValueStr() string
	SetValue(string) error
	SendStats(*statsd.Client)
}

func NewDevices(roomCfg *config.Room) []Sensor {
	sens := []Sensor{}
	for _, sensor := range roomCfg.Sensors {
		switch sensor.ID {
		case dht:
			sens = append(sens, NewHDTSensor(sensor.Name))
		case movesensor:
			sens = append(sens, NewMovementSensor(sensor.Name))
		}
	}
	return sens
}
