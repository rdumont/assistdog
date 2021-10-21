package defaults

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCompareString(t *testing.T) {
	t.Run("returns nil for equal strings", func(t *testing.T) {
		err := CompareString("abc", "abc")

		require.NoError(t, err)
	})

	t.Run("returns error for actual that isn't a string", func(t *testing.T) {
		err := CompareString("abc", 123)

		require.EqualError(t, err, "123 is not a string")
	})

	t.Run("returns error for different strings", func(t *testing.T) {
		err := CompareString("abc", "ABC")

		require.EqualError(t, err, "expected abc, but got ABC")
	})
}

func TestCompareInt(t *testing.T) {
	t.Run("returns nil for equal integers", func(t *testing.T) {
		err := CompareInt("123", 123)

		require.NoError(t, err)
	})

	t.Run("returns error for actual that isn't an integer", func(t *testing.T) {
		err := CompareInt("123", "not an integer")

		require.EqualError(t, err, "not an integer is not an int")
	})

	t.Run("returns error for different integers", func(t *testing.T) {
		err := CompareInt("123", 456)

		require.EqualError(t, err, "expected 123, but got 456")
	})
}

func TestCompareInt32(t *testing.T) {
	t.Run("returns nil for equal int32 values", func(t *testing.T) {
		err := CompareInt32("123", int32(123))

		require.NoError(t, err)
	})

	t.Run("returns error for actual that isn't an int32", func(t *testing.T) {
		err := CompareInt32("123", 123)

		require.EqualError(t, err, "123 is not an int32")
	})

	t.Run("returns error for different int32 values", func(t *testing.T) {
		err := CompareInt32("123", int32(456))

		require.EqualError(t, err, "expected 123, but got 456")
	})
}

func TestCompareTime(t *testing.T) {
	validTime, err := time.Parse(time.RFC3339, "2020-11-05T16:01:54Z")
	require.NoError(t, err)

	t.Run("returns nil for equal times", func(t *testing.T) {
		err := CompareTime("2020-11-05T16:01:54Z", validTime)

		require.NoError(t, err)
	})

	t.Run("returns error for actual that is not a time", func(t *testing.T) {
		err := CompareTime("2020-11-05T16:01:54Z", "NOT A TIME")

		require.EqualError(t, err, "NOT A TIME is not time.Time")
	})

	t.Run("returns error for different times", func(t *testing.T) {
		differentTime := validTime.Add(1 * time.Hour)
		require.NoError(t, err)

		err = CompareTime("2020-11-05T16:01:54Z", differentTime)

		require.EqualError(t, err, "expected 2020-11-05 16:01:54 +0000 UTC, but got 2020-11-05 17:01:54 +0000 UTC")
	})
}
