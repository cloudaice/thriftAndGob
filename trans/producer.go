package trans

import (
	"compress/zlib"
	"encoding/gob"
	"net"
)

type Producer struct {
	conn   net.Conn
	client *gob.Encoder
	compW  *zlib.Writer
}

func NewProducer(Addr string, compression bool) (*Producer, error) {
	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		return nil, err
	}

	if !compression {
		enc := gob.NewEncoder(conn)
		return &Producer{conn, enc, nil}, nil
	}
	compW, err := zlib.NewWriterLevel(conn, zlib.BestCompression)
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(compW)
	return &Producer{conn, enc, compW}, nil
}

func (pdr *Producer) SendOne(msg *LogEntry) error {
	return pdr.SendArray([]*LogEntry{msg})
}

func (pdr *Producer) SendArray(msgs []*LogEntry) error {
	return pdr.client.Encode(msgs)
}

func (pdr *Producer) Close() error {
	if pdr.compW == nil {
		return pdr.conn.Close()
	}
	err := pdr.compW.Close()
	if err != nil {
		return err
	}
	return pdr.conn.Close()
}
