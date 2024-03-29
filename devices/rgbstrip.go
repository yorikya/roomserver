package devices

import (
	"fmt"
	"log"

	"github.com/smira/go-statsd"
	"github.com/yorikya/roomserver/style"
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
		RGBColor{"ON", 255, 255, 255, 50},
		RGBColor{"OFF", 0, 0, 0, 50},
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
	RoomName string
	Name     string
	Sensor   string
	ValueStr string
	value    float64
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

func (s *RGBStrip) CreateCMD(cmd string) (string, string, []string, error) {
	col, ok := s.getColor(cmd)
	if ok {
		msg := UpdateMsg(s.RoomName, s.Name, style.GetTextStyle(style.StylColGreen), col.GetTag())
		return col.ToCMD(), col.GetTag(), []string{msg}, nil
	}
	return cmd, CUSTOM, []string{}, nil
}

func (s *RGBStrip) GetRoomName() string {
	return s.RoomName
}

func (s *RGBStrip) GetSensor() string {
	return s.Sensor
}

func (s *RGBStrip) GetName() string {
	return s.Name
}

func (s *RGBStrip) SetValue(newValstr string) error {
	s.ValueStr = newValstr
	return nil
}

func (s *RGBStrip) GetValueStr() string {
	return s.ValueStr
}

func (s *RGBStrip) GetOptions(_ string) []string {
	copt := []string{}
	for _, c := range stripColors {
		copt = append(copt, c.GetTag())
	}
	return copt
}

func (s *RGBStrip) SendStats(c *statsd.Client) {
	log.Println("RGB need implement this function")
}

func (_ *RGBStrip) InRangeThreshold() bool {
	return false
}

func (s *RGBStrip) TurnOff() {
	log.Println("RGBStrip turn off")
}

func NewRGBStrip(roomName, sensor string) *RGBStrip {
	return &RGBStrip{
		RoomName: roomName,
		Name:     RGBstrip,
		Sensor:   sensor,
		colors:   stripColors,
		ValueStr: "UNSET",
	}
}
