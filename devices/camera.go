package devices

import (
	"fmt"

	"github.com/smira/go-statsd"
)

type Camera struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
}

func (_ *Camera) InRangeThreshold() bool {
	return false
}

func (_ *Camera) CreateCMD(cmd string) (string, string, error) {
	return cmd, CUSTOM, nil
}

func (s *Camera) GetID() string {
	return s.ID
}

func (s *Camera) GetSensor() string {
	return s.Sensor
}

func (s *Camera) GetName() string {
	return s.Name
}

func (s *Camera) GetOptions(_ string) []string {
	return []string{}
}

func (s *Camera) SetValue(newValstr string) error {

	return nil
}

func (s *Camera) GetValueStr() string {
	return s.ValueStr
}

func (s *Camera) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.Name, s.Sensor), s.value)
}

func NewCamera(id, sensor string) *Camera {
	return &Camera{
		ID:       id,
		Name:     Camera2MP,
		Sensor:   sensor,
		ValueStr: "Not Connected",
	}
}
