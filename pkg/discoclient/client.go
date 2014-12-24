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
