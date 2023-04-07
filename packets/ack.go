package packets

import (
	"bytes"
	"encoding/binary"
	"errors"

	"tftp"
)

// Ack layout:
//
//	 2 bytes
//	┌───────┬──────────┐
//	│OpCode │ Block #  │
//	└───────┴──────────┘
//			 2 bytes
//

type Ack uint16

func (a *Ack) MarshalBinary() ([]byte, error) {
	b := make([]byte, 4)

	binary.BigEndian.PutUint16(b, uint16(tftp.OpAck))
	binary.BigEndian.PutUint16(b[2:], uint16(*a))

	return b, nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	r := bytes.NewReader(p)

	var code tftp.OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read operation code
	if err != nil {
		return err
	}

	if code != tftp.OpAck {
		return errors.New("invalid ACK")
	}

	return binary.Read(r, binary.BigEndian, a) // read block number
}
