package dockerclient

import "encoding/json"

type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

type Container struct {
	Command string
	Created int
	Id      string
	Image   string
	Names   []string
	Ports   []Port
}

func (c *Client) GetContainers() ([]Container, error) {

	var containers []Container

	body, err := c.do("GET", "/containers/json")
	if err != nil {
		return containers, err
	}

	err = json.Unmarshal(body, &containers)
	if err != nil {
		return containers, err
	}

	return containers, nil
}
