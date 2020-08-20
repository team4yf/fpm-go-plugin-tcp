package plugin

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckSum(t *testing.T) {
	checksum := CheckSum([]byte{0x01, 0x02, 0x03, 0x04})
	fmt.Printf("check sum:%X \n", checksum)
	assert.Equal(t, (uint16)(0x2BA1), checksum, "should be 2ba1")

	data, _ := hex.DecodeString("0003000000AB1E0BA2E0AAA400260405B0000000700000000000011200000107000000110000001F0000001C000000260000000E")
	fmt.Printf("data:%X \n", data)
	checksum = CheckSum(data)
	fmt.Printf("check sum:%X \n", checksum)
	assert.Equal(t, (uint16)(0x2940), checksum, "should be 2940")

}
