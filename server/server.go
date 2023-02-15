package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"tftp"
	"tftp/packets"
)

type Server struct{}

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
	for {
		buf := make([]byte, tftp.DatagramSize)

		n, addr, err := cn.ReadFrom(buf)
		if err != nil {
			log.Printf("ERROR 1: %s", err)
			return err
		}

		if n == 0 {
			continue
		}

		var opcode tftp.OpCode

		r := bytes.NewReader(buf[:n])
		err = binary.Read(r, binary.BigEndian, &opcode)
		if err != nil {
			return err
		}

		switch opcode {
		case tftp.OpWRQ:
			wrq := packets.WriteReq{}
			err = wrq.UnmarshalBinary(buf[:n])
			if err != nil {
				log.Printf("write err: %s", err)
				continue
			}
			go s.handleWrite(addr.String(), wrq)
		case tftp.OpRRQ:
			rrq := packets.ReadReq{}
			err = rrq.UnmarshalBinary(buf[:n])
			if err != nil {
				log.Printf("read err: %s", err)
				continue
			}
			go s.handleRead(addr.String(), rrq)
		default:
			n, err = cn.WriteTo([]byte{}, addr)
			if err != nil {
				log.Printf("err %s, %d bytes sent back to client", err, n)
			}
		}

		fmt.Println("handling ...")
	}
}

func (s *Server) handleWrite(addr string, wrq packets.WriteReq) {

}

func (s *Server) handleRead(addr string, rrq packets.ReadReq) {
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

	fileBytes, err := os.ReadFile(rrq.Filename)
	if err != nil {
		fmt.Printf("err %v \n", err)
		return
	}

	fileReader := bytes.NewReader(fileBytes)

	// fixme(1): implement retries

	blockN := uint16(1)
	for {
		buf := make([]byte, tftp.DataBlockSize)
		n, err := fileReader.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		dataPkt := packets.Data{
			Block:   blockN,
			Payload: bytes.NewReader(buf[:n]),
		}

		b, err := dataPkt.MarshalBinary()
		if err != nil {
			fmt.Println(err)
			return
		}

		// use conn.SetDeadline()
		// check for os.ErrDeadlineExceeded
		// extend deadline

		fmt.Printf("sending block %d from %s to addr: %s \n", blockN, conn.LocalAddr(), conn.RemoteAddr())
		f, err := conn.Write(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("sent %d bytes to client\n", f)

		resBuf := make([]byte, 1000)
		g, err := conn.Read(resBuf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("--------------------------------------------------")
		fmt.Printf("received %d bytes from %s to %s:\n", g, conn.RemoteAddr(), conn.LocalAddr())
		fmt.Printf("raw %b \n", resBuf[:g])
		var ack packets.Ack
		ackErr := ack.UnmarshalBinary(resBuf[:g])
		if ackErr != nil {
			fmt.Printf("unexpected payload %s, Acknowledge expected \n", string(resBuf[:g]))
			return
		}
		fmt.Printf("received ACK from client \n")

		if n < tftp.DataBlockSize {
			fmt.Println(">>>>>>>>>>>>>>>last packet was sent to Client. Terminating.")
			return
		}

		blockN++
	}
}
