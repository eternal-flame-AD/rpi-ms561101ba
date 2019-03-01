package ms561101ba

import (
	"errors"
	"os"
	"time"

	"github.com/eternal-flame-AD/go-wiringpi/i2c"
)

// AddrLow is used when CSB is pulled low
const AddrLow uint8 = 0x76

// AddrHigh is used when CSB is pulled high
const AddrHigh uint8 = 0x77

// BaroSensor is the access handle to the sensor instance
type BaroSensor struct {
	fd int
}

// PROMData is the 8 uint16 data stored in the sensor for calibration
type PROMData [8]uint16

func osrToOffset(osr uint16) (offset uint8, success bool) {
	cur := uint16(256)
	for osr > cur {
		cur <<= 1
		offset += 2
	}
	return offset, cur == osr
}

// Open opens an I2C interface to the given address
func Open(addr uint8) (*BaroSensor, error) {
	fd, err := i2c.Open(addr)
	if err != nil {
		return nil, err
	}
	return &BaroSensor{fd}, nil
}

// Close closed the file descriptor to the sensor
func (c *BaroSensor) Close() error {
	f := os.NewFile(uintptr(c.fd), "baro")
	return f.Close()
}

// ReadPROM reads the sensor PROM data for calibration constants
func (c *BaroSensor) ReadPROM() (PROMData, error) {
	var res PROMData
	f := os.NewFile(uintptr(c.fd), "baro")
	for i := range res {
		/*
			d, err := i2c.ReadReg16(c.fd, 0xA0+2*i)
			if err != nil {
				return res, err
			}
			res[i] = d>>8 + d<<8
		*/

		if _, err := f.Write([]byte{byte(0xA0 + 2*i)}); err != nil {
			return res, err
		}
		var d [2]byte
		if _, err := f.Read(d[:]); err != nil {
			return res, err
		}
		res[i] = uint16(d[0])<<8 + uint16(d[1])
	}
	return res, verifyCRC4(res)
}

// Reset sends a reset signal to the sensor
func (c *BaroSensor) Reset() error {
	return i2c.Write(c.fd, 0x1E)
}

func (c *BaroSensor) readADC() (uint32, error) {
	f := os.NewFile(uintptr(c.fd), "baro")
	if _, err := f.Write([]byte{0}); err != nil {
		return 0, err
	}
	var d [3]byte
	if _, err := f.Read(d[:]); err != nil {
		return 0, err
	}
	return uint32(d[0])<<16 + uint32(d[1])<<8 + uint32(d[2]), nil
}

func (c *BaroSensor) doConv(base uint8, osr uint16) error {
	cmd := uint8(base)
	offset, success := osrToOffset(osr)
	if !success {
		return errors.New("osr is invalid")
	}
	cmd += offset
	if err := i2c.Write(c.fd, cmd); err != nil {
		return err
	}
	time.Sleep(10000 * time.Microsecond)
	return nil
}

// ReadPressureADC reads 24-bit digital pressure value
// OSR can be 256, 512, ..., 4096
func (c *BaroSensor) ReadPressureADC(osr uint16) (uint32, error) {
	if err := c.doConv(0x40, osr); err != nil {
		return 0, err
	}
	return c.readADC()
}

// ReadTemperatureADC reads 24-bit digital temperature value
// OSR can be 256, 512, ..., 4096
func (c *BaroSensor) ReadTemperatureADC(osr uint16) (uint32, error) {
	if err := c.doConv(0x50, osr); err != nil {
		return 0, err
	}
	return c.readADC()
}
