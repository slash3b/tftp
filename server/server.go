package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"tftp"
	"tftp/packets"
)

type Server struct {
	// Retries a number of times to wait for packet acknowledgement
	//Retries uint8
	//Payload []byte
	//Timeout time.Duration
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
		buf := make([]byte, tftp.DatagramSize)

		n, addr, err := cn.ReadFrom(buf)
		fmt.Printf("%d bytes came from %s \n", n, addr)
		fmt.Printf("raw: %08b\n", buf[:n])
		if err != nil {
			log.Printf("ERROR 1: %s", err)
			return err
		}

		if n == 0 {
			continue
		}

		rrq := packets.ReadReq{}
		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("ERROR 2: %s", err)
			continue
		}

		fmt.Println("handling ...")
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

	// ackPkt := packets.Ack(0)
	// errPkt := packets.Err{}

	//dir, err := os.Getwd()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	fileBytes, err := os.ReadFile("lorem.txt")
	if err != nil {
		fmt.Printf("err %v \n", err)
		return
	}
	rd := bytes.NewReader(fileBytes)

	blockN := uint16(1)
	for {
		// every packet should be acknowledged
		buf := make([]byte, tftp.BlockSize)
		//io.CopyN()
		n, err := rd.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		dataPkt := packets.Data{
			Block:   blockN,
			Payload: bytes.NewReader(buf[:n]),
		}
		//fmt.Println("sending:", string(buf[:n]))

		b, err := dataPkt.MarshalBinary()
		if err != nil {
			fmt.Println(err)
			return
		}

		// use conn.SetDeadline()
		// check for os.ErrDeadlineExceeded
		// extend deadline
		fmt.Printf("sending block %d from %s to addr: %s \n", blockN, conn.LocalAddr(), conn.RemoteAddr())
		// fmt.Printf("raw bytes: %b \n", b)
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

		if n < tftp.BlockSize {
			fmt.Println(">>>>>>>>>>>>>>>last packet was sent")
			return
		}
		blockN++
	}

}
