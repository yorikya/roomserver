package devices

import (
	"fmt"
	"log"
	"sync"

	"github.com/smira/go-statsd"
)

var (
	LUM_1900  = RGBColor{255, 147, 41, 50}
	LUM_2600  = RGBColor{255, 197, 143, 50}
	LUM_2850  = RGBColor{255, 214, 170, 50}
	LUM_3200  = RGBColor{255, 241, 224, 50}
	LUM_5200  = RGBColor{255, 250, 244, 50}
	LUM_5400  = RGBColor{255, 255, 251, 50}
	LUM_6000  = RGBColor{255, 255, 255, 50}
	LUM_7000  = RGBColor{201, 226, 255, 50}
	LUM_20000 = RGBColor{64, 156, 255, 50}

	WHITE = RGBColor{255, 255, 255, 0}
)

type RGBColor struct {
	RColor, GColor, BColor, Fade int
}

func (c RGBColor) ToCMD() string {
	return fmt.Sprintf("%d,%d,%d,%d", c.RColor, c.GColor, c.BColor, c.Fade)
}

type RGBStrip struct {
	ID       string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
	mu       *sync.Mutex
}

func (s *RGBStrip) CreateCMD(cmd string) (string, error) {
	switch cmd {
	case "1900L":
		return LUM_1900.ToCMD(), nil
	case "2600L":
		return LUM_2600.ToCMD(), nil
	case "2850L":
		return LUM_2850.ToCMD(), nil
	case "3200L":
		return LUM_3200.ToCMD(), nil
	case "5200L":
		return LUM_5200.ToCMD(), nil
	case "5400L":
		return LUM_5400.ToCMD(), nil
	case "6000L":
		return LUM_6000.ToCMD(), nil
	case "7000L":
		return LUM_7000.ToCMD(), nil
	case "20000L":
		return LUM_20000.ToCMD(), nil
	}
	return WHITE.ToCMD(), nil
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
	log.Println("RGB SetValue need implement this function")
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
		ID:     id,
		Name:   "rgbstrip",
		Sensor: sensor,
		mu:     &sync.Mutex{},
	}
}
