package message

import (
	"io"

	"github.com/jmaralo/sntp/internal/timestamp"
)

type leapIndicator uint8

const (
	NoWarning leapIndicator = iota
	LastMinute61
	LastMinute59
	Alarm
)

type mode uint8

const (
	Reserved mode = iota
	SymmetricActive
	SymmetricPassive
	Client
	Server
	Broadcast
	ReservedNTPControl
	ReservedPrivate
)

const MESSAGE_SIZE = 48

// Message is a NTP message as described in [RFC 1305].
//
// [RFC 1305]: https://tools.ietf.org/html/rfc1305#section-3.1
type Message struct {
	// LeapIndicator is a two-bit code warning of an impending leap second to be inserted/deleted in the last minute of the current day.
	LeapIndicator leapIndicator
	// VersionNumber is a three-bit integer indicating the NTP/SNTP version number.
	VersionNumber uint8
	// Mode is a three-bit integer indicating the mode, which represents the purpose or function of the message.
	Mode mode
	// Stratum is a eight-bit integer indicating the stratum level of the local clock.
	Stratum uint8
	// PollInterval is a eight-bit signed integer indicating the maximum interval between successive messages, in log2 seconds.
	PollInterval int8
	// Precision is a eight-bit signed integer indicating the precision of the local clock, in log2 seconds.
	Precision int8
	// RootDelay is a 32-bit signed fixed-point number indicating the total roundtrip delay to the primary reference source, in seconds.
	RootDelay int32
	// RootDispersion is a 32-bit unsigned fixed-point number indicating the maximum error due to the clock frequency tolerance, in seconds.
	RootDispersion uint32
	// ReferenceID is a 32-bit code identifying the particular reference source.
	ReferenceID uint32
	// ReferenceTimestamp is a 64-bit timestamp identifying the last update of the local clock.
	ReferenceTimestamp timestamp.Timestamp
	// OriginateTimestamp is a 64-bit timestamp indicating the local time at which the request departed the client for the server.
	OriginateTimestamp timestamp.Timestamp
	// ReceiveTimestamp is a 64-bit timestamp indicating the local time at which the request arrived at the server.
	ReceiveTimestamp timestamp.Timestamp
	// TransmitTimestamp is a 64-bit timestamp indicating the local time at which the reply departed the server for the client.
	TransmitTimestamp timestamp.Timestamp
	// TODO: Authenticator
}

// Read reads a Message from a [io.Reader].
//
// Read will return a partial message if the reader does not contain enough data.
func Read(reader io.Reader) (Message, error) {
	data := make([]byte, 0, MESSAGE_SIZE)

	readBuf := make([]byte, MESSAGE_SIZE)
	totalRead := 0
	var readErr error
	var n int
	for totalRead < MESSAGE_SIZE {
		n, readErr = reader.Read(readBuf[:MESSAGE_SIZE-totalRead])
		data = append(data, readBuf[:n]...)
		totalRead += n
		if readErr != nil {
			break
		}
	}

	message := Message{}

	if len(data) >= 1 {
		message.LeapIndicator = leapIndicator(data[0] >> 6 & 0b11)
		message.VersionNumber = data[0] >> 3 & 0b111
		message.Mode = mode(data[0] & 0b111)
	}

	if len(data) >= 2 {
		message.Stratum = uint8(data[1])
	}

	if len(data) >= 3 {
		message.PollInterval = int8(data[2])
	}

	if len(data) >= 4 {
		message.Precision = int8(data[3])
	}

	if len(data) >= 8 {
		message.RootDelay = int32(data[4])<<24 | int32(data[5])<<16 | int32(data[6])<<8 | int32(data[7])
	}

	if len(data) >= 12 {
		message.RootDispersion = uint32(data[8])<<24 | uint32(data[9])<<16 | uint32(data[10])<<8 | uint32(data[11])
	}

	if len(data) >= 16 {
		message.ReferenceID = uint32(data[12])<<24 | uint32(data[13])<<16 | uint32(data[14])<<8 | uint32(data[15])
	}

	if len(data) >= 24 {
		message.ReferenceTimestamp = timestamp.Timestamp{
			Seconds:  uint32(data[16])<<24 | uint32(data[17])<<16 | uint32(data[18])<<8 | uint32(data[19]),
			Fraction: uint32(data[20])<<24 | uint32(data[21])<<16 | uint32(data[22])<<8 | uint32(data[23]),
		}
	}

	if len(data) >= 32 {
		message.OriginateTimestamp = timestamp.Timestamp{
			Seconds:  uint32(data[24])<<24 | uint32(data[25])<<16 | uint32(data[26])<<8 | uint32(data[27]),
			Fraction: uint32(data[28])<<24 | uint32(data[29])<<16 | uint32(data[30])<<8 | uint32(data[31]),
		}
	}

	if len(data) >= 40 {
		message.ReceiveTimestamp = timestamp.Timestamp{
			Seconds:  uint32(data[32])<<24 | uint32(data[33])<<16 | uint32(data[34])<<8 | uint32(data[35]),
			Fraction: uint32(data[36])<<24 | uint32(data[37])<<16 | uint32(data[38])<<8 | uint32(data[39]),
		}
	}

	if len(data) >= 48 {
		message.TransmitTimestamp = timestamp.Timestamp{
			Seconds:  uint32(data[40])<<24 | uint32(data[41])<<16 | uint32(data[42])<<8 | uint32(data[43]),
			Fraction: uint32(data[44])<<24 | uint32(data[45])<<16 | uint32(data[46])<<8 | uint32(data[47]),
		}
	}

	return message, readErr
}

// Write writes a message to a [io.Writer].
func (message *Message) Write(writer io.Writer) error {
	_, err := writer.Write([]byte{
		(byte(message.LeapIndicator) << 6) | (message.VersionNumber << 3) | byte(message.Mode),
		byte(message.Stratum),
		byte(message.PollInterval),
		byte(message.Precision),
		byte(message.RootDelay >> 24 & 0xff),
		byte(message.RootDelay >> 16 & 0xff),
		byte(message.RootDelay >> 8 & 0xff),
		byte(message.RootDelay & 0xff),
		byte(message.RootDispersion >> 24 & 0xff),
		byte(message.RootDispersion >> 16 & 0xff),
		byte(message.RootDispersion >> 8 & 0xff),
		byte(message.RootDispersion & 0xff),
		byte(message.ReferenceID >> 24 & 0xff),
		byte(message.ReferenceID >> 16 & 0xff),
		byte(message.ReferenceID >> 8 & 0xff),
		byte(message.ReferenceID & 0xff),
	})
	if err != nil {
		return err
	}

	err = message.ReferenceTimestamp.Write(writer)
	if err != nil {
		return err
	}

	err = message.OriginateTimestamp.Write(writer)
	if err != nil {
		return err
	}

	err = message.ReceiveTimestamp.Write(writer)
	if err != nil {
		return err
	}

	return message.TransmitTimestamp.Write(writer)
}
