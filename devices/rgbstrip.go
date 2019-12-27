package devices

import (
	"fmt"
	"log"
	"sync"

	"github.com/smira/go-statsd"
)

var (
	stripColors = []RGBColor{
		RGBColor{"LUM_1900", 255, 147, 41, 50},
		RGBColor{"LUM_2600", 255, 197, 143, 50},
		RGBColor{"LUM_2850", 255, 214, 170, 50},
		RGBColor{"LUM_3200", 255, 241, 224, 50},
		RGBColor{"LUM_5200", 255, 250, 244, 50},
		RGBColor{"LUM_5400", 255, 255, 251, 50},
		RGBColor{"LUM_6000", 255, 255, 255, 50},
		RGBColor{"LUM_7000", 201, 226, 255, 50},
		RGBColor{"LUM_20000", 64, 156, 255, 50},
		RGBColor{"WHITE", 255, 255, 255, 0},
	}
)

type RGBColor struct {
	tag                          string
	RColor, GColor, BColor, Fade int
}

func (c RGBColor) ToCMD() string {
	return fmt.Sprintf("%d,%d,%d,%d", c.RColor, c.GColor, c.BColor, c.Fade)
}

func (c RGBColor) GetTag() string {
	return c.tag
}

type RGBStrip struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
	mu       *sync.Mutex
	colors   []RGBColor
}

func (s *RGBStrip) getColor(tag string) (RGBColor, bool) {
	for _, col := range s.colors {
		if col.tag == tag {
			return col, true
		}
	}
	return RGBColor{}, false
}

func (s *RGBStrip) CreateCMD(cmd string) (string, string, error) {
	col, ok := s.getColor(cmd)
	if ok {
		return col.ToCMD(), col.GetTag(), nil
	}
	return cmd, CUSTOM, nil
}

func (s *RGBStrip) GetID() string {
	return s.ID
}

func (s *RGBStrip) GetSensor() string {
	return s.Sensor
}

func (s *RGBStrip) GetName() string {
	return s.Name
}

func (s *RGBStrip) SetValue(newValstr string) error {
	s.mu.Lock()
	s.ValueStr = newValstr
	s.mu.Unlock()
	return nil
}

func (s *RGBStrip) GetValueStr() string {
	return s.ValueStr
}

func (s *RGBStrip) SendStats(c *statsd.Client) {
	log.Println("RGB need implement this function")
}

func NewRGBStrip(id, sensor string) *RGBStrip {
	return &RGBStrip{
		ID:       id,
		Name:     "rgbstrip",
		Sensor:   sensor,
		mu:       &sync.Mutex{},
		colors:   stripColors,
		ValueStr: "UNSET",
	}
}
