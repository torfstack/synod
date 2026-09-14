package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitAndCombineThresholdSecret(t *testing.T) {
	secret := []byte("a 32-byte threshold encryption key")
	shares, err := SplitSecret(secret, 3, 5)
	require.NoError(t, err)
	require.Len(t, shares, 5)

	recovered, err := CombineShares([][]byte{shares[0], shares[2], shares[4]})
	require.NoError(t, err)
	require.Equal(t, secret, recovered)
}

func TestCombineSharesNeedsEnoughDistinctShares(t *testing.T) {
	shares, err := SplitSecret([]byte("secret"), 2, 2)
	require.NoError(t, err)

	_, err = CombineShares([][]byte{shares[0], shares[0]})
	require.Error(t, err)
}

func TestSplitSecretValidatesThreshold(t *testing.T) {
	_, err := SplitSecret([]byte("secret"), 1, 2)
	require.Error(t, err)
	_, err = SplitSecret([]byte("secret"), 3, 2)
	require.Error(t, err)
}
