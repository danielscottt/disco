package dockerclient

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
)

type Client struct {
	Path string

	clientconn *httputil.ClientConn
}

func NewClient(path string) (*Client, error) {
	c := &Client{
		Path: path,
	}
	if c.Path == "" {
		return c, errors.New("Docker API Path cannot be blank")
	}
	return c, nil
}

func (c *Client) do(method string, path string) ([]byte, error) {

	var body []byte

	err := c.setConnection()
	defer c.clientconn.Close()
	if err != nil {
		return body, err
	}

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return body, err
	}
	resp, err := c.clientconn.Do(req)
	if err != nil {
		return body, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func (c *Client) setConnection() error {
	sock, err := net.Dial("unix", c.Path)
	if err != nil {
		return err
	}
	c.clientconn = httputil.NewClientConn(sock, nil)
	return nil
}
