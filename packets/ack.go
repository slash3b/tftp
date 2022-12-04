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

// question: do we really need to unmarshall this?
// looks like it can be generated once
func (a *Ack) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	err := binary.Write(buf, binary.BigEndian, tftp.OpAck)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

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
