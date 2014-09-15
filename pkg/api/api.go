package discoapi

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

type DiscoAPI struct {
	NodeId     string
	SocketPath string
	DataPath   string

	listener   net.Listener
	connection net.Conn
}

func NewDiscoAPI(id string) (*DiscoAPI, error) {

	var (
		d                    *DiscoAPI
		socketPath, dataPath string
	)

	if os.Getenv("DISCO_SOCKET") != "" {
		socketPath = os.Getenv("DISCO_SOCKET")
	} else {
		return d, errors.New("Disco socket path not set. Cannot start.")
	}
	if os.Getenv("DISCO_DATA_PATH") != "" {
		dataPath = os.Getenv("DISCO_DATA_PATH")
	} else {
		return d, errors.New("Disco data path not set. Cannot start.")
	}

	log.Print("Disco socket Path: [", socketPath, "]")
	log.Print("Disco data Path: [", dataPath, "]")

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return d, err
	}

	d = &DiscoAPI{
		listener:   l,
		NodeId:     id,
		SocketPath: socketPath,
		DataPath:   dataPath,
	}

	return d, nil
}

func (d *DiscoAPI) Start() {

	var err error

	for {
		d.connection, err = d.listener.Accept()
		if err != nil {
			log.Println("Error reading socket")
			continue
		}
		go d.handleSocketRequest()
	}
}

func (d *DiscoAPI) Stop() {
	if d.connection != nil {
		d.connection.Close()
	}
	d.listener.Close()
}

func (d *DiscoAPI) handleSocketRequest() {

	defer d.connection.Close()

	buf := make([]byte, 1024)
	e, err := d.connection.Read(buf)
	if err != nil {
		return
	}
	req := buf[0:e]

	splitPayload := bytes.SplitN(req, []byte("\n"), 2)

	var path, payload []byte
	path = splitPayload[0]
	if len(splitPayload) > 2 {
		payload = splitPayload[1]
	}

	d.routeRequest(path, payload)
}

func (d *DiscoAPI) routeRequest(path, payload []byte) {

	p := string(path)
	log.Print("Received request [", p, "]")

	switch p {
	case "/disco/local/node_id":
		d.reply([]byte(d.NodeId))
	case "/disco/api/add_container":
		addContainer(d, payload)
	default:
		log.Print("Request path [", p, "] not found")
		err := fmt.Sprintf("Error: Invalid request path [%s]", p)
		d.reply([]byte(err))
	}
}

func (d *DiscoAPI) reply(response []byte) {
	_, err := d.connection.Write(response)
	if err != nil {
		log.Println("Error in replying to request")
		return
	}
}

func addContainer(d *DiscoAPI, payload []byte) {
	log.Println((*d).NodeId, "Payload:", string(payload))
}
