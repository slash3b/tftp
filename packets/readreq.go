package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"

	"tftp"
)

/*
Read Request packet:

   2 bytes             1 byte          1 byte
   ┌───────┬──────────┬──────┬────────┬──────┐
   │OpCode │ Filename │  0   │ Mode   │   0  │
   └───────┴──────────┴──────┴────────┴──────┘
            n bytes           n bytes
*/

type ReadReq struct {
	Filename string
	Mode     tftp.TransferMode
}

func (r *ReadReq) MarshalBinary() ([]byte, error) {
	dataBytesLen := 2 + len(r.Filename) + 1 + len(r.Mode) + 1
	data := make([]byte, 0, dataBytesLen)

	buf := bytes.NewBuffer(data)

	err := binary.Write(buf, binary.BigEndian, tftp.OpRRQ)
	if err != nil {
		return []byte{}, err
	}

	_, err = buf.WriteString(r.Filename)
	if err != nil {
		return []byte{}, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return []byte{}, err
	}

	_, err = buf.WriteString(string(tftp.Octet))
	if err != nil {
		return []byte{}, err
	}

	err = buf.WriteByte(0)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (q *ReadReq) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)

	var code tftp.OpCode

	err := binary.Read(r, binary.BigEndian, &code) // read operation code
	if err != nil {
		return err
	}

	if code != tftp.OpRRQ {
		return errors.New("invalid RRQ")
	}

	q.Filename, err = r.ReadString(0) // read filename
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") // remove the 0-byte
	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ")
	}

	mode, err := r.ReadString(0) // read mode
	if err != nil {
		return errors.New("invalid RRQ")
	}

	mode = strings.TrimRight(mode, "\x00") // remove the 0-byte
	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}

	q.Mode = tftp.TransferMode(strings.ToLower(mode))
	if q.Mode != tftp.Octet {
		return errors.New("only binary transfers supported")
	}

	return nil
}
