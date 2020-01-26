package devices

import (
	"fmt"
	"math"
	"strconv"

	"github.com/smira/go-statsd"
)

type HDTSensor struct {
	RoomName         string
	Name             string
	Sensor           string
	ValueStr         string
	value, goodRange float64
}

func (s *HDTSensor) InRangeThreshold() bool {
	if math.Abs(s.value-s.goodRange) <= 2.0 {
		return true
	}
	return false
}

func (s *HDTSensor) TurnOff() {

}

func (_ *HDTSensor) CreateCMD(cmd string) (string, string, []string, error) {
	return cmd, CUSTOM, []string{}, nil
}

func (s *HDTSensor) GetRoomName() string {
	return s.RoomName
}

func (s *HDTSensor) GetSensor() string {
	return s.Sensor
}

func (s *HDTSensor) GetName() string {
	return s.Name
}

func (s *HDTSensor) GetOptions(_ string) []string {
	return []string{}
}

func (s *HDTSensor) SetValue(newValstr string) error {
	newValue, err := strconv.ParseFloat(newValstr, 64)
	if err != nil {
		return fmt.Errorf("NewHDTSensor SetValue error parse float %s", err)
	}
	s.ValueStr = newValstr
	s.value = newValue
	return nil
}

func (s *HDTSensor) GetValueStr() string {
	return s.ValueStr
}

func (s *HDTSensor) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.Name, s.Sensor), s.value)
}

func NewDHTHumiditySensor(roomName, sensor string, goodRange float64) *HDTSensor {
	return &HDTSensor{
		RoomName:  roomName,
		Name:      DHT_Humidity,
		Sensor:    sensor,
		ValueStr:  "UNSET",
		goodRange: goodRange,
	}
}

func NewDHTTemperatureSensor(roomName, sensor string, goodRange float64) *HDTSensor {
	return &HDTSensor{
		RoomName:  roomName,
		Name:      DHT_Temperature,
		Sensor:    sensor,
		ValueStr:  "UNSET",
		goodRange: goodRange,
	}
}
