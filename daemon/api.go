package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

type DiscoAPI struct {
	NodeId   string
	DataPath string

	listener   net.Listener
	connection net.Conn
	stop       chan bool
}

func NewDiscoAPI(id string) (*DiscoAPI, error) {

	var d *DiscoAPI

	log.Print("Disco socket Path: [", config.Disco.DiscoSocket, "]")

	l, err := net.Listen("unix", "/var/run/disco.sock")
	if err != nil {
		return d, err
	}

	d = &DiscoAPI{
		listener: l,
		NodeId:   id,
		DataPath: PREFIX,
	}
	d.stop = make(chan bool, 1)

	return d, nil
}

func (d *DiscoAPI) Start() {

	var err error

	for {
		d.connection, err = d.listener.Accept()
		if err != nil {
			select {
			case <-d.stop:
				return
			default:
				log.Println("Accepting socket failed:", err.Error())
				continue
			}
		}
		go d.handleSocketRequest()
	}
}

func (d *DiscoAPI) Stop() {
	d.stop <- true
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

	getCont := regexp.MustCompile("/disco/api/get_container")
	rmCont := regexp.MustCompile("/disco/api/remove_container")
	addCont := regexp.MustCompile("/disco/api/add_container")

	switch {
	case p == "/disco/local/node_id":
		d.Reply([]byte(d.NodeId))
	case p == "/disco/api/get_containers":
		getContainersAPI(d)
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

func getContainersAPI(d *DiscoAPI) {
	response := []byte{'['}
	rep, err := persist.Read(d.DataPath + "/containers")
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	for i, f := range rep.Children {
		data, err := scanContainer(d, f)
		if err != nil {
			d.Reply([]byte(err.Error()))
			return
		}
		response = append(response, data...)
		if i != (len(rep.Children) - 1) {
			response = append(response, ',')
		}
	}
	response = append(response, ']')
	d.Reply(response)
}

func getContainer(d *DiscoAPI, p string) {
	name := getName(p)
	data, err := scanContainer(d, PREFIX+"/containers/"+name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(data)
}

func addContainer(d *DiscoAPI, p string, payload []byte) {
	name := getName(p)
	persist.Create(PREFIX+"/containers/"+name, string(payload), true)
	d.Reply([]byte("success"))
}

func removeContainer(d *DiscoAPI, p string) {
	name := getName(p)
	_, err := persist.Delete(PREFIX + "/containers/" + name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply([]byte("success"))
}

func scanContainer(d *DiscoAPI, path string) ([]byte, error) {
	var data []byte
	rep, err := persist.Read(path)
	if err != nil {
		return data, err
	}
	data = []byte(rep.Value)
	return data, nil
}

func getName(path string) string {
	pathArr := strings.Split(path, "/")
	return pathArr[len(pathArr)-1]
}
