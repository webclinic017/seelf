package infra_test

import (
	"testing"

	"github.com/YuukanOO/seelf/internal/auth/infra"
	"github.com/YuukanOO/seelf/pkg/testutil"
)

func Test_KeyGenerator(t *testing.T) {
	t.Run("should generate an API key", func(t *testing.T) {
		generator := infra.NewKeyGenerator()
		key, err := generator.Generate()
		testutil.IsNil(t, err)
		testutil.HasNChars(t, 64, key)
	})
}
