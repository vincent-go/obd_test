package mock

// diagnoseSupport from response data of mcm and acm
// encode the result as defined in GB17691
func diagnoseSupport(data [4]byte) uint16 {
	// 1 "Catalyst monitoring status"
	// 2 "Heated catalyst monitoring"
	// 3 "Evaporative system monitoring status"
	// 4 "Secondary air system monitooring status"
	// 5 "A/C system refrigerant moitoring status"
	// 6 "Exhaust gas sensor monitoring status"
	// 7 "Exhasut gas sensor heater monitoring status"
	// 8 "EGR/VVT system monitoring"
	// 9 "Cold start aid system monitoring status"
	// 10 "Boost pressure control system monitoring status"
	// 11 "Diesel Particulate Filter (DPF) monitoring status"
	// 12 "NOx converting catalyst and/or NOx absorber monitoring status"
	// 13 "NMHC converting catalyst monitoring status"
	// 14 "Misfire monitoring support"
	// 15 "Fuel system monitoring support"
	// 16 "Comprehensive component monitoring support"
	var result uint16
	bs := make([]byte, 16)
	bs[5] = data[2] & 0b100000
	bs[7] = data[2] & 0b10000000
	bs[9] = data[2] & 0b1000
	bs[10] = data[2] & 0b1000000
	bs[11] = data[2] & 0b10
	bs[12] = data[2] & 1
	bs[13] = data[1] & 1
	bs[14] = data[1] & 0b10
	bs[15] = data[1] & 0b100
	for i, d := range bs {
		result = result | uint16(d<<i)
	}
	return result
}

// diagnoseReadiness from response data of mcm and acm
// encode the result as defined in GB17691 just like diagnose support.
func diagnoseReadiness(data [4]byte) uint16 {
	var result uint16
	bs := make([]byte, 16)
	bs[5] = (data[2] & 0b100000) & data[3]
	bs[7] = (data[2] & 0b10000000) & data[3]
	bs[9] = (data[2] & 0b1000) & data[3]
	bs[10] = (data[2] & 0b1000000) & data[3]
	bs[11] = (data[2] & 0b10) & data[3]
	bs[12] = (data[2] & 1) & data[3]
	bs[13] = (data[1] & 1) & (data[1] & 0b10000)
	bs[14] = (data[1] & 0b10) & (data[1] & 0b100000)
	bs[15] = (data[1] & 0b100) & (data[1] & 0b1000000)
	for i, d := range bs {
		result = result | uint16(d<<i)
	}
	return result
}

// sourceSupportAndReadiness calculate [4]bytes from diagnose support status and readiness
func sourceSupportAndReadiness(support, readiness uint16) [4]byte {
	bs := make([]byte, 32)
	bs[18] = byte(support & 0b10000000000)
	bs[16] = byte(support & 0b100000000)
	bs[20] = byte(support & 0b1000000)
	bs[17] = byte(support & 0b100000)
	bs[22] = byte(support & 0b10000)
	bs[23] = byte(support & 0b1000)
	bs[15] = byte(support & 0b100)
	bs[14] = byte(support & 0b10)
	bs[13] = byte(support & 1)
	bs[26] = readinessBit(bs[18], byte(readiness&0b10000000000))
	bs[24] = readinessBit(bs[16], byte(readiness&0b100000000))
	bs[28] = readinessBit(bs[20], byte(readiness&0b1000000))
	bs[25] = readinessBit(bs[17], byte(readiness&0b100000))
	bs[30] = readinessBit(bs[22], byte(readiness&0b10000))
	bs[31] = readinessBit(bs[23], byte(readiness&0b1000))
	bs[11] = readinessBit(bs[15], byte(readiness&0b100))
	bs[10] = readinessBit(bs[14], byte(readiness&0b10))
	bs[9] = readinessBit(bs[13], byte(readiness&1))
	result := [4]byte{}
	for i := 0; i < 4; i++ {
		b := concatBits(bs[8*i : 8*i+8])
		result[i] = b
	}
	return result
}

func concatBits(bits []byte) byte {
	var result byte
	if len(bits) != 8 {
		panic("input bits not enough to form a byte")
	}
	for i, b := range bits {
		result = result | b<<i
	}
	return result
}

// readinessBit set readiness bit accoridng to readiness bit and support bit
func readinessBit(supportBit, readinessBit byte) byte {
	if supportBit == 0 {
		// if item not support, set readiness as 1
		return 1
	}
	if readinessBit == 1 {
		// if item supported, and readiness bit is 1, set readiness bit as 1
		return 1
	}
	return 0
}

func mergeDiagnoseSupport(mcm, acm uint16) uint16 {
	return mcm | acm
}

func mergeReadiness(mcmSupport, mcmReadiness, acmSupport, acmReadiness uint16) uint16 {
	return (mcmSupport & mcmReadiness) | (acmSupport & acmReadiness)
}
