package discoapi

import (
	"log"
	"net"
)

var nodeId string

func StartListener(id string) {

	nodeId = id

	l, err := net.Listen("unix", "/var/run/disco.sock")
	if err != nil {
		log.Fatalf("Error opening Disco's socket", err)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error reading socket")
			continue
		}
		go handleSocketRequest(conn)
	}
}

func handleSocketRequest(c net.Conn) {
	buf := make([]byte, 512)
	e, err := c.Read(buf)
	if err != nil {
		return
	}
	req := buf[0:e]
	routeRequest(c, string(req))
}

func routeRequest(c net.Conn, req string) {
	switch req {
	case "/disco/local/nodeId":
		reply(c, []byte(nodeId))
	}
}

func reply(c net.Conn, response []byte) {
	_, err := c.Write(response)
	if err != nil {
		log.Println("Error in replying to request")
		return
	}
}
