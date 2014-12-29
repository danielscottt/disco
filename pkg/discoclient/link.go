package discoclient

import (
	"github.com/danielscottt/disco/pkg/disco"
)

func (c *Client) LinkContainers(name string, source, target disco.Container) (
	*disco.Container, *disco.Link, error) {
	l := &disco.Link{
		Id:      uuid.New(),
		Source:  source,
		Target:  target,
		EnvName: name,
	}
}
