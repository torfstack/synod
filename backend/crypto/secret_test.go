package crypto

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PolynomialSecretSharing(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "can reconstruct polynomial from points",
			run: func(t *testing.T) {
				secret := big.NewInt(42)
				degree := 3
				p, err := NewPolynomialFromSecret(degree, *secret)
				assert.NoError(t, err)
				assert.Len(t, p.coefficients, degree)
				assert.Equal(t, *secret, p.coefficients[0])

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
				assert.NoError(t, err)
				assert.Equal(t, secret, valueAtZero)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t)
		})
	}
}
