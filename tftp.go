package tftp

const DatagramSize = 516
const BlockSize = DatagramSize - 4

// note: each block must be acknowledged before next comes in
// note: data block that less than 512 bytes means the last block
const DataBlockSize = 512

type TFTPBlock struct {
	// Header bears opcode
	Header
	DataBlockSize
}

/*
WRQ write to
RRQ read from
each data packet has a block number and starts from 1
*/

/*
WRQ
- srv responds with ack


RRQ
- srv responds with first data packet for read
*/

type OpCode uint16

/*
   opcode  operation
           1     Read request (RRQ)
           2     Write request (WRQ)
           3     Data (DATA)
           4     Acknowledgment (ACK)
           5     Error (ERROR)
*/
const (
	OpRRQ OpCode = iota + 1 // RRQ stands for read request
	OpWRQ
	OpData
	OpAck
	OpErr
)

type ErrCode uint16

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

type TransferMode string

const Octet TransferMode = "octet"
const Netascii TransferMode = "netascii"
