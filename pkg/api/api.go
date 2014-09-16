package discoapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func NewDiscoAPI(id, dataPath, socketPath string) (*DiscoAPI, error) {

	var d *DiscoAPI

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
	if len(splitPayload) >= 2 {
		payload = splitPayload[1]
	}

	d.routeRequest(path, payload)
}

func (d *DiscoAPI) routeRequest(path, payload []byte) {

	p := string(path)
	log.Print("Received request [", p, "]")

	switch p {
	case "/disco/local/node_id":
		d.Reply([]byte(d.NodeId))
	case "/disco/api/add_container":
		addContainer(d, payload)
	case "/disco/api/remove_container":
		removeContainer(d, payload)
	default:
		log.Print("Request path [", p, "] not found")
		err := fmt.Sprintf("Error: Invalid request path [%s]", p)
		d.Reply([]byte(err))
	}
}

func (d *DiscoAPI) Reply(response []byte) {
	_, err := d.connection.Write(response)
	if err != nil {
		log.Println("Error in replying to request")
		return
	}
}

func addContainer(d *DiscoAPI, payload []byte) {
	splitPayload := bytes.SplitN(payload, []byte("\n"), 2)
	id := splitPayload[0]
	p := splitPayload[1]
	path := fmt.Sprintf("%s/%s:%s", d.DataPath, d.NodeId, id)
	ioutil.WriteFile(path, p, 644)
	d.Reply([]byte("success"))
}

func removeContainer(d *DiscoAPI, payload []byte) {
	path := fmt.Sprintf("%s/%s:%s", d.DataPath, d.NodeId, string(payload))
	os.Remove(path)
	d.Reply([]byte("success"))
}
