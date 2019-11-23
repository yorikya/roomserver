package devices

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/smira/go-statsd"
)

type MovementSensor struct {
	id       string
	name     string
	valueStr string
	value    int64
	mu       *sync.Mutex
}

func (s *MovementSensor) GetID() string {
	return s.id
}

func (s *MovementSensor) GetName() string {
	return s.name
}

func (s *MovementSensor) SetValue(newValstr string) error {
	n, err := strconv.ParseInt(newValstr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed parse int in SetValue movement sensor function, %s", err)
	}
	s.mu.Lock()
	s.valueStr = newValstr
	s.value = n
	s.mu.Unlock()
	return nil
}

func (s *MovementSensor) GetValueStr() string {
	return s.valueStr
}

func (s *MovementSensor) SendStats(c *statsd.Client) {
	c.Gauge(fmt.Sprintf("%s.%s", s.id, s.name), s.value)
}

func NewMovementSensor(name string) *MovementSensor {
	return &MovementSensor{
		id:   "movesensor",
		name: name,
		mu:   &sync.Mutex{},
	}
}
