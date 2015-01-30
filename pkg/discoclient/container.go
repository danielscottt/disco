package discoclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/danielscottt/disco/pkg/disco"
)

func (c *Client) ContainerExists(name string) (bool, error) {
	reply, err := c.do(fmt.Sprintf("/disco/api/get_container/%s", name))
	if err != nil {
		return false, err
	}
	if strings.Contains(string(reply), "100: Key not found") {
		return false, nil
	}
	return true, nil
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
