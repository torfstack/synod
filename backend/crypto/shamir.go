package crypto

import (
	"crypto/rand"
	"errors"
)

func SplitSecret(secret []byte, threshold, count int) ([][]byte, error) {
	if len(secret) == 0 || threshold < 2 || threshold > count || count > 255 {
		return nil, errors.New("invalid secret sharing parameters")
	}
	shares := make([][]byte, count)
	for i := range shares {
		shares[i] = make([]byte, len(secret)+1)
		shares[i][0] = byte(i + 1)
	}
	coefficients := make([]byte, threshold-1)
	for offset, value := range secret {
		if _, err := rand.Read(coefficients); err != nil {
			return nil, err
		}
		for i := range shares {
			x := shares[i][0]
			y := value
			power := x
			for _, coefficient := range coefficients {
				y ^= gfMultiply(coefficient, power)
				power = gfMultiply(power, x)
			}
			shares[i][offset+1] = y
		}
	}
	return shares, nil
}

func CombineShares(shares [][]byte) ([]byte, error) {
	if len(shares) < 2 || len(shares[0]) < 2 {
		return nil, errors.New("not enough shares")
	}
	length := len(shares[0])
	seen := make(map[byte]struct{}, len(shares))
	for _, share := range shares {
		if len(share) != length || share[0] == 0 {
			return nil, errors.New("invalid share")
		}
		if _, exists := seen[share[0]]; exists {
			return nil, errors.New("duplicate share")
		}
		seen[share[0]] = struct{}{}
	}
	secret := make([]byte, length-1)
	for offset := range secret {
		for i, share := range shares {
			xi := share[0]
			basis := byte(1)
			for j, other := range shares {
				if i == j {
					continue
				}
				basis = gfMultiply(basis, gfDivide(other[0], other[0]^xi))
			}
			secret[offset] ^= gfMultiply(share[offset+1], basis)
		}
	}
	return secret, nil
}

func gfMultiply(a, b byte) byte {
	var result byte
	for b != 0 {
		if b&1 != 0 {
			result ^= a
		}
		high := a & 0x80
		a <<= 1
		if high != 0 {
			a ^= 0x1b
		}
		b >>= 1
	}
	return result
}

func gfDivide(a, b byte) byte {
	if b == 0 {
		return 0
	}
	return gfMultiply(a, gfPow(b, 254))
}

func gfPow(value byte, exponent int) byte {
	result := byte(1)
	for exponent > 0 {
		if exponent&1 != 0 {
			result = gfMultiply(result, value)
		}
		value = gfMultiply(value, value)
		exponent >>= 1
	}
	return result
}
