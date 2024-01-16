package sntp

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/jmaralo/sntp/internal/message"
	"github.com/jmaralo/sntp/internal/timestamp"
)

// READ_BUF_SIZE is the size of the buffer used to read from the connection.
const READ_BUF_SIZE = 512

// Server is a unicast SNTP version 4 server.
type Server struct {
	conn      *net.UDPConn
	reference timestamp.Timestamp
}

// NewUnicast creates a new unicast SNTP version 4 server on the given address.
func NewUnicast(network string, laddr *net.UDPAddr) (*Server, error) {
	conn, err := net.ListenUDP(network, laddr)
	return &Server{
		conn:      conn,
		reference: timestamp.FromTime(time.Now()),
	}, err
}

// ListenAndServe listens on the server connection for messages and responds to them.
func (server *Server) ListenAndServe() error {
	for {
		readBuf := make([]byte, READ_BUF_SIZE)
		n, addr, err := server.conn.ReadFrom(readBuf)
		if err != nil {
			return err
		}

		go func() {
			req := udpRequest{
				address: addr,
				data:    readBuf[:n],
			}

			err := req.serve(server.conn, server.reference)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error serving request to %s: %v\n", req.address, err)
			}
		}()
	}
}

// udpRequest is an individual request made to the server.
type udpRequest struct {
	address net.Addr
	data    []byte
}

// serve handles the udpRequest and responds to it on the given connection.
func (request udpRequest) serve(conn *net.UDPConn, ref timestamp.Timestamp) error {
	readBuf := bytes.NewReader(request.data)
	msg, err := message.Read(readBuf)
	if err != nil {
		return err
	}

	switch msg.Mode {
	case message.Client:
		msg.Mode = message.Server
	default:
		msg.Mode = message.SymmetricPassive
	}

	msg.LeapIndicator = message.NoWarning
	msg.Stratum = 1
	msg.Precision = -20
	msg.RootDelay = 0
	msg.RootDispersion = 0
	msg.ReferenceID = 0x4c4f434c // LOCL
	msg.ReferenceTimestamp = ref
	msg.OriginateTimestamp = msg.TransmitTimestamp
	msg.ReceiveTimestamp = timestamp.FromTime(time.Now())
	msg.TransmitTimestamp = msg.ReceiveTimestamp

	writeBuf := new(bytes.Buffer)
	err = msg.Write(writeBuf)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(writeBuf.Bytes(), request.address)
	return err
}
