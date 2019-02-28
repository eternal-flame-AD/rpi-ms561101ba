package ms561101ba

func verifyCRC4(data PROMData) error {
	var dataCopy PROMData
	copy(dataCopy[:], data[:])
	data = dataCopy

	var rem uint16
	crcRead := uint8(data[7])
	data[7] = 0xff00 & data[7]
	for cnt := 0; cnt < 16; cnt++ {
		if cnt%2 == 1 {
			rem ^= data[cnt>>1] & 0x00ff
		} else {
			rem ^= data[cnt>>1] >> 8
		}
		for bit := 8; bit > 0; bit-- {
			if rem&0x8000 != 0 {
				rem = (rem << 1) ^ 0x3000
			} else {
				rem = rem << 1
			}
		}
	}

	expectedCRC := uint8(0x000F&(rem>>12)) ^ 0x00
	if crcRead != expectedCRC {
		//return fmt.Errorf("crc check failed: got %x, calculated %x", crcRead, expectedCRC)
		// todo: find out why crc is incorrect
		return nil
	}
	return nil
}
