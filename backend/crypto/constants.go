package crypto

import "errors"

const (
	AesKeyLengthInBytes     = 32
	KeyDerivationIterations = 600000
	KDFSaltLength           = 16
)

var (
	MarkerBytes = []byte{0x64, 0x61, 0x74, 0x71}

	AsymmetricMarkerBytes = []byte{0x00, 0x00, 0x00, 0x01}
	AesGcmMarkerBytes     = []byte{0x00, 0x00, 0x00, 0x02}

	// KDFSaltPrefix is prepended to password-protected key material.
	KDFSaltPrefix = []byte{0x73, 0x79, 0x6e, 0x64}

	ErrCryptoInvalidMarker     = errors.New("invalid encryption marker")
	ErrCryptoAlgorithmMarker   = errors.New("invalid algorithm marker")
	ErrCryptoInvalidCiphertext = errors.New("invalid ciphertext")
)
