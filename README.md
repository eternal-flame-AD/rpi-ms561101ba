# rpi-ms561101ba

Reading [ms561101ba](https://www.te.com/usa-en/product-CAT-BLPS0036.html) sensors with wiringpi and golang on raspberry pi.


## Example

```golang
package main

import (
	"log"
	"time"

	ms561101ba "github.com/eternal-flame-AD/rpi-ms561101ba"
)

func main() {
	handle, err := ms561101ba.Open(0x77)
	if err != nil {
		panic(err)
	}
	handle.Reset()
	time.Sleep(100 * time.Millisecond)
	p, err := handle.ReadPROM()
	if err != nil {
		panic(err)
	}
	for {
		d1, err := handle.ReadPressureADC(4096)
		if err != nil {
			log.Printf("error while reading pressure ADC value: %s", err.Error())
		}
		d2, err := handle.ReadTemperatureADC(4096)
		if err != nil {
			log.Printf("error while reading temperature ADC value: %s", err.Error())
		}
		t, p := ms561101ba.CalcTemp(d2, p), ms561101ba.CalcPressure(d1, d2, p)
		log.Printf("T=%.2f degC P=%.2f hPa", t, p)
	}
}
```