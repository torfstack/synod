package crypto

import "errors"

var (
	MarkerBytes = []byte{0x64, 0x61, 0x74, 0x71}

	RsaOaepMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
	AesGcmMarkerBytes  = []byte{0x00, 0x00, 0x00, 0x02}

	RsaKeyLengthInBits  = 2048
	AesKeyLengthInBytes = 32

	KeyDerivationSalt       = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	KeyDerivationIterations = 600000

	ErrCryptoInvalidMarker   = errors.New("invalid encryption marker")
	ErrCryptoAlgorithmMarker = errors.New("invalid algorithm marker")
)
