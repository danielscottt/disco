package discoapi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
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

	getCont := regexp.MustCompile("/disco/api/get_container")
	rmCont := regexp.MustCompile("/disco/api/remove_container")
	addCont := regexp.MustCompile("/disco/api/add_container")

	switch {
	case p == "/disco/local/node_id":
		d.Reply([]byte(d.NodeId))
	case p == "/disco/api/get_containers":
		getContainers(d)
	case addCont.MatchString(p):
		addContainer(d, p, payload)
	case rmCont.MatchString(p):
		removeContainer(d, p)
	case getCont.MatchString(p):
		getContainer(d, p)
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

func getContainers(d *DiscoAPI) {
	response := []byte{'['}
	ls, err := ioutil.ReadDir(d.DataPath)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	for i, f := range ls {
		data, err := scanContainerFile(d, f.Name())
		if err != nil {
			d.Reply([]byte(err.Error()))
			return
		}
		response = append(response, data...)
		if i != (len(ls) - 1) {
			response = append(response, ',')
		}
	}
	response = append(response, ']')
	d.Reply(response)
}

func getContainer(d *DiscoAPI, p string) {
	name := getName(p)
	data, err := scanContainerFile(d, name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(data)
}

func addContainer(d *DiscoAPI, p string, payload []byte) {
	name := getName(p)
	path := fmt.Sprintf("%s/%s", d.DataPath, name)
	ioutil.WriteFile(path, payload, 644)
	d.Reply([]byte("success"))
}

func removeContainer(d *DiscoAPI, p string) {
	name := getName(p)
	path := fmt.Sprintf("%s/%s", d.DataPath, name)
	os.Remove(path)
	d.Reply([]byte("success"))
}

func scanContainerFile(d *DiscoAPI, name string) ([]byte, error) {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", d.DataPath, name))
	if err != nil {
		return data, err
	}
	return data, nil
}

func getName(path string) string {
	pathArr := strings.Split(path, "/")
	return pathArr[len(pathArr)-1]
}
