package vehicle

import (
	"encoding/binary"
)

// EncryptMethod used to differentiate encrypt method for Packet
type EncryptMethod uint8

const (
	// EytNone means no encryption
	EytNone EncryptMethod = 0x01
	// EytRSA use RSA encryption method
	EytRSA EncryptMethod = 0x02
	// EytSM2 use SM2 encryption method
	EytSM2 EncryptMethod = 0x03
	// EytErr indicate error during encryption
	EytErr EncryptMethod = 0xFE
	// EytInvaild indicate  improper encryption
	EytInvaild EncryptMethod = 0xFF
)

// Packet is the package send to the server, it includes the information of the RTDU
// and info about the truck
// BCC is will be automatically generated when encoding Packet to []byte
type Packet struct {
	Start       string
	CommandUnit uint8
	VIN         string
	CALID       uint8
	Encryption  uint8
	DULength    uint16
	DU          DataUnit
	BCC         uint8
}

// Encode will enocde the data packet to []byte,
// this []byte is required by the regulation
func (p *Packet) Encode() []byte {
	b := make([]byte, 24)
	copy(b, p.Start)
	b[2] = byte(p.CommandUnit)
	copy(b[3:], p.VIN)
	b[20] = byte(p.CALID)
	b[21] = byte(p.Encryption)
	binary.BigEndian.PutUint16(b[22:], p.DULength)
	du := p.DU.encode()
	b = append(b, du...)
	bcc := bccCal(b[2:])
	b = append(b, bcc)
	return b
}

// Encrypt will encrypt the data using the EncryptMethod
// TODO:
func (p *Packet) Encrypt(em EncryptMethod) []byte {
	b := p.Encode()
	return b
}

func bccCal(bs []byte) byte {
	var r byte
	l := len(bs)
	if l == 0 {
		return r
	}
	if l == 1 {
		return bs[0]
	}
	r = bs[0] ^ bs[1]
	for i := 2; i < l; i++ {
		r = r ^ bs[i]
	}
	return r
}
