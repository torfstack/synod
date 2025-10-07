package crypto

import (
	"crypto/rand"
	"math/big"
)

var (
	MaxInt big.Int // 2^130 - 1
)

func init() {
	m := new(big.Int)
	MaxInt = *m.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(m, big.NewInt(1))
}

// Polynomial represents a polynomial with coefficients
// the constant term is the first element in the slice
// for our purposes: len(coefficients) == degree+1
type Polynomial struct {
	coefficients []big.Int
}

// NewPolynomialFromSecret creates a new polynomial with the given degree
// and secret as the constant term, all other coefficients are random
func NewPolynomialFromSecret(degree int, secret big.Int) (Polynomial, error) {
	coefficients := make([]big.Int, degree)
	coefficients[0] = secret
	for i := 1; i < degree; i++ {
		b, err := RandomBigInt()
		if err != nil {
			return Polynomial{}, err
		}
		coefficients[i] = *b
	}
	return Polynomial{coefficients: coefficients}, nil
}

// ReconstructPolynomialAndEvaluateAtZero takes a slice of points and returns the polynomial
// that passes through all of them, and evaluates it at 0
func ReconstructPolynomialAndEvaluateAtZero(points []Point) (*big.Int, error) {
	degree := len(points) - 1
	valueAtZero := new(big.Rat)
	for i := 0; i <= degree; i++ {
		term := new(big.Rat).SetFrac(&points[i].Output, big.NewInt(1))
		for j := 0; j <= degree; j++ {
			if i == j {
				continue
			}
			num := new(big.Rat).SetFrac(&points[j].Input, big.NewInt(-1))
			den := new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Sub(&points[i].Input, &points[j].Input))
			term.Mul(term, num)
			term.Mul(term, den)
		}
		valueAtZero.Add(valueAtZero, term)
	}

	return valueAtZero.Num(), nil
}

// RandomBigInt returns a random big.Int, a 130-bits integer, i.e 2^130 - 1
func RandomBigInt() (*big.Int, error) {
	return rand.Int(rand.Reader, &MaxInt)
}

// Evaluate returns the result of evaluating the polynomial at x
func (p Polynomial) Evaluate(x big.Int) Point {
	result := big.NewInt(0)
	for i, c := range p.coefficients {
		term := new(big.Int).Exp(&x, big.NewInt(int64(i)), nil)
		term.Mul(term, &c)
		result.Add(result, term)
	}
	return Point{Input: x, Output: *result}
}

type Point struct {
	Input  big.Int
	Output big.Int
}
