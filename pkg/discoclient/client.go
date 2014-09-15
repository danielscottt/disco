package discoclient

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (c *Client) do(payload string) ([]byte, error) {

	var buff []byte

	conn, err := net.Dial("unix", c.Path)
	if err != nil {
		return buff, err
	}
	defer conn.Close()
	conn.Write([]byte(payload))
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
	id, err := c.do("/disco/local/node_id")
	if err != nil {
		return string(id), err
	}
	return string(id), nil
}

func (c *Client) RegisterContainer(ct *Container) ([]byte, error) {

	id, err := c.GetNodeId()
	if err != nil {
		return []byte(id), err
	}

	cJson, err := ct.Marshal()
	if err != nil {
		return cJson, err
	}

	reply, err := c.do(fmt.Sprintf("/disco/api/add_container\n%s", cJson))
	if err != nil {
		return reply, err
	}
	if string(reply) != "success" {
		return reply, errors.New(string(reply))
	}

	return reply, nil
}

func (c *Client) RemoveContainer(conId string) ([]byte, error) {

	id, err := c.GetNodeId()
	if err != nil {
		return []byte(id), err
	}

	reply, err := c.do(fmt.Sprintf("/disco/api/remove_container\n%s", conId))
	if err != nil {
		return reply, err
	}
	if string(reply) != "success" {
		return reply, errors.New(string(reply))
	}

	return reply, nil
}

type Container struct {
	Host  string
	Ports []dockerclient.Port
	Id    string
	Names []string
}

func NewContainer(names []string, id string, ports []dockerclient.Port) *Container {
	c := &Container{
		Names: names,
		Id:    id,
		Ports: ports,
	}
	c.Host, _ = os.Hostname()
	return c
}

func (c *Container) Marshal() ([]byte, error) {
	cJson, err := json.Marshal(c)
	if err != nil {
		return cJson, err
	}
	return cJson, nil
}
