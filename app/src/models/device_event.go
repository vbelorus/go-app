package models

import (
	"fmt"
	"strings"
	"time"
)

type DeviceEvent struct {
	ClientTime  ClientTime `json:"client_time" validate:"required"`
	DeviceId    string     `json:"device_id" validate:"required"`
	DeviceOs    string     `json:"device_os" validate:"required"`
	Session     string     `json:"session" validate:"required"`
	Sequence    int32      `json:"sequence" validate:"gte=0"`
	Event       string     `json:"event" validate:"required"`
	ParamInt    int32      `json:"param_int" validate:"gte=0"`
	ParamString string     `json:"param_str" validate:"required"`
	Ip          string     `validate:"required"`
	ServerTime  time.Time  `validate:"required"`
}

type ClientTime time.Time

const ctLayout = "2006-01-02 15:04:05"

func (ct *ClientTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(ctLayout, s)
	*ct = ClientTime(nt)
	return
}

func (ct ClientTime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *ClientTime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(ctLayout))
}

func (ct *ClientTime) GetTime() time.Time {
	return time.Time(*ct)
}
