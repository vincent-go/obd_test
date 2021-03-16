package mock

func mergeIUPR(mcmIUPR, acmIUPR [18]uint16) [18]uint16 {
	iupr := [18]uint16{}
	// for "OBD monitoring conditions encountered counts" and "ignation cycle counter", just consider MCM counters
	iupr[0], iupr[1] = mcmIUPR[0], mcmIUPR[1]
	// for the IUPRs that need to consider MCM and ACM, apply selection logic
	for i := 1; i < 8; i++ {
		m, a := [2]uint16{}, [2]uint16{}
		m[0], m[1] = mcmIUPR[2*i], mcmIUPR[2*i+1]
		a[0], a[1] = acmIUPR[2*i], acmIUPR[2*i+1]

		iuprGroup := selectionLogic(m, a)

		iupr[2*i], iupr[2*i+1] = iuprGroup[0], iuprGroup[1]
	}
	// for "Fuel Monitor Completion Condition Counts" and "Fuel Monitor ConditionsEncountered Counts", set to 0, as it is not supported
	iupr[16], iupr[17] = 0, 0
	return iupr
}

// selectionLogic is the selection logic for iupr items that needs to compare between MCM and ACM counters
func selectionLogic(mc, ac [2]uint16) [2]uint16 {
	if mc[0]+mc[1]+ac[0]+ac[1] == 0 {
		return [2]uint16{0, 0}
	}
	if mc[0]+mc[1] == 0 {
		return ac
	}
	if ac[0]+ac[1] == 0 {
		return mc
	}
	if (mc[0]*ac[0] > 0) && (mc[1]+ac[1] == 0) {
		if mc[0] > ac[0] {
			return ac
		}
		return mc
	}
	if (mc[0] > 0) && (mc[1] == 0) {
		return ac
	}
	if (ac[0] > 0) && (ac[1] == 0) {
		return mc
	}
	if mc[1]*ac[1] > 0 {
		if float64(mc[0])/float64(mc[1]) > float64(ac[0])/float64(ac[1]) {
			return ac
		}
		return mc
	}
	return mc
}
