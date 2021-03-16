package vehicle

import (
	"encoding/binary"
	"log"
)

// DataUnit is the struct that contains real time data terminal got from CAN
// In the original regulation, the defination of the DataUnit is WRONG!
// This defination is from a draft version of the new regulation which most likely will be the one
// used in production environment.
// There is one more thing added to the DataUnit in the new regulation:
// signature of the DataUnit, TODO: add this when fully understand the mechanism
type DataUnit struct {
	// UT is the upload time of the DataUnit
	UT   DateTime
	SN   uint16
	Body []Body
	// Sign signature
}

// encode will encode DataUnit
func (du *DataUnit) encode() []byte {
	b := make([]byte, 8)
	copy(b[:6], du.UT.encode())
	binary.BigEndian.PutUint16(b[6:], du.SN)
	for _, body := range du.Body {
		b = append(b, body.encodeBody()...)
	}
	return b
}

// Body is the interface for realtime data stream and OBD data stream
// Body is the main data in DataUnit
type Body interface {
	encodeBody() []byte
}

// LengthCal calculate the length of the DataUnit
// TODO: test the function
func (du *DataUnit) LengthCal() uint16 {
	if du == nil {
		log.Fatal("DataUnit is nil?")
	}
	// DataUnit have 8 bytes header
	l := uint16(8)
	for _, body := range du.Body {
		switch by := body.(type) {
		case *OBDBody:
			l = l + uint16(by.DTCCount)*4 + 103
		case *RTBody:
			l += 44
		}
	}
	return l
}
