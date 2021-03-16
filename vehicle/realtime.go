package vehicle

import (
	"bytes"
	"encoding/binary"
)

// RTBody is the data type of collected data from the vehicle,
// it includes all required engine/vehicle real-DateTime data
// several RTBody will be composed together and send to server
type RTBody struct {
	Type            uint8
	ST              DateTime
	VehicleSpeed    uint16
	AmbientP        uint8
	NetTorque       uint8
	Fraction        uint8
	EngineSpeed     uint16
	FuelConsumption uint16
	RawNOx          uint16
	TpNOx           uint16
	Urea            uint8
	AirMF           uint16
	SCRTempIn       uint16
	SCRTempOut      uint16
	DPFDeltaP       uint16
	CoolantTemp     uint8
	FuelTankLevel   uint8
	LocationStatus  uint8
	Latitude        uint32
	Longitude       uint32
	Mileage         uint32
}

// encodeBody will encode the RTBody to binary with bigEndian byte order
func (rt *RTBody) encodeBody() []byte {
	buf := new(bytes.Buffer)
	b := make([]byte, 44)
	binary.Write(buf, binary.BigEndian, rt)
	buf.Read(b)
	return b
}

