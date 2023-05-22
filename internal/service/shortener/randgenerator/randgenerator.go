package randgenerator

import (
	"math/rand"
	"time"
)

// RandGenerator implements Generator interface.
// It uses random to generate path.
type RandGenerator struct {
	rand     *rand.Rand
	alphabet []byte
	shortLen int
}

func New(alphabet []byte, shortLen int) *RandGenerator {
	return &RandGenerator{
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		alphabet: alphabet,
		shortLen: shortLen,
	}
}

func (p *RandGenerator) Generate(_ string) string {
	b := make([]byte, p.shortLen)
	for i := range b {
		b[i] = p.alphabet[p.rand.Intn(len(p.alphabet))]
	}
	return string(b)
}
