package discoclient

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/danielscottt/disco/pkg/disco"
)

func (c *Client) LinkContainers(name string, source, target *disco.Container) (*disco.Link, error) {
	l := &disco.Link{
		Id:     uuid.New(),
		Source: source,
		Target: target,
		Name:   name,
	}
	if l.Source.Links == nil {
		l.Source.Links = make(map[string][]string)
	}
	l.Source.Links["source"] = append(l.Source.Links["source"], l.Id)
	if l.Target.Links == nil {
		l.Target.Links = make(map[string][]string)
	}
	l.Target.Links["target"] = append(l.Source.Links["target"], l.Id)
	lj, err := l.Marshal()
	if err != nil {
		return l, err
	}
	reply, err := c.do("/disco/link" + "\n" + string(lj))
	if err != nil {
		return l, err
	}
	if string(reply) != "success" {
		return l, err
	}
	return l, nil
}
