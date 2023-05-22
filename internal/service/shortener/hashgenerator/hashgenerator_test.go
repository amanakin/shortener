package hashgenerator

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashGenerator(t *testing.T) {
	length := 10
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	generator := New([]byte("abcdefghijklmnopqrstuvwxyz"), length)

	t.Run("same result", func(t *testing.T) {
		res1 := generator.Generate("https://golang.org")
		res2 := generator.Generate("https://golang.org")

		require.Equal(t, res1, res2)
		require.Equal(t, length, len(res1))
	})

	t.Run("uniques", func(t *testing.T) {
		r := rand.New(rand.NewSource(42))

		generated := make(map[string]struct{})

		N := 10000
		for i := 0; i < N; i++ {
			input := make([]byte, length)
			for j := 0; j < length; j++ {
				input[j] = alphabet[r.Intn(len(alphabet))]
			}

			res := generator.Generate(string(input))
			if _, ok := generated[res]; ok {
				t.Errorf("generated %s twice", res)
			}
			generated[res] = struct{}{}
		}
	})
}
