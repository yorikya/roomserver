package devices

import (
	"github.com/smira/go-statsd"
)

type Sensor interface {
	GetID() string
	GetName() string
	GetValueStr() string
	SetValue(string) error
	SendStats(*statsd.Client)
}
