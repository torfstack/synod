package util

import "encoding/binary"

func IntToBytes(u uint32) []byte {
	return binary.LittleEndian.AppendUint32(nil, u)
}

func BytesToInt(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}
