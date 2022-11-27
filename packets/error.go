package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"tftp"
)

/*
  2 bytes              n bytes       1 byte
 ┌───────┬──────────┬───────────────┬──────┐
 │OpCode │ ErrCode  │  Message      │  0   │
 └───────┴──────────┴───────────────┴──────┘
           2 bytes
*/

type Err struct {
	Error   tftp.ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {
	// operation code + error code + message + 0 byte
	bufCap := 2 + 2 + len(e.Message) + 1
	b := new(bytes.Buffer)
	b.Grow(bufCap)

	err := binary.Write(b, binary.BigEndian, tftp.OpErr) // write operation code
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error) // write error code
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message) // write message
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // write 0 byte
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *Err) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)
	var code tftp.OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read operation code
	if err != nil {
		return err
	}

	if code != tftp.OpErr {
		return errors.New("invalid ERROR")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error) // read error message
	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)
	e.Message = strings.TrimRight(e.Message, "\x00") // remove the 0-byte
	return err
}
