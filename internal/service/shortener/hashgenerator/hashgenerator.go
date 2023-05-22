package hashgenerator

import (
	"crypto/sha256"
	"hash"
	"math/big"
	"strings"
)

// HashGenerator implements more secure Generator interface.
type HashGenerator struct {
	hash     hash.Hash
	alphabet []byte
	shortLen int
}

func New(alphabet []byte, shortLen int) *HashGenerator {
	return &HashGenerator{
		hash:     sha256.New(),
		alphabet: alphabet,
		shortLen: shortLen,
	}
}

func (g *HashGenerator) Generate(input string) string {
	g.hash.Reset()
	g.hash.Write([]byte(input))
	hashed := g.hash.Sum(nil)

	base := big.NewInt(int64(len(g.alphabet)))
	num := big.NewInt(0).SetBytes(hashed)
	index := big.NewInt(0)

	var sb strings.Builder
	for len(num.Bits()) > 0 && sb.Len() < g.shortLen {
		index.Mod(num, base)
		num.Div(num, base)

		sb.WriteByte(g.alphabet[index.Int64()])
	}

	for sb.Len() < g.shortLen {
		sb.WriteByte(g.alphabet[0])
	}

	return sb.String()
}
