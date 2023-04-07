package packets_test

import (
	"testing"
	"tftp/packets"

	"github.com/stretchr/testify/assert"
)

func TestAckMarshall(t *testing.T) {
	ackBlock := packets.Ack(123)
	b, err := ackBlock.MarshalBinary()

	assert.NoError(t, err)
	expectedBytes := []byte{0x00, 0x04, 0x00, 0x7b}
	assert.Equal(t, expectedBytes, b)
}

func TestAckUnmarshall(t *testing.T) {
	b := []byte{0x00, 0x04, 0x00, 0x7b}

	var ack packets.Ack
	err := ack.UnmarshalBinary(b)
	assert.NoError(t, err)

	assert.Equal(t, ack, packets.Ack(0x7b))
}
