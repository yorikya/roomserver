package httphandlers

import "fmt"

const (
	H2Fmt         = "<h2 class='%s'>%s</h2>"
	stylColGreen  = "green"
	stylColRed    = "red"
	stylColViolet = "violet"
	stylColGray   = "gray"
	stylColYellow = "yellow"
	stylColSky    = "sky"
	stylColBlack  = "black"
)

func getTextStyle(color string) string {
	switch color {
	case stylColGreen:
		return "text-success"
	case stylColRed:
		return "text-danger"
	case stylColViolet:
		return "text-primary"
	case stylColGray:
		return "text-secondary"
	case stylColYellow:
		return "text-warning"
	case stylColSky:
		return "text-info"
	case stylColBlack:
		return "text-dark"
	default:
		return "text-dark"
	}
}

func getStyleSuccessOrFailText(success bool) string {
	if success {
		return getTextStyle(stylColGreen)
	}
	return getTextStyle(stylColRed)
}

func newH2(color, value string) string {
	return fmt.Sprintf(H2Fmt, getTextStyle(color), value)
}

func newH2GreenSRedF(success bool, value string) string {
	if success {
		return newH2("green", value)
	}
	return newH2("red", value)
}
