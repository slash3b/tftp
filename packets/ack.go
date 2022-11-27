package packets

import (
	"bytes"
	"encoding/binary"
	"errors"

	"tftp"
)

/*
   2 bytes
   ┌───────┬──────────┐
   │OpCode │ Block #  │
   └───────┴──────────┘
             2 bytes
*/

type Ack uint16

func (a Ack) MarshalBinary() ([]byte, error) {
	bufCap := 2 + 2 // operation code + block number

	b := new(bytes.Buffer)

	b.Grow(bufCap)

	err := binary.Write(b, binary.BigEndian, tftp.OpAck) // write operation code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a) // write block number
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil

}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var code tftp.OpCode
	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code) // read operation code
	if err != nil {
		return err
	}

	if code != tftp.OpAck {
		return errors.New("invalid ACK")
	}

	return binary.Read(r, binary.BigEndian, a) // read block number
}
