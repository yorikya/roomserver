package devices

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/smira/go-statsd"
)

type HDTSensor struct {
	id       string
	name     string
	valueStr string
	value    float64
	mu       *sync.Mutex
}

func (s *HDTSensor) GetID() string {
	return s.id
}

func (s *HDTSensor) GetName() string {
	return s.name
}

func (s *HDTSensor) SetValue(newValstr string) error {
	newValue, err := strconv.ParseFloat(newValstr, 64)
	if err != nil {
		return fmt.Errorf("NewHDTSensor SetValue error parse float %s", err)
	}
	s.mu.Lock()
	s.valueStr = newValstr
	s.value = newValue
	s.mu.Unlock()
	return nil
}

func (s *HDTSensor) GetValueStr() string {
	return s.valueStr
}

func (s *HDTSensor) SendStats(c *statsd.Client) {
	c.FGauge(fmt.Sprintf("%s.%s", s.id, s.name), s.value)
}

func NewHDTSensor(name string) *HDTSensor {
	return &HDTSensor{
		id:   "dht",
		name: name,
		mu:   &sync.Mutex{},
	}
}
