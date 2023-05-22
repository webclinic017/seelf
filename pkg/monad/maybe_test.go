package monad_test

import (
	"testing"
	"time"

	"github.com/YuukanOO/seelf/pkg/monad"
	"github.com/YuukanOO/seelf/pkg/testutil"
)

func Test_Maybe(t *testing.T) {
	t.Run("should have a default state without value", func(t *testing.T) {
		var m monad.Maybe[time.Time]

		testutil.IsFalse(t, m.HasValue())
	})

	t.Run("could be created empty", func(t *testing.T) {
		m := monad.None[time.Time]()
		testutil.IsFalse(t, m.HasValue())
	})

	t.Run("could be created with a defined value", func(t *testing.T) {
		m := monad.Value("ok")

		testutil.Equals(t, "ok", m.MustGet())
		testutil.IsTrue(t, m.HasValue())
	})

	t.Run("could be assigned a value", func(t *testing.T) {
		var (
			m   monad.Maybe[time.Time]
			now = time.Now().UTC()
		)

		m = m.WithValue(now)
		testutil.Equals(t, now, m.MustGet())
		testutil.IsTrue(t, m.HasValue())
	})

	t.Run("could unset its value", func(t *testing.T) {
		m := monad.Value(time.Now().UTC())

		m = m.None()

		testutil.IsFalse(t, m.HasValue())
	})

	t.Run("should panic if trying to access a value with MustGet", func(t *testing.T) {
		defer func() {
			err := recover()
			testutil.IsNotNil(t, err)
			testutil.Equals(t, "trying to access a monad's value but none is set", err.(string))
		}()

		var m monad.Maybe[time.Time]

		m.MustGet()
	})

	t.Run("could returns its value if its set", func(t *testing.T) {
		var now = time.Now().UTC()

		m := monad.Value(now)

		testutil.Equals(t, now, m.MustGet())
	})

	t.Run("could returns its value or fallback if not set", func(t *testing.T) {
		var (
			woValue monad.Maybe[string]
			wValue  = monad.Value("got a value")
		)

		testutil.Equals(t, "got a value", wValue.Get("default"))
		testutil.Equals(t, "default", woValue.Get("default"))
	})

	t.Run("should implements the valuer interface", func(t *testing.T) {
		var m monad.Maybe[time.Time]

		driverValue, err := m.Value()

		testutil.IsNil(t, err)
		testutil.IsNil(t, driverValue)

		now := time.Now().UTC()
		m = m.WithValue(now)
		driverValue, err = m.Value()

		testutil.IsNil(t, err)
		testutil.IsTrue(t, driverValue == now)
	})
}
