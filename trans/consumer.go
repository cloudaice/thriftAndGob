package trans

import (
	"compress/zlib"
	"encoding/gob"
	"io"
	"net"
)

type Consumer struct {
	conn   net.Conn
	server *gob.Decoder
	compR  io.ReadCloser
}

// NewConsumer return a pointer of Consumer instance,
// compression is a optional, it can be true if need zlib compression
// else false.
func NewConsumer(conn net.Conn, compression bool) (*Consumer, error) {
	if !compression {
		dec := gob.NewDecoder(conn)
		return &Consumer{conn, dec, nil}, nil
	}

	compR, err := zlib.NewReader(conn)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(compR)
	return &Consumer{conn, dec, compR}, nil
}

// Receive return a list of *LogEntry,
// if can't read more than 1024 length once a time.
func (cmr *Consumer) Receive() ([]*LogEntry, error) {
	buf := make([]*LogEntry, 1024)
	err := cmr.server.Decode(&buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (cmr *Consumer) Close() error {
	if cmr.compR == nil {
		return cmr.conn.Close()
	}

	err := cmr.compR.Close()
	if err != nil {
		return err
	}
	return cmr.conn.Close()
}
