package discoclient

import (
	"encoding/json"
	"net"
	"os"

	"github.com/danielscottt/disco/pkg/dockerclient"
)

type Client struct {
	Path string
}

func NewClient(path string) *Client {
	return &Client{
		Path: path,
	}
}

func (c *Client) do(path string) ([]byte, error) {

	var buff []byte

	conn, err := net.Dial("unix", c.Path)
	if err != nil {
		return buff, err
	}
	defer conn.Close()
	conn.Write([]byte(path))
	if err != nil {
		return buff, err
	}
	buff = make([]byte, 1024)
	size, err := conn.Read(buff[:])
	if err != nil {
		return buff, err
	}
	return buff[0:size], nil
}

func (c *Client) GetNodeId() (string, error) {
	id, err := c.do("/disco/local/node_dfsfid")
	if err != nil {
		return string(id), err
	}
	return string(id), nil
}

func (c *Client) RegisterContainer() {
	c.do("/disco/api/add_container\nhello world")
}

type RegisteredContainer struct {
	Host  string
	Ports []dockerclient.Port
	Id    string
	Names []string
}

func NewRegisteredContainer(names []string, id string, ports []dockerclient.Port) *RegisteredContainer {
	c := &RegisteredContainer{
		Names: names,
		Id:    id,
		Ports: ports,
	}
	c.Host, _ = os.Hostname()
	return c
}

func (r *RegisteredContainer) Marshal() ([]byte, error) {
	cJson, err := json.Marshal(r)
	if err != nil {
		return cJson, err
	}
	return cJson, nil
}
