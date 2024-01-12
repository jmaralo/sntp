package timestamp_test

import (
	"testing"
	"time"

	"github.com/jmaralo/sntp/internal/timestamp"
)

func TestTimestamp(t *testing.T) {
	const (
		EXPECTED_SECONDS  = 0xdcd2aa81
		EXPECTED_FRACTION = 0x7bff7370
	)
	instant := time.Date(2017, 5, 26, 13, 22, 9, 484366621, time.UTC)

	// This test checks that when applying the conversion back and forth, the instant is kept.
	t.Run("idempotency", func(t *testing.T) {
		ts := timestamp.FromTime(instant)
		newInstant := ts.ToTime()

		if !newInstant.Equal(instant) {
			t.Errorf("conversion does not work both ways:\n-- %v\n-> %v\n-> %v\n", instant, ts, newInstant)
		}

		firstTimestamp := timestamp.Timestamp{
			Seconds:  EXPECTED_SECONDS,
			Fraction: EXPECTED_FRACTION,
		}
		newInstant = firstTimestamp.ToTime()
		newTimestamp := timestamp.FromTime(newInstant)

		if newTimestamp.Seconds != firstTimestamp.Seconds {
			t.Errorf("seconds are not correct:\nexpected: %032b\ngot:      %032b\n", firstTimestamp.Seconds, newTimestamp.Seconds)
		}

		if newTimestamp.Fraction != firstTimestamp.Fraction {
			t.Errorf("fraction is not correct:\nexpected: %032b\ngot:      %032b\n", firstTimestamp.Fraction, newTimestamp.Fraction)
		}
	})

	t.Run("unix to NTP", func(t *testing.T) {

		timestamp := timestamp.FromTime(instant)

		if timestamp.Seconds != EXPECTED_SECONDS {
			t.Errorf("seconds are not correct:\nexpected: %032b\ngot:      %032b\n", EXPECTED_SECONDS, timestamp.Seconds)
		}

		if timestamp.Fraction != EXPECTED_FRACTION {
			t.Errorf("fraction is not correct:\nexpected: %032b\ngot:      %032b\n", EXPECTED_FRACTION, timestamp.Fraction)
		}
	})

	t.Run("NTP to unix", func(t *testing.T) {
		timestamp := timestamp.Timestamp{
			Seconds:  EXPECTED_SECONDS,
			Fraction: EXPECTED_FRACTION,
		}
		newInstant := timestamp.ToTime()

		if !newInstant.Equal(instant) {
			t.Errorf("instant is not correct:\nexpected: %v\ngot:      %v\n", instant, newInstant)
		}
	})
}
