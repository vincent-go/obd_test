package mock

import "github.com/vincent-go/obd_test/vehicle"

// GenRTPkt will generate a fixed realtime Packet(except date time) for consumption
func GenRTPkt() vehicle.Packet {
	rtBody := vehicle.RTBody{
		Type:            2,
		ST:              vehicle.FetchDateTime(),
		VehicleSpeed:    200,
		AmbientP:        100,
		NetTorque:       85,
		Fraction:        12,
		EngineSpeed:     2000,
		FuelConsumption: 20,
		RawNOx:          12,
		TpNOx:           1,
		Urea:            2,
		AirMF:           1856,
		SCRTempIn:       367,
		SCRTempOut:      378,
		DPFDeltaP:       10,
		CoolantTemp:     87,
		FuelTankLevel:   43,
		LocationStatus:  1,
		Latitude:        0x12344576,
		Longitude:       0x67342345,
		Mileage:         0x00000008,
	}
	du := vehicle.DataUnit{
		UT:   vehicle.FetchDateTime(),
		SN:   123,
		Body: []vehicle.Body{&rtBody},
	}
	pkt := vehicle.Packet{
		Start:       "##",
		CommandUnit: 0x01,
		VIN:         "12345678901234567",
		CALID:       0x01,
		Encryption:  0x00,
		DULength:    du.LengthCal(),
		DU:          du,
		BCC:         0x56,
	}
	return pkt
}

// GenOBDData generate OBD data, including data from MCM and ACM
// some data only requested from MCM, in this case, ACM will not include these items (default value)
func GenOBDData() (mcm, acm vehicle.OBDBody) {
	mcm = vehicle.OBDBody{
		VIN:               "ABCDEFG0123456789",
		Protocal:          1,
		MIL:               0,
		DiagnoseSupport:   0b1111110111110111,
		DiagnoseReadiness: 0b1111101111111011,
		CALID:             "123456789012345678",
		CVN:               "123456789012345678",
		IUPR:              [18]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8},
		DTCCount:          2,
		DTCs:              []vehicle.DTC{{1, 2, 3, 4}, {2, 3, 4, 5}},
	}
	acm = vehicle.OBDBody{
		DiagnoseSupport:   0b1111110111110111,
		DiagnoseReadiness: 0b1111101111111011,
		CALID:             "123456789012345678",
		CVN:               "123456789012345678",
		IUPR:              [18]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8},
		DTCCount:          2,
		DTCs:              []vehicle.DTC{{1, 2, 3, 4}, {2, 3, 4, 5}},
	}
	return mcm, acm
}

// PackOBDData pack OBD packet to send to the server
func PackOBDData(mcm, acm vehicle.OBDBody) vehicle.Packet {
	merged := mergeOBDBody(mcm, acm)
	du := vehicle.DataUnit{
		UT:   vehicle.FetchDateTime(),
		SN:   123,
		Body: []vehicle.Body{&merged},
	}
	pkt := vehicle.Packet{
		Start:       "##",
		CommandUnit: 0x01,
		VIN:         "12345678901234567",
		CALID:       0x01,
		Encryption:  0x00,
		DULength:    du.LengthCal(),
		DU:          du,
		BCC:         0x56,
	}
	return pkt
}

func mergeOBDBody(mcm, acm vehicle.OBDBody) vehicle.OBDBody {
	merged := vehicle.OBDBody{
		VIN:      mcm.VIN,
		Protocal: mcm.Protocal,
		MIL:      mcm.MIL,
		CALID:    mcm.CALID,
		CVN:      mcm.CVN,
	}
	// all items need to be merged includes:
	// DiagnoseSupport
	merged.DiagnoseSupport = mergeDiagnoseSupport(mcm.DiagnoseSupport, acm.DiagnoseSupport)
	// DiagnoseReadiness
	merged.DiagnoseReadiness = mergeReadiness(mcm.DiagnoseSupport, mcm.DiagnoseReadiness, acm.DiagnoseSupport, acm.DiagnoseReadiness)
	// IUPR
	merged.IUPR = mergeIUPR(mcm.IUPR, acm.IUPR)
	// DTCs
	merged.DTCs = append(mcm.DTCs, acm.DTCs...)
	// DTCCount
	merged.DTCCount = uint8(len(merged.DTCs))
	return merged
}
