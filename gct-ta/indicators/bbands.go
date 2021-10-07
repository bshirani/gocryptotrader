package indicators

// BBANDS returns the  Bollinger Bands for the given period
func BBANDS(inReal []float64, inTimePeriod int, inNbDevUp, inNbDevDn float64, inMAType MaType) (upper, middle, lower []float64) {
	outRealUpperBand := make([]float64, len(inReal))
	outRealMiddleBand := MA(inReal, inTimePeriod, inMAType)
	outRealLowerBand := make([]float64, len(inReal))

	tempBuffer2 := stdDev(inReal, inTimePeriod, 1.0)

	switch inNbDevUp {
	case inNbDevDn:
		if inNbDevUp == 1.0 {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i]
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = tempReal2 + tempReal
				outRealLowerBand[i] = tempReal2 - tempReal
			}
		} else {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i] * inNbDevUp
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = tempReal2 + tempReal
				outRealLowerBand[i] = tempReal2 - tempReal
			}
		}
	case 1.0:
		for i := 0; i < len(inReal); i++ {
			tempReal := tempBuffer2[i]
			tempReal2 := outRealMiddleBand[i]
			outRealUpperBand[i] = tempReal2 + tempReal
			outRealLowerBand[i] = tempReal2 - (tempReal * inNbDevDn)
		}
	default:
		if inNbDevDn == 1.0 {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i]
				tempReal2 := outRealMiddleBand[i]
				outRealLowerBand[i] = tempReal2 - tempReal
				outRealUpperBand[i] = tempReal2 + (tempReal * inNbDevUp)
			}
		} else {
			for i := 0; i < len(inReal); i++ {
				tempReal := tempBuffer2[i]
				tempReal2 := outRealMiddleBand[i]
				outRealUpperBand[i] = tempReal2 + (tempReal * inNbDevUp)
				outRealLowerBand[i] = tempReal2 - (tempReal * inNbDevDn)
			}
		}
	}

	return outRealUpperBand, outRealMiddleBand, outRealLowerBand
}
