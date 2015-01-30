package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fsouza/go-dockerclient"

	"github.com/danielscottt/disco/pkg/disco"
)

func (d *DiscoAPI) collectDockerContainers() {
	ds, err := d.docker.ListContainers(docker.ListContainersOptions{})
	cs := make([]*disco.Container, len(ds))
	var wg sync.WaitGroup
	for i, dcont := range ds {
		wg.Add(1)
		// parallelize container data collection
		go func(index int, dc docker.APIContainers) {
			c := &disco.Container{
				HostNode: node.Id,
				Name:     dcont.Names[0][1:],
				Id:       dcont.ID,
				Image:    dcont.Image,
			}
			defer wg.Done()
			inspect, err := d.docker.InspectContainer(c.Id)
			if err != nil {
				log.Printf("error inspecting container " + c.Name)
				return
			}
			c.Env = inspect.Config.Env
			c.IPAddress = inspect.NetworkSettings.IPAddress
			c.Ports = make([]disco.Port, len(dcont.Ports))
			for i, p := range dcont.Ports {
				c.Ports[i] = disco.Port{
					Private: int(p.PrivatePort),
					Public:  int(p.PublicPort),
				}
			}
			cs[i] = c
		}(i, dcont)
	}
	wg.Wait()
	csj, err := json.Marshal(cs)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(csj)
}

func (d *DiscoAPI) createDockerContainer(p string, payload []byte) {
	var con disco.Container
	err := json.Unmarshal(payload, &con)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	create := docker.CreateContainerOptions{
		Name: con.Name,
	}
	config := &docker.Config{
		Image:           con.Image,
		Env:             con.Env,
		NetworkDisabled: false,
		AttachStdin:     false,
		AttachStdout:    false,
		AttachStderr:    false,
	}
	create.Config = config
	c, err := d.docker.CreateContainer(create)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	hc := &docker.HostConfig{
		NetworkMode:     "bridge",
		PublishAllPorts: true,
	}
	hc.LxcConf = make([]docker.KeyValuePair, 0)
	hc.PortBindings = make(map[docker.Port][]docker.PortBinding)
	err = d.docker.StartContainer(c.ID, hc)
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
	reply = append(reply, conj...)
	d.Reply(reply)
}
