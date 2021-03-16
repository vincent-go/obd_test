package vehicle

import "encoding/binary"

// DTC is 4 byte slice, normally it is just ASCII encoded number (to be confirmed)
type DTC [4]byte

func (d *DTC) encodeDTC() []byte {
	return d[:]
}

// OBDBody is one type of the collected data from the vehicle,
// it includes all OBD related data
type OBDBody struct {
	Type              uint8
	ST                DateTime
	Protocal          uint8
	MIL               uint8
	DiagnoseSupport   uint16
	DiagnoseReadiness uint16
	VIN               string
	CALID             string
	CVN               string
	IUPR              [18]uint16
	DTCCount          uint8
	DTCs              []DTC
}

// encode will encode OBDBody to binary with bigEndian byte order
func (o *OBDBody) encodeBody() []byte {
	l := (o.DTCCount * 4) + 103
	b := make([]byte, l)
	b[0] = byte(o.Type)
	copy(b[1:], o.ST.encode())
	b[7] = byte(o.Protocal)
	b[8] = byte(o.MIL)
	binary.BigEndian.PutUint16(b[9:], o.DiagnoseSupport)
	binary.BigEndian.PutUint16(b[11:], o.DiagnoseReadiness)
	copy(b[13:], o.VIN)
	copy(b[30:], o.CALID)
	copy(b[48:], o.CVN)
	for i := 0; i < 18; i++ {
		binary.BigEndian.PutUint16(b[66+2*i:68+2*i], o.IUPR[i])
	}
	b[102] = byte(o.DTCCount)
	for i, dtc := range o.DTCs {
		start := 103 + i*4
		copy(b[start:], dtc.encodeDTC())
	}
	return b
}
