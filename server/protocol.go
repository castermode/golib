package server

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	log "github.com/nicholaskh/log4go"
)

type Protocol struct {
	net.Conn
}

func NewProtocol() *Protocol {
	this := new(Protocol)

	return this
}

func (this *Protocol) SetConn(conn net.Conn) {
	this.Conn = conn
}

//len+payload
func (this *Protocol) Marshal(payload []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	dataBuff := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(payload)))
	binary.Write(dataBuff, binary.BigEndian, payload)
	buf.Write(dataBuff.Bytes())

	return buf.Bytes()
}

func (this *Protocol) Read() ([]byte, error) {
	buf := make([]byte, 4)
	n, err := this.Conn.Read(buf)
	if err != nil {
		log.Error("[Protocol] Read data length error: %s", err.Error())
		return []byte{}, err
	}
	buf = buf[0:n]
	b_buf := bytes.NewBuffer(buf)
	var dataLength int32
	binary.Read(b_buf, binary.BigEndian, &dataLength)
	data := make([]byte, dataLength)
	n, err = this.Conn.Read(data)
	if err != nil {
		log.Error("[Protocol] Read data error: %s", err.Error())
		return []byte{}, err
	}
	if int32(n) != dataLength {
		err = errors.New("[Protocol] Data payload length not correct")
		log.Error("[Protocol] Data payload length not correct, expect %d, give %d", dataLength, n)
		return []byte{}, err
	}
	return data, nil
}
