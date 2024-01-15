package timestamp

import (
	"io"
	"time"
)

// Timestamp is a 64-bit fixed-point number following the NTP timestamp specification [RFC 1305].
//
// [RFC 1305]: https://tools.ietf.org/html/rfc1305#section-3.1
type Timestamp struct {
	Seconds  uint32
	Fraction uint32
}

// FromTime converts a [time.Time] to a Timestamp.
//
// Note this implementation returns a Fraction value off by 4-5 nanoseconds.
func FromTime(t time.Time) Timestamp {
	seconds := uint64(t.Unix())
	fraction := ((uint64(t.UnixNano()+1) - (seconds * 1e9)) << 32) / 1e9
	return Timestamp{
		Seconds:  uint32(seconds + 2208988800),
		Fraction: uint32(fraction),
	}
}

// ToTime converts a Timestamp to a [time.Time].
func (t Timestamp) ToTime() time.Time {
	seconds := t.Seconds - 2208988800
	nanos := (uint64(t.Fraction) * 1e9) >> 32
	return time.Unix(int64(seconds), int64(nanos))
}

func (t Timestamp) Write(writer io.Writer) error {
	data := []byte{
		byte(t.Seconds >> 24 & 0xff),
		byte(t.Seconds >> 16 & 0xff),
		byte(t.Seconds >> 8 & 0xff),
		byte(t.Seconds & 0xff),
		byte(t.Fraction >> 24 & 0xff),
		byte(t.Fraction >> 16 & 0xff),
		byte(t.Fraction >> 8 & 0xff),
		byte(t.Fraction & 0xff),
	}

	totalWrite := 0
	for totalWrite < len(data) {
		n, err := writer.Write(data[totalWrite:])
		totalWrite += n
		if err != nil {
			return err
		}
	}
	return nil
}
