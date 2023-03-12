package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"tftp"
)

/*
  2 bytes                n bytes
 ┌───────┬──────────┬───────────────┐
 │OpCode │ Block #  │    payload    │
 └───────┴──────────┴───────────────┘
          2 bytes
*/

// all data packets must me acknowledged

// Data represents data flow
type Data struct {
	Block   uint16
	Payload io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 0, tftp.DatagramSize)
	b := bytes.NewBuffer(buf)

	// write operation
	err := binary.Write(b, binary.BigEndian, tftp.OpData)
	if err != nil {
		return []byte{}, err
	}

	// write block number
	// d.Block // decide later where and how to increase blockNumber
	// we can do it here if we reuse that struct which makes sense
	err = binary.Write(b, binary.BigEndian, d.Block)
	if err != nil {
		return []byte{}, err
	}

	_, err = io.CopyN(b, d.Payload, tftp.DataBlockSize)
	if err != nil && err != io.EOF {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {
	if l := len(p); l < 4 || l > tftp.DatagramSize {
		return errors.New("invalid DATA")
	}

	var opcode tftp.OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode)
	if err != nil || opcode != tftp.OpData {
		return errors.New("invalid DATA")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block)
	if err != nil {
		return errors.New("invalid DATA")
	}

	d.Payload = bytes.NewBuffer(p[4:])

	return nil
}
