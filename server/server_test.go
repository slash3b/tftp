package server_test

import (
	"flag"
	"log"
	"testing"

	"tftp/server"
)

//var (
//	address = flag.String("a", "127.0.0.1:69", "listen address")
//	payload = flag.String("p", "payload", "file to serve to clients")
//)

func TestServer(t *testing.T) {
	flag.Parse()

	// p, err := ioutil.ReadFile("./payload")

	//if err != nil {
	//	log.Fatal(err)
	//}

	s := server.Server{}

	log.Fatal(s.Listen("127.0.0.1:8080"))
}
