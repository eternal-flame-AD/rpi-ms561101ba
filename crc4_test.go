package ms561101ba

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRC4(t *testing.T) {
	assert.Nil(t, verifyCRC4(PROMData{0x3132, 0x3334, 0x3536, 0x3738, 0x3940, 0x4142, 0x4344, 0x450B}))
}
