package devices

import (
	"fmt"
	"log"

	"github.com/smira/go-statsd"
)

type MainLight struct {
	RoomName string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
}

func (_ *MainLight) InRangeThreshold() bool {
	return false
}

func (s *MainLight) TurnOff() {
	log.Println("MainLight turn off")
}

func (_ *MainLight) CreateCMD(cmd string) (string, string, []string, error) {
	return cmd, CUSTOM, []string{}, nil
}

func (s *MainLight) GetRoomName() string {
	return s.RoomName
}

func (s *MainLight) GetSensor() string {
	return s.Sensor
}

func (s *MainLight) GetName() string {
	return s.Name
}

func (s *MainLight) GetOptions(_ string) []string {
	return []string{}
}

func (s *MainLight) SetValue(newValstr string) error {

	return nil
}

func (s *MainLight) GetValueStr() string {
	return s.ValueStr
}

func (s *MainLight) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.Name, s.Sensor), s.value)
}

func NewLight(roomName, sensor string) *MainLight {
	return &MainLight{
		RoomName: roomName,
		Name:     Light,
		Sensor:   sensor,
		ValueStr: "UNSET",
	}
}
