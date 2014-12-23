package discoclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/danielscottt/disco/pkg/disco"
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

func (c *Client) RegisterContainer(con *disco.Container) error {

	cJson, err := con.Marshal()
	if err != nil {
		return err
	}

	reply, err := c.do(fmt.Sprintf("/disco/api/add_container/%s\n%s", con.Name, cJson))
	if err != nil {
		return err
	}
	if string(reply) != "success" {
		return errors.New(string(reply))
	}

	return nil
}

func (c *Client) RemoveContainer(name string) error {

	reply, err := c.do("/disco/api/remove_container/" + name)
	if err != nil {
		return err
	}
	if string(reply) != "success" {
		return errors.New(string(reply))
	}

	return nil
}

func (c *Client) GetContainer(name string) (*disco.Container, error) {
	var con disco.Container
	reply, err := c.do(fmt.Sprintf("/disco/api/get_container/%s", name))
	if err != nil {
		return &con, err
	}
	if err := json.Unmarshal(reply, &con); err != nil {
		return &con, err
	}
	return &con, nil
}

func (c *Client) GetContainers() ([]disco.Container, error) {
	var cons []disco.Container
	reply, err := c.do("/disco/api/get_containers")
	if err != nil {
		return cons, err
	}
	if err := json.Unmarshal(reply, &cons); err != nil {
		return cons, err
	}
	return cons, nil
}

func (c *Client) CollectDockerContainers() ([]disco.Container, error) {
	var cons []disco.Container
	reply, err := c.do("/disco/api/docker/collect")
	if err != nil {
		return cons, err
	}
	if err := json.Unmarshal(reply, &cons); err != nil {
		return cons, err
	}
	return cons, nil
}
