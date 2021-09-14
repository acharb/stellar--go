package keypair

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"sync"

	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

type Full struct {
	seed string

	// cacheOnce synchronizes the first call to keys() and ensures concurrent
	// calls to any function that calls keys() do not read or write the cached
	// fields while they are being written for the first time.
	cacheOnce sync.Once

	// cachedPublicKey is a cached copy of the ed25519 public key after first
	// call to keys(). Code should never access this field, call keys() instead.
	cachedPublicKey ed25519.PublicKey

	// cachedPrivateKey is a cached copy of the ed25519 private key after first
	// call to keys(). Code should never access this field, call keys() instead.
	cachedPrivateKey ed25519.PrivateKey
}

func (kp *Full) Address() string {
	return strkey.MustEncode(strkey.VersionByteAccountID, kp.publicKey()[:])
}

// FromAddress gets the address-only representation, or public key, of this
// Full keypair.
func (kp *Full) FromAddress() *FromAddress {
	return &FromAddress{address: kp.Address()}
}

func (kp *Full) Hint() (r [4]byte) {
	copy(r[:], kp.publicKey()[28:])
	return
}

func (kp *Full) Seed() string {
	return kp.seed
}

func (kp *Full) Verify(input []byte, sig []byte) error {
	if len(sig) != 64 {
		return ErrInvalidSignature
	}
	if !ed25519.Verify(kp.publicKey(), input, sig) {
		return ErrInvalidSignature
	}
	return nil
}

func (kp *Full) Sign(input []byte) ([]byte, error) {
	_, priv := kp.keys()
	return ed25519.Sign(priv, input), nil
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

func (kp *Full) publicKey() ed25519.PublicKey {
	pub, _ := kp.keys()
	return pub
}

func (kp *Full) keys() (ed25519.PublicKey, ed25519.PrivateKey) {
	kp.cacheOnce.Do(func() {
		reader := bytes.NewReader(kp.rawSeed())
		pub, priv, err := ed25519.GenerateKey(reader)
		if err != nil {
			panic(err)
		}
		kp.cachedPublicKey = pub
		kp.cachedPrivateKey = priv
	})
	return kp.cachedPublicKey, kp.cachedPrivateKey
}

func (kp *Full) rawSeed() []byte {
	return strkey.MustDecode(strkey.VersionByteSeed, kp.seed)
}
