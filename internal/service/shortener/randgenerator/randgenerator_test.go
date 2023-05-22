package randgenerator

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandGenerator(t *testing.T) {
	length := 10
	t.Run("uniques", func(t *testing.T) {
		generated := make(map[string]struct{})

		generator := RandGenerator{
			alphabet: []byte("abcdefghijklmnopqrstuvwxyz"),
			shortLen: length,
			rand:     rand.New(rand.NewSource(42)),
		}

		N := 10000
		for i := 0; i < N; i++ {
			res := generator.Generate("")
			require.Equal(t, length, len(res))
			if _, ok := generated[res]; ok {
				t.Errorf("generated %s twice", res)
			}
			generated[res] = struct{}{}
		}
	})
}
