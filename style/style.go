package style

import "fmt"

const (
	H2Fmt         = "<h2 class='%s'>%s</h2>"
	StylColGreen  = "green"
	StylColRed    = "red"
	StylColViolet = "violet"
	StylColGray   = "gray"
	StylColYellow = "yellow"
	StylColSky    = "sky"
	StylColBlack  = "black"
)

func GetTextStyle(color string) string {
	switch color {
	case StylColGreen:
		return "text-success"
	case StylColRed:
		return "text-danger"
	case StylColViolet:
		return "text-primary"
	case StylColGray:
		return "text-secondary"
	case StylColYellow:
		return "text-warning"
	case StylColSky:
		return "text-info"
	case StylColBlack:
		return "text-dark"
	default:
		return "text-dark"
	}
}

func SuccessOrFailText(success bool) string {
	if success {
		return GetTextStyle(StylColGreen)
	}
	return GetTextStyle(StylColRed)
}

func NewH2(color, value string) string {
	return fmt.Sprintf(H2Fmt, GetTextStyle(color), value)
}

func NewH2GreenSRedF(success bool, value string) string {
	if success {
		return NewH2(StylColGreen, value)
	}
	return NewH2(StylColRed, value)
}
