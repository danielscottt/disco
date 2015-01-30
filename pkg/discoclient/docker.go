package discoclient

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/danielscottt/disco/pkg/disco"
)

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

func (c *Client) CreateDockerContainer(con *disco.Container) (*disco.Container, error) {
	var updated disco.Container
	cj, err := con.Marshal()
	reply, err := c.do("/disco/api/docker/create/" + con.Name + "\n" + string(cj))
	if err != nil {
		return &updated, err
	}
	replySplit := strings.Split(string(reply), "\n")
	if replySplit[0] != "success" {
		return &updated, errors.New(string(reply))
	}
	err = json.Unmarshal([]byte(replySplit[1]), &updated)
	if err != nil {
		return &updated, err
	}
	return &updated, nil
}
