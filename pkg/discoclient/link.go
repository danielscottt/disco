package discoclient

import (
	"github.com/danielscottt/disco/pkg/disco"
)

func (c *Client) LinkContainers(name string, source, target disco.Container) (*disco.Link, error) {
	l := &disco.Link{
		Id:     uuid.New(),
		Source: source,
		Target: target,
		Name:   name,
	}
	reply, err := c.do("/disco/links/" + l.Id + "\n" + l.Marshal())
	if err != nil {
		return l, err
	}
	if reply != "success" {
		return l, err
	}
	return l, nil
}
