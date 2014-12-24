package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"

	"github.com/fsouza/go-dockerclient"

	p "github.com/danielscottt/disco/pkg/persist"
)

type DiscoAPI struct {
	NodeId   string
	DataPath string
	Persist  p.Controller

	docker     *docker.Client
	listener   net.Listener
	connection net.Conn
	stop       chan bool
}

type ApiConfig struct {
	Id        string
	DockerUri string
	Persist   p.Controller
}

func NewDiscoAPI(config *ApiConfig) (*DiscoAPI, error) {

	var d *DiscoAPI

	log.Print("Disco socket Path: [/var/dun/disco.sock]")

	l, err := net.Listen("unix", "/var/run/disco.sock")
	if err != nil {
		return d, err
	}

	d = &DiscoAPI{
		listener: l,
		NodeId:   config.Id,
		DataPath: PREFIX,
		Persist:  config.Persist,
	}
	d.docker, err = docker.NewClient(config.DockerUri)
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
	d.Persist.Delete(PREFIX+"/nodes/"+node.Id, false)
	resp, _ := d.Persist.Read(PREFIX + "/containers/nodes/" + node.Id)
	for _, c := range resp.Children {
		d.clearOutContainer(c)
	}
	d.Persist.Delete(PREFIX+"/containers/nodes/"+node.Id, true)
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
		d.getContainers()
	case p == "/disco/api/docker/collect":
		d.collectDockerContainers()
	case addCont.MatchString(p):
		d.addContainer(p, payload)
	case rmCont.MatchString(p):
		d.removeContainer(p)
	case getCont.MatchString(p):
		d.getContainer(p)
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
