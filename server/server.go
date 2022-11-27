package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"tftp"
	"tftp/packets"
)

type Server struct {
	// Retries a number of times to wait for packet acknowledgement
	Retries uint8
	Payload []byte
	Timeout time.Duration
}

func (s *Server) Listen(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	log.Printf("listening on %s ...", conn.LocalAddr())

	return s.Serve(conn)
}

func (s *Server) Serve(cn net.PacketConn) error {

	// basically verify that server state is valid ??? wtf ?
	// take all request

	// assuming all request are read requests

	for {
		buf := make([]byte, 0, tftp.DatagramSize)
		_, addr, err := cn.ReadFrom(buf)
		if err != nil {
			log.Printf("ERROR: %s", err)
			return err
		}

		rrq := packets.ReadReq{}
		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("ERROR: %s", err)
			continue
		}

		go s.handle(addr.String(), rrq)
	}
}

func (s *Server) handle(addr string, rrq packets.ReadReq) {
	log.Printf("[%s] requested file: %s", addr, rrq.Filename)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Printf("[%s] dial: %v", addr, err)
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("%s", err)
		}
	}()

	ackPkt := packets.Ack(0)
	errPkt := packets.Err{}
	dataPkt := packets.Data{Payload: bytes.NewReader(s.Payload)}

	buf := make([]byte, tftp.DatagramSize)

nextpacket:
	for n := tftp.DatagramSize; n == tftp.DatagramSize; {
		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v", addr, err)
			return
		}

	retry:
		for i := s.Retries; i > 0; i-- {
			n, err = conn.Write(data) // send the data packet
			if err != nil {
				log.Printf("[%s] write: %v", addr, err)
				return
			}

			// wait for the client's ACK packet
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))
			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue retry
				}
				log.Printf("[%s] waiting for ACK: %v", addr, err)
				return
			}

			switch {
			case ackPkt.UnmarshalBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block {
					// received ACK; send next data packet
					continue nextpacket
				}
			case errPkt.UnmarshalBinary(buf) == nil:
				log.Printf("[%s] received error: %v",
					addr, errPkt.Message)
				return
			default:
				log.Printf("[%s] bad packet", addr)
			}
		}
		log.Printf("[%s] exhausted retries", addr)
		return
	}
	log.Printf("[%s] sent %d blocks", addr, dataPkt.Block)
}
