//go:build regressiontests
// +build regressiontests

package refactortests

import (
	"bytes"
	"strconv"
	"testing"

	fuzz "github.com/google/gofuzz"
	keypairmaster "github.com/stellar/go-master/keypair"
	keypair "github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inputSize = 32

func TestRefactorOfKeypairIsConsistent(t *testing.T) {
	f := fuzz.New()
	for i := 0; i < 10000; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// Generate a valid address and seed to use as inputs to tests.
			rawSeed := [32]byte{}
			f.Fuzz(&rawSeed)
			sk, err := keypairmaster.FromRawSeed(rawSeed)
			require.NoError(t, err)
			address := sk.Address()
			seed := sk.Seed()

			t.Logf("address = %s", address)
			t.Logf("seed = %s", seed)

			// Test consistency of FromAddresses.
			{
				want := keypairmaster.MustParseAddress(address)
				got := keypair.MustParseAddress(address)
				t.Run("EqualFromAddress(want, got)", func(t *testing.T) {
					t.Parallel()
					for i := 0; i < 10; i++ {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							EqualFromAddress(t, sk, want, got)
						})
					}
				})
				t.Run("EqualFromAddress(want.FromAddress(), got.FromAddress())", func(t *testing.T) {
					t.Parallel()
					for i := 0; i < 10; i++ {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							EqualFromAddress(t, sk, want.FromAddress(), got.FromAddress())
						})
					}
				})
			}

			// Test consistency of Fulls.
			{
				want := keypairmaster.MustParseFull(seed)
				got := keypair.MustParseFull(seed)
				t.Run("EqualFull(want, got)", func(t *testing.T) {
					t.Parallel()
					for i := 0; i < 10; i++ {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							EqualFull(t, want, got)
						})
					}
				})
				t.Run("EqualFromAddress(want.FromAddress(), got.FromAddress())", func(t *testing.T) {
					t.Parallel()
					for i := 0; i < 10; i++ {
						t.Run(strconv.Itoa(i), func(t *testing.T) {
							EqualFromAddress(t, sk, want.FromAddress(), got.FromAddress())
						})
					}
				})
			}
		})
	}
}

func EqualFromAddress(t *testing.T, wantFull *keypairmaster.Full, want *keypairmaster.FromAddress, got *keypair.FromAddress) {
	t.Helper()

	// Check basic functions for consistency.
	assert.Equal(t, want.Address(), got.Address())
	assert.Equal(t, want.Hint(), got.Hint())

	f := fuzz.New()
	input := []byte{}
	f.NumElements(inputSize, inputSize*10).Fuzz(&input)
	t.Log("input:", input)
	input2 := []byte{}
	for {
		f.NumElements(inputSize, inputSize*10).Fuzz(&input2)
		if !bytes.Equal(input, input2) {
			break
		}
	}
	t.Log("input2:", input2)

	// Getting signature to verify.
	sig, err := wantFull.Sign(input)
	require.NoError(t, err)

	// Check that verifying is consistent.
	wantErr := want.Verify(input, sig)
	gotErr := got.Verify(input, sig)
	assert.Equal(t, wantErr, gotErr)
	assert.NoError(t, wantErr)
	assert.NoError(t, gotErr)

	// Check that verification failure is consistent.
	wantErr = want.Verify(input2, sig)
	gotErr = got.Verify(input2, sig)
	assert.Equal(t, wantErr, gotErr)
	assert.Error(t, wantErr)
	assert.Error(t, gotErr)
}

func EqualFull(t *testing.T, want *keypairmaster.Full, got *keypair.Full) {
	t.Helper()

	// Check basic functions for consistency.
	assert.Equal(t, want.Address(), got.Address())
	assert.Equal(t, want.Seed(), got.Seed())
	assert.Equal(t, want.Hint(), got.Hint())

	f := fuzz.New()
	input := []byte{}
	f.NumElements(inputSize, inputSize*10).Fuzz(&input)
	t.Log("input:", input)
	input2 := []byte{}
	for {
		f.NumElements(inputSize, inputSize*10).Fuzz(&input2)
		if !bytes.Equal(input, input2) {
			break
		}
	}
	t.Log("input2:", input2)

	// Check signing is consistent.
	wantSig, err := want.Sign(input)
	require.NoError(t, err)
	gotSig, err := got.Sign(input)
	require.NoError(t, err)
	assert.Equal(t, wantSig, gotSig)

	wantSigBase64, err := want.SignBase64(input)
	require.NoError(t, err)
	gotSigBase64, err := got.SignBase64(input)
	require.NoError(t, err)
	assert.Equal(t, wantSigBase64, gotSigBase64)

	wantSigDec, err := want.SignBase64(input)
	require.NoError(t, err)
	gotSigDec, err := got.SignBase64(input)
	require.NoError(t, err)
	assert.Equal(t, wantSigDec, gotSigDec)

	// Check that verifying is consistent.
	wantErr := want.Verify(input, wantSig)
	gotErr := got.Verify(input, gotSig)
	assert.Equal(t, wantErr, gotErr)
	assert.NoError(t, wantErr)
	assert.NoError(t, gotErr)

	// Check that verifying is consistent with signatures from other package.
	wantErr = want.Verify(input, gotSig)
	gotErr = got.Verify(input, wantSig)
	assert.Equal(t, wantErr, gotErr)
	assert.NoError(t, wantErr)
	assert.NoError(t, gotErr)

	// Check that verification failure is consistent.
	wantErr = want.Verify(input2, wantSig)
	gotErr = got.Verify(input2, gotSig)
	assert.Equal(t, wantErr, gotErr)
	assert.Error(t, wantErr)
	assert.Error(t, gotErr)
}
