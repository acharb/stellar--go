package keypair

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

type Full struct {
	seed string

	// publicKey is the ed25519 public key derived from the seed. It must be set
	// during construction of the value.
	publicKey ed25519.PublicKey

	// privateKey is the ed25519 private key derived from the seed. It must be
	// set during construction of the value.
	privateKey ed25519.PrivateKey
}

func newFull(seed string) (*Full, error) {
	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, seed)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawSeed)
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	return &Full{
		seed:       seed,
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func newFullFromRawSeed(rawSeed [32]byte) (*Full, error) {
	seed, err := strkey.Encode(strkey.VersionByteSeed, rawSeed[:])
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(rawSeed[:])
	pub, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		panic(err)
	}
	return &Full{
		seed:       seed,
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func (kp *Full) Address() string {
	return strkey.MustEncode(strkey.VersionByteAccountID, kp.publicKey[:])
}

// FromAddress gets the address-only representation, or public key, of this
// Full keypair.
func (kp *Full) FromAddress() *FromAddress {
	return newFromAddressWithPublicKey(kp.Address(), kp.publicKey)
}

func (kp *Full) Hint() (r [4]byte) {
	copy(r[:], kp.publicKey[28:])
	return
}

func (kp *Full) Seed() string {
	return kp.seed
}

func (kp *Full) Verify(input []byte, sig []byte) error {
	if len(sig) != 64 {
		return ErrInvalidSignature
	}
	if !ed25519.Verify(kp.publicKey, input, sig) {
		return ErrInvalidSignature
	}
	return nil
}

func (kp *Full) Sign(input []byte) ([]byte, error) {
	return ed25519.Sign(kp.privateKey, input), nil
}

// SignBase64 signs the input data and returns a base64 encoded string, the
// common format in which signatures are exchanged.
func (kp *Full) SignBase64(input []byte) (string, error) {
	sig, err := kp.Sign(input)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func (kp *Full) SignDecorated(input []byte) (xdr.DecoratedSignature, error) {
	sig, err := kp.Sign(input)
	if err != nil {
		return xdr.DecoratedSignature{}, err
	}

	return xdr.DecoratedSignature{
		Hint:      xdr.SignatureHint(kp.Hint()),
		Signature: xdr.Signature(sig),
	}, nil
}

func (kp *Full) Equal(f *Full) bool {
	if kp == nil && f == nil {
		return true
	}
	if kp == nil || f == nil {
		return false
	}
	return kp.seed == f.seed
}
