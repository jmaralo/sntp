package message_test

import (
	"bytes"
	"testing"

	"github.com/jmaralo/sntp/internal/message"
	"github.com/jmaralo/sntp/internal/timestamp"
)

func TestMessage(t *testing.T) {
	t.Run("idempotency", func(t *testing.T) {
		msg := message.Message{
			LeapIndicator:  message.NoWarning,
			VersionNumber:  4,
			Mode:           message.Client,
			Stratum:        2,
			PollInterval:   6,
			Precision:      -20,
			RootDelay:      0x0001,
			RootDispersion: 0x0002,
			ReferenceID:    0x0003,
			ReferenceTimestamp: timestamp.Timestamp{
				Seconds:  0x0004,
				Fraction: 0x0005,
			},
			OriginateTimestamp: timestamp.Timestamp{
				Seconds:  0x0006,
				Fraction: 0x0007,
			},
			ReceiveTimestamp: timestamp.Timestamp{
				Seconds:  0x0008,
				Fraction: 0x0009,
			},
			TransmitTimestamp: timestamp.Timestamp{
				Seconds:  0x000a,
				Fraction: 0x000b,
			},
		}

		buffer := new(bytes.Buffer)

		err := msg.Write(buffer)
		if err != nil {
			t.Fatalf("error writing message: %v", err)
		}

		newMsg, err := message.Read(buffer)
		if err != nil {
			t.Fatalf("error reading message: %v", err)
		}

		if newMsg.LeapIndicator != msg.LeapIndicator {
			t.Errorf("leap indicator is not correct:\nexpected: %02b\ngot:      %02b\n", msg.LeapIndicator, newMsg.LeapIndicator)
		}

		if newMsg.VersionNumber != msg.VersionNumber {
			t.Errorf("version number is not correct:\nexpected: %03b\ngot:      %03b\n", msg.VersionNumber, newMsg.VersionNumber)
		}

		if newMsg.Mode != msg.Mode {
			t.Errorf("mode is not correct:\nexpected: %03b\ngot:      %03b\n", msg.Mode, newMsg.Mode)
		}

		if newMsg.Stratum != msg.Stratum {
			t.Errorf("stratum is not correct:\nexpected: %08b\ngot:      %08b\n", msg.Stratum, newMsg.Stratum)
		}

		if newMsg.PollInterval != msg.PollInterval {
			t.Errorf("poll interval is not correct:\nexpected: %08b\ngot:      %08b\n", msg.PollInterval, newMsg.PollInterval)
		}

		if newMsg.Precision != msg.Precision {
			t.Errorf("precision is not correct:\nexpected: %08b\ngot:      %08b\n", msg.Precision, newMsg.Precision)
		}

		if newMsg.RootDelay != msg.RootDelay {
			t.Errorf("root delay is not correct:\nexpected: %032b\ngot:      %032b\n", msg.RootDelay, newMsg.RootDelay)
		}

		if newMsg.RootDispersion != msg.RootDispersion {
			t.Errorf("root dispersion is not correct:\nexpected: %032b\ngot:      %032b\n", msg.RootDispersion, newMsg.RootDispersion)
		}

		if newMsg.ReferenceID != msg.ReferenceID {
			t.Errorf("reference ID is not correct:\nexpected: %032b\ngot:      %032b\n", msg.ReferenceID, newMsg.ReferenceID)
		}

		if newMsg.ReferenceTimestamp.Seconds != msg.ReferenceTimestamp.Seconds {
			t.Errorf("reference timestamp seconds are not correct:\nexpected: %032b\ngot:      %032b\n", msg.ReferenceTimestamp.Seconds, newMsg.ReferenceTimestamp.Seconds)
		}

		if newMsg.ReferenceTimestamp.Fraction != msg.ReferenceTimestamp.Fraction {
			t.Errorf("reference timestamp fraction is not correct:\nexpected: %032b\ngot:      %032b\n", msg.ReferenceTimestamp.Fraction, newMsg.ReferenceTimestamp.Fraction)
		}

		if newMsg.OriginateTimestamp.Seconds != msg.OriginateTimestamp.Seconds {
			t.Errorf("originate timestamp seconds are not correct:\nexpected: %032b\ngot:      %032b\n", msg.OriginateTimestamp.Seconds, newMsg.OriginateTimestamp.Seconds)
		}

		if newMsg.OriginateTimestamp.Fraction != msg.OriginateTimestamp.Fraction {
			t.Errorf("originate timestamp fraction is not correct:\nexpected: %032b\ngot:      %032b\n", msg.OriginateTimestamp.Fraction, newMsg.OriginateTimestamp.Fraction)
		}

		if newMsg.ReceiveTimestamp.Seconds != msg.ReceiveTimestamp.Seconds {
			t.Errorf("receive timestamp seconds are not correct:\nexpected: %032b\ngot:      %032b\n", msg.ReceiveTimestamp.Seconds, newMsg.ReceiveTimestamp.Seconds)
		}

		if newMsg.ReceiveTimestamp.Fraction != msg.ReceiveTimestamp.Fraction {
			t.Errorf("receive timestamp fraction is not correct:\nexpected: %032b\ngot:      %032b\n", msg.ReceiveTimestamp.Fraction, newMsg.ReceiveTimestamp.Fraction)
		}

		if newMsg.TransmitTimestamp.Seconds != msg.TransmitTimestamp.Seconds {
			t.Errorf("transmit timestamp seconds are not correct:\nexpected: %032b\ngot:      %032b\n", msg.TransmitTimestamp.Seconds, newMsg.TransmitTimestamp.Seconds)
		}

		if newMsg.TransmitTimestamp.Fraction != msg.TransmitTimestamp.Fraction {
			t.Errorf("transmit timestamp fraction is not correct:\nexpected: %032b\ngot:      %032b\n", msg.TransmitTimestamp.Fraction, newMsg.TransmitTimestamp.Fraction)
		}
	})
}
