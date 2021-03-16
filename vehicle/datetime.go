package vehicle

import (
	"bytes"
	"encoding/binary"
	"time"
)

// DateTime defines the datetime format requested by the regulation
type DateTime struct {
	Year  uint8 // only contains the last 2 digits in the year, for example, year 2021 will be recorded as 21
	Month uint8
	Day   uint8
	Hour  uint8
	Min   uint8
	Sec   uint8
}

// encode will encode datetime to []byte
func (dt *DateTime) encode() []byte {
	buf := new(bytes.Buffer)
	b := make([]byte, 6)
	binary.Write(buf, binary.BigEndian, dt)
	buf.Read(b)
	return b
}

// FetchDateTime get the currect datatime with the format demanded by the regulation
func FetchDateTime() DateTime {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return DateTime{
		Year:  uint8(year % 100),
		Month: uint8(month),
		Day:   uint8(day),
		Hour:  uint8(hour),
		Min:   uint8(min),
		Sec:   uint8(sec),
	}
}
