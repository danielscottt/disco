package main

import (
	"encoding/json"

	"github.com/fsouza/go-dockerclient"

	"github.com/danielscottt/disco/pkg/disco"
)

func (d *DiscoAPI) collectDockerContainers() {
	ds, err := d.docker.ListContainers(docker.ListContainersOptions{})
	cs := make([]*disco.Container, len(ds))
	for i, dcont := range ds {
		c := &disco.Container{
			HostNode: node.Id,
			Name:     dcont.Names[0][1:],
			Id:       dcont.ID,
		}
		c.Ports = make([]disco.Port, len(dcont.Ports))
		for i, p := range dcont.Ports {
			c.Ports[i] = disco.Port{
				Private: int(p.PrivatePort),
				Public:  int(p.PublicPort),
			}
		}
		cs[i] = c
	}
	csj, err := json.Marshal(cs)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(csj)
}
