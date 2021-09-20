package benchmarks

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stellar/go/gxdr"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/require"
	goxdr "github.com/xdrpp/goxdr/xdr"
)

const input64 = "AAAAAgAAAADy2f6v1nv9lXdvl5iZvWKywlPQYsZ1JGmmAfewflnbUAAABLACG4bdAADOYQAAAAEAAAAAAAAAAAAAAABhSLZ9AAAAAAAAAAEAAAABAAAAAF8wDgs7+R5R2uftMvvhHliZOyhZOQWsWr18/Fu6S+g0AAAAAwAAAAJHRE9HRQAAAAAAAAAAAAAAUwsPRQlK+jECWsJLURlsP0qsbA/aIaB/z50U79VSRYsAAAAAAAAAAAAAAYMAAA5xAvrwgAAAAAAAAAAAAAAAAAAAAAJ+WdtQAAAAQCTonAxUHyuVsmaSeGYuVsGRXgxs+wXvKgSa+dapZWN4U9sxGPuApjiv/UWb47SwuFQ+q40bfkPYT1Tff4RfLQe6S+g0AAAAQBlFjwF/wpGr+DWbjCyuolgM1VP/e4ubfUlVnDAdFjJUIIzVakZcr5omRSnr7ClrwEoPj49h+vcLusagC4xFJgg="

var input = func() []byte {
	input, err := base64.StdEncoding.DecodeString(input64)
	if err != nil {
		panic(err)
	}
	return input
}()

func BenchmarkXDRUnmarshal(b *testing.B) {
	te := xdr.TransactionEnvelope{}

	// Make sure the input is valid.
	err := te.UnmarshalBinary(input)
	require.NoError(b, err)

	// Benchmark.
	for i := 0; i < b.N; i++ {
		_ = te.UnmarshalBinary(input)
	}
}

func BenchmarkGXDRUnmarshal(b *testing.B) {
	te := gxdr.TransactionEnvelope{}

	// Make sure the input is valid, note goxdr will panic if there's a
	// marshaling error.
	te.XdrMarshal(&goxdr.XdrIn{In: bytes.NewReader(input)}, "")

	// Benchmark.
	r := bytes.NewReader(input)
	for i := 0; i < b.N; i++ {
		r.Reset(input)
		te.XdrMarshal(&goxdr.XdrIn{In: r}, "")
	}
}

func BenchmarkXDRMarshal(b *testing.B) {
	te := xdr.TransactionEnvelope{}

	// Make sure the input is valid.
	err := te.UnmarshalBinary(input)
	require.NoError(b, err)
	output, err := te.MarshalBinary()
	require.NoError(b, err)
	require.Equal(b, input, output)

	// Benchmark.
	for i := 0; i < b.N; i++ {
		_, _ = te.MarshalBinary()
	}
}

func BenchmarkGXDRMarshal(b *testing.B) {
	te := gxdr.TransactionEnvelope{}

	// Make sure the input is valid, note goxdr will panic if there's a
	// marshaling error.
	te.XdrMarshal(&goxdr.XdrIn{In: bytes.NewReader(input)}, "")
	output := bytes.Buffer{}
	te.XdrMarshal(&goxdr.XdrOut{Out: &output}, "")

	// Benchmark.
	for i := 0; i < b.N; i++ {
		output.Reset()
		te.XdrMarshal(&goxdr.XdrOut{Out: &output}, "")
	}
}
