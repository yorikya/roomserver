package devices

import (
	"fmt"
	"sync"

	"github.com/smira/go-statsd"
)

type Camera struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
	mu       *sync.Mutex
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
	// newValue, err := strconv.ParseFloat(newValstr, 64)
	// if err != nil {
	// 	return fmt.Errorf("NewHDTSensor SetValue error parse float %s", err)
	// }
	// s.mu.Lock()
	// s.ValueStr = newValstr
	// s.value = newValue
	// s.mu.Unlock()
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
		Name:     "camera",
		Sensor:   sensor,
		mu:       &sync.Mutex{},
		ValueStr: "Not Connected",
	}
}
