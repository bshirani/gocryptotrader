package indicators

// ATR returns the Average True Range for the given period
func ATR(inHigh, inLow, inClose []float64, inTimePeriod int) []float64 {
	outReal := make([]float64, len(inClose))

	if inTimePeriod < 1 || inTimePeriod > len(inHigh) || inTimePeriod > len(inLow) || inTimePeriod > len(inClose) {
		return outReal
	}

	if inTimePeriod == 1 {
		return trueRange(inHigh, inLow, inClose)
	}

	today := inTimePeriod + 1
	tr := trueRange(inHigh, inLow, inClose)
	prevATRTemp := SMA(tr, inTimePeriod)
	prevATR := prevATRTemp[inTimePeriod]
	outReal[inTimePeriod] = prevATR

	for outIdx := inTimePeriod + 1; outIdx < len(inClose); outIdx++ {
		prevATR *= float64(inTimePeriod) - 1.0
		prevATR += tr[today]
		prevATR /= float64(inTimePeriod)
		outReal[outIdx] = prevATR
		today++
	}

	return outReal
}
