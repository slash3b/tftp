package packets_test

import (
	"fmt"
	"testing"

	"tftp"
	"tftp/packets"
)

func TestReadReq(t *testing.T) {
	req := packets.ReadReq{
		Filename: "aaa",
		Mode:     tftp.Octet,
	}

	b, err := req.MarshalBinary()
	if err != nil {
		t.Error(fmt.Sprintf("unexpected err: %s", err))
	}

	t.Logf("%b", b)
}
