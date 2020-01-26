package devices

import (
	"fmt"

	"github.com/smira/go-statsd"
)

type MainDoor struct {
	RoomName string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
}

func (_ *MainDoor) InRangeThreshold() bool {
	return false
}

func (s *MainDoor) TurnOff() {

}

func (_ *MainDoor) CreateCMD(cmd string) (string, string, []string, error) {
	return cmd, CUSTOM, []string{}, nil
}

func (s *MainDoor) GetRoomName() string {
	return s.RoomName
}

func (s *MainDoor) GetSensor() string {
	return s.Sensor
}

func (s *MainDoor) GetName() string {
	return s.Name
}

func (s *MainDoor) GetOptions(_ string) []string {
	return []string{}
}

func (s *MainDoor) SetValue(newValstr string) error {

	return nil
}

func (s *MainDoor) GetValueStr() string {
	return s.ValueStr
}

func (s *MainDoor) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.Name, s.Sensor), s.value)
}

func NewDoor(roomName, sensor string) *MainDoor {
	return &MainDoor{
		RoomName: roomName,
		Name:     Door,
		Sensor:   sensor,
		ValueStr: "UNSET",
	}
}
