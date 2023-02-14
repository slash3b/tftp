package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"tftp"
)

/*
Write Request packet:

   2 bytes             1 byte          1 byte
   ┌───────┬──────────┬──────┬────────┬──────┐
   │OpCode │ Filename │  0   │ Mode   │   0  │
   └───────┴──────────┴──────┴────────┴──────┘
            n bytes           n bytes
*/

type WriteReq struct {
	Filename string
	Mode     tftp.TransferMode
}

func (w *WriteReq) MarshalBinary() ([]byte, error) {
	dataBytesLen := 2 + len(w.Filename) + 1 + len(w.Mode) + 1
	data := make([]byte, 0, dataBytesLen)

	buf := bytes.NewBuffer(data)

	err := binary.Write(buf, binary.BigEndian, tftp.OpWRQ)
	if err != nil {
		return []byte{}, err
	}

	_, err = buf.WriteString(w.Filename)
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

func (w *WriteReq) UnmarshalBinary(p []byte) error {
	buf := bytes.NewBuffer(p)

	var code tftp.OpCode

	err := binary.Read(buf, binary.BigEndian, &code) // read operation code
	if err != nil {
		return err
	}

	if code != tftp.OpWRQ {
		return fmt.Errorf("invalid write request, expected op code %d, received %d", tftp.OpWRQ, code)
	}

	w.Filename, err = buf.ReadString(0) // read filename
	if err != nil {
		return errors.New("could not read filename")
	}

	w.Filename = strings.TrimRight(w.Filename, "\x00")
	if len(w.Filename) == 0 {
		return errors.New("invalid WRQ, filename is empty")
	}

	mode, err := buf.ReadString(0)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid WRQ: could not read mode, due to err, %#v", err))
	}

	mode = strings.TrimRight(mode, "\x00")
	if len(mode) == 0 {
		return errors.New("invalid WRQ, mode is empty upon trimming")
	}

	w.Mode = tftp.TransferMode(strings.ToLower(mode))
	if w.Mode != tftp.Octet {
		return errors.New("only binary transfers supported")
	}

	fmt.Printf("%#v", w)

	return nil
}
