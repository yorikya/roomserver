package devices

import (
	"fmt"
	"strconv"

	"github.com/smira/go-statsd"
)

type HDTSensor struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
}

func (_ *HDTSensor) CreateCMD(cmd string) (string, string, error) {
	return cmd, CUSTOM, nil
}

func (s *HDTSensor) GetID() string {
	return s.ID
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

func NewDHTHumiditySensor(id, sensor string) *HDTSensor {
	return &HDTSensor{
		ID:       id,
		Name:     DHT_Humidity,
		Sensor:   sensor,
		ValueStr: "UNSET",
	}
}

func NewDHTTemperatureSensor(id, sensor string) *HDTSensor {
	return &HDTSensor{
		ID:       id,
		Name:     DHT_Temperature,
		Sensor:   sensor,
		ValueStr: "UNSET",
	}
}
