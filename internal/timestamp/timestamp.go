package timestamp

import (
	"time"
)

type Timestamp struct {
	Seconds  uint32
	Fraction uint32
}

func FromTime(t time.Time) Timestamp {
	seconds := uint64(t.Unix())
	fraction := ((uint64(t.UnixNano()+1) - (seconds * 1e9)) << 32) / 1e9
	return Timestamp{
		Seconds:  uint32(seconds + 2208988800),
		Fraction: uint32(fraction),
	}
}

func (t Timestamp) ToTime() time.Time {
	seconds := t.Seconds - 2208988800
	nanos := (uint64(t.Fraction) * 1e9) >> 32
	return time.Unix(int64(seconds), int64(nanos))
}
