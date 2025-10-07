package crypto

import (
	"math/big"
	"testing"
)

func Test_NewPolynomialFromSecret(t *testing.T) {
	secret := big.NewInt(42)
	degree := 3
	p, err := NewPolynomialFromSecret(degree, *secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(p.coefficients) != degree {
		t.Errorf("expected degree %d, got %d", degree, len(p.coefficients))
	}
	if p.coefficients[0].Cmp(secret) != 0 {
		t.Errorf("expected secret %v, got %v", secret, p.coefficients[0])
	}

	pointAtOne := p.Evaluate(*big.NewInt(1))
	pointAtTwo := p.Evaluate(*big.NewInt(2))
	pointAtThree := p.Evaluate(*big.NewInt(3))
	pointAtFour := p.Evaluate(*big.NewInt(4))

	valueAtZero, err := ReconstructPolynomialAndEvaluateAtZero([]Point{
		pointAtOne,
		pointAtTwo,
		pointAtThree,
		pointAtFour,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if valueAtZero.Cmp(secret) != 0 {
		t.Errorf("expected secret %v, got %v", secret, valueAtZero)
	}
}
