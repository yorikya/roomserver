package devices

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/smira/go-statsd"
)

type HDTSensor struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
	mu       *sync.Mutex
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

func (s *HDTSensor) SetValue(newValstr string) error {
	newValue, err := strconv.ParseFloat(newValstr, 64)
	if err != nil {
		return fmt.Errorf("NewHDTSensor SetValue error parse float %s", err)
	}
	s.mu.Lock()
	s.ValueStr = newValstr
	s.value = newValue
	s.mu.Unlock()
	return nil
}

func (s *HDTSensor) GetValueStr() string {
	return s.ValueStr
}

func (s *HDTSensor) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.Name, s.Sensor), s.value)
}

func NewHDTSensor(id, sensor string) *HDTSensor {
	return &HDTSensor{
		ID:     id,
		Name:   "dht",
		Sensor: sensor,
		mu:     &sync.Mutex{},
	}
}
