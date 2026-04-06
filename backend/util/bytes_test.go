package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntToBytes_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input uint32
	}{
		{"zero", 0},
		{"one", 1},
		{"max_uint32", ^uint32(0)},
		{"arbitrary", 0xDEADBEEF},
		{"small", 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := IntToBytes(tt.input)
			require.Len(t, b, 4, "IntToBytes must always return 4 bytes")
			got := BytesToInt(b)
			require.Equal(t, tt.input, got)
		})
	}
}

func TestIntToBytes_LittleEndian(t *testing.T) {
	// Value 1 in little-endian is [0x01, 0x00, 0x00, 0x00]
	b := IntToBytes(1)
	require.Equal(t, []byte{0x01, 0x00, 0x00, 0x00}, b)
}

func TestBytesToInt_KnownValue(t *testing.T) {
	// [0x01, 0x00, 0x00, 0x00] is 1 in little-endian
	b := []byte{0x01, 0x00, 0x00, 0x00}
	require.Equal(t, uint32(1), BytesToInt(b))
}
