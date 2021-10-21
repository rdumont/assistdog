package defaults

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseString(t *testing.T) {
	res, err := ParseString("abc")

	assert.NoError(t, err)
	assert.Equal(t, "abc", res)
}

func TestParseInt(t *testing.T) {
	t.Run("parses valid integer", func(t *testing.T) {
		res, err := ParseInt("123")

		require.NoError(t, err)
		require.Equal(t, 123, res)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		_, err := ParseInt("abc")

		require.EqualError(t, err, `strconv.Atoi: parsing "abc": invalid syntax`)
	})
}

func TestParseInt32(t *testing.T) {
	t.Run("parses valid int32", func(t *testing.T) {
		res, err := ParseInt32("123")

		require.NoError(t, err)
		require.Equal(t, int32(123), res)
	})

	t.Run("returns error for invalid integer", func(t *testing.T) {
		_, err := ParseInt32("abc")

		require.EqualError(t, err, `strconv.ParseInt: parsing "abc": invalid syntax`)
	})

	t.Run("returns error for int32 overflow", func(t *testing.T) {
		_, err := ParseInt32("2147483648")

		require.EqualError(t, err, `strconv.ParseInt: parsing "2147483648": value out of range`)
	})
}

func TestParseTime(t *testing.T) {
	t.Run("parses supported layouts", func(t *testing.T) {
		expected, err := time.Parse(time.RFC3339Nano, "2020-11-05T16:01:54.0123Z")
		require.NoError(t, err)

		cases := []struct {
			layout    string
			raw       string
			tolerance time.Duration
		}{
			{layout: time.RFC822, raw: "05 Nov 20 16:01 MST", tolerance: time.Minute},
			{layout: time.RFC3339, raw: "2020-11-05T16:01:54Z", tolerance: time.Second},
			{layout: time.RFC3339, raw: "2020-11-05T16:01:54+00:00", tolerance: time.Second},
			{layout: time.RFC3339Nano, raw: "2020-11-05T16:01:54.0123Z", tolerance: time.Nanosecond},
		}

		for _, tc := range cases {
			t.Run(tc.layout, func(t *testing.T) {
				res, err := ParseTime(tc.raw)

				require.NoError(t, err)
				assert.WithinDuration(t, expected, res.(time.Time), tc.tolerance)
			})
		}
	})

	t.Run("returns error for unsupported layout", func(t *testing.T) {
		_, err := ParseTime("abc")

		require.EqualError(t, err, "unrecognized time format abc")
	})
}
