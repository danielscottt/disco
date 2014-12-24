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
			Image:    dcont.Image,
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

func (d *DiscoAPI) createDockerContainer(p string, payload []byte) {
	var con disco.Container
	err := json.Unmarshal(&c)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	create := docker.CreateContainerOptions{
		Name: c.Name,
	}
	config := &dockerclient.Config{
		Image:           con.Image,
		Env:             con.Env,
		NetworkDisabled: false,
		AttachStdin:     false,
		AttachStdout:    false,
		AttachStderr:    false,
	}
	create.Config = config
	c, err := docker.CreateContainer(create)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	hc := &dockerclient.HostConfig{
		NetworkMode:     "bridge",
		PublishAllPorts: true,
	}
	hc.LxcConf = make([]dockerclient.KeyValuePair, 0)
	hc.PortBindings = make(map[dockerclient.Port][]dockerclient.PortBinding)
	err = docker.StartContainer(c.ID, hc)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	con.Id = c.ID
	conj, err := con.Marshal()
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	reply := []byte("success\n")
	reply = append(reply, conj)
	d.Reply(reply)
}
