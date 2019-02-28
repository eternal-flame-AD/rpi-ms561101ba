package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	ms561101ba "github.com/eternal-flame-AD/rpi-ms561101ba"
)

const DeviceAddr = 0x77

var prefix string
var metricAddr string

type Measurement struct {
	Prefix      string
	Timestamp   int64
	Pressure    float64
	Temperature float64
}

var metrics Measurement
var metricsLock = sync.RWMutex{}
var metricTpl = template.Must(template.New("metric").Parse(strings.TrimSpace(`
# HELP {{.Prefix}}_pressure pressure in hPa
# TYPE {{.Prefix}}_pressure gauge
# HELP {{.Prefix}}_temperature temperature in degC
# TYPE {{.Prefix}}_temperature gauge
{{.Prefix}}_pressure{device="ms561101ba"} {{printf "%.2f" .Pressure}} {{.Timestamp}}
{{.Prefix}}_temperature{device="ms561101ba"} {{printf "%.2f" .Temperature}} {{.Timestamp}}
`)))

func init() {
	p := flag.String("p", "ms56xx", "metric prefix")
	m := flag.String("m", ":9100", "metrics address")
	flag.Parse()

	prefix = *p
	metricAddr = *m
}

func takeReading(device *ms561101ba.BaroSensor, prom ms561101ba.PROMData) {
	d1, err := device.ReadPressureADC(4096)
	if err != nil {
		log.Printf("error while reading pressure ADC value: %s", err.Error())
		return
	}
	d2, err := device.ReadTemperatureADC(4096)
	if err != nil {
		log.Printf("error while reading temperature ADC value: %s", err.Error())
		return
	}
	t, p := ms561101ba.CalcTemp(d2, prom), ms561101ba.CalcPressure(d1, d2, prom)
	log.Printf("T=%.2f degC P=%.2f hPa", t, p)
	now := time.Now()
	metricsLock.Lock()
	defer metricsLock.Unlock()
	metrics = Measurement{
		Temperature: t,
		Pressure:    p,
		Timestamp:   now.Unix() * 1000,
		Prefix:      prefix,
	}
}

func main() {
	handle, err := ms561101ba.Open(0x77)
	if err != nil {
		panic(err)
	}
	handle.Reset()
	time.Sleep(100 * time.Millisecond)
	prom, err := handle.ReadPROM()
	if err != nil {
		panic(err)
	}

	takeReading(handle, prom)
	go func() {
		for range time.NewTicker(1 * time.Second).C {
			takeReading(handle, prom)
		}
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain; version=0.0.4")
		metricsLock.RLock()
		defer metricsLock.RUnlock()
		metricTpl.Execute(w, metrics)
		w.Write([]byte{'\n'})
	})

	log.Println("promethus metrics running at " + metricAddr)
	log.Fatal(http.ListenAndServe(metricAddr, nil))

}
