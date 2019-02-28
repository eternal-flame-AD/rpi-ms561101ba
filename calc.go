package ms561101ba

func calcDt(d2 uint32, c PROMData) int64 {
	return int64(d2) - int64(c[5])*(1<<8)
}

// CalcTemp takes the digital temperature value and the PROM data and calculates a temperature value.
// Precision is .2f
func CalcTemp(d2 uint32, c PROMData) float64 {
	dt := calcDt(d2, c)
	return (2000. + float64(dt)*float64(c[6])/(1<<23)) / 100
}

// CalcPressure takes the digital pressure and temperature value and the PROM data and calculates a temperature-compensated pressure value.
// Precision is .2f
func CalcPressure(d1 uint32, d2 uint32, c PROMData) float64 {
	dt := calcDt(d2, c)
	off := int64(c[2])*(1<<16) + (int64(c[4])*int64(dt))>>7
	sens := int64(c[1])*(1<<15) + (int64(c[3])*int64(dt))>>8
	p := (int64(d1)*sens>>21 - off) >> 15
	return float64(p) / 100
}
