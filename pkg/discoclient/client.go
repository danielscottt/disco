package discoclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/fsouza/go-dockerclient"
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

func (c *Client) RegisterContainer(con *docker.APIContainers) error {

	id, err := c.GetNodeId()
	if err != nil {
		return err
	}

	ct := &Container{
		HostNode: id,
		Name:     (*con).Names[0][1:],
		Id:       (*con).ID,
		Ports:    (*con).Ports,
	}

	cJson, err := ct.Marshal()
	if err != nil {
		return err
	}

	reply, err := c.do(fmt.Sprintf("/disco/api/add_container/%s\n%s", ct.Name, cJson))
	if err != nil {
		return err
	}
	if string(reply) != "success" {
		return errors.New(string(reply))
	}

	return nil
}

func (c *Client) RemoveContainer(name string) error {

	reply, err := c.do(fmt.Sprintf("/disco/api/remove_container\n%s", name))
	if err != nil {
		return err
	}
	if string(reply) != "success" {
		return errors.New(string(reply))
	}

	return nil
}

func (c *Client) GetContainer(name string) (*Container, error) {
	var con Container
	reply, err := c.do(fmt.Sprintf("/disco/api/get_container/%s", name))
	if err != nil {
		return &con, err
	}
	if err := json.Unmarshal(reply, &con); err != nil {
		return &con, err
	}
	return &con, nil
}

func (c *Client) GetContainers() ([]Container, error) {
	var cons []Container
	reply, err := c.do("/disco/api/get_containers")
	if err != nil {
		return cons, err
	}
	if err := json.Unmarshal(reply, &cons); err != nil {
		return cons, err
	}
	return cons, nil
}

type Container struct {
	Name     string
	HostNode string
	Ports    []docker.APIPort
	Id       string
	Links    []Link
}

func (c *Container) Marshal() ([]byte, error) {
	cJson, err := json.Marshal(c)
	if err != nil {
		return cJson, err
	}
	return cJson, nil
}

type Link struct {
	Id      string
	Source  *Container
	Target  *Container
	PortMap map[string]Port
}

type Port struct {
	Name    string
	Private int
	Public  int
}
