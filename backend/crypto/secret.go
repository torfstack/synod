package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// ShamirPrime is the Mersenne prime 2^127 - 1, used as the finite field modulus
// for Shamir secret sharing. All secret values and share coordinates must be
// strictly less than this value.
var ShamirPrime *big.Int

func init() {
	ShamirPrime = new(big.Int).Sub(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(127), nil),
		big.NewInt(1),
	)
}

// Polynomial represents a polynomial over GF(ShamirPrime).
// The constant term (index 0) is the secret; all other coefficients are random.
// A degree-k polynomial has exactly k+1 coefficients.
type Polynomial struct {
	coefficients []big.Int
}

// NewPolynomialFromSecret creates a degree-k polynomial with the given secret as
// the constant term and k random coefficients drawn from [1, ShamirPrime-1].
func NewPolynomialFromSecret(degree int, secret big.Int) (Polynomial, error) {
	coefficients := make([]big.Int, degree+1)
	coefficients[0] = secret
	for i := 1; i <= degree; i++ {
		b, err := RandomBigInt()
		if err != nil {
			return Polynomial{}, err
		}
		coefficients[i] = *b
	}
	return Polynomial{coefficients: coefficients}, nil
}

// ReconstructPolynomialAndEvaluateAtZero reconstructs the secret via Lagrange
// interpolation over GF(ShamirPrime). It requires exactly degree+1 distinct points.
func ReconstructPolynomialAndEvaluateAtZero(points []Point) (*big.Int, error) {
	p := ShamirPrime
	result := big.NewInt(0)
	n := len(points)

	for i := 0; i < n; i++ {
		numerator := new(big.Int).Set(&points[i].Output)
		denominator := big.NewInt(1)

		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			// numerator   *= (0 - x_j) mod p
			negXj := new(big.Int).Neg(&points[j].Input)
			negXj.Mod(negXj, p)
			numerator.Mul(numerator, negXj)
			numerator.Mod(numerator, p)

			// denominator *= (x_i - x_j) mod p
			diff := new(big.Int).Sub(&points[i].Input, &points[j].Input)
			diff.Mod(diff, p)
			denominator.Mul(denominator, diff)
			denominator.Mod(denominator, p)
		}

		// term = numerator * modInverse(denominator, p)
		inv := new(big.Int).ModInverse(denominator, p)
		if inv == nil {
			return nil, errors.New("zero denominator: share x-coordinates must be distinct and non-zero")
		}
		term := new(big.Int).Mul(numerator, inv)
		term.Mod(term, p)

		result.Add(result, term)
		result.Mod(result, p)
	}

	return result, nil
}

// RandomBigInt returns a random element from [1, ShamirPrime-1].
func RandomBigInt() (*big.Int, error) {
	return rand.Int(rand.Reader, ShamirPrime)
}

// Evaluate evaluates the polynomial at x over GF(ShamirPrime).
func (p Polynomial) Evaluate(x big.Int) Point {
	result := big.NewInt(0)
	for i, c := range p.coefficients {
		xPow := new(big.Int).Exp(&x, big.NewInt(int64(i)), ShamirPrime)
		term := new(big.Int).Mul(xPow, &c)
		term.Mod(term, ShamirPrime)
		result.Add(result, term)
		result.Mod(result, ShamirPrime)
	}
	return Point{Input: x, Output: *result}
}

type Point struct {
	Input  big.Int
	Output big.Int
}
