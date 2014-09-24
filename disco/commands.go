package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danielscottt/commando"
)

var nodeId, list, link *commando.Command

func getNodeId() {
	id, err := disco.GetNodeId()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(id)
}

func listContainers() {
	cons, err := disco.GetContainers()
	if err != nil {
		fmt.Println(err)
		return
	}
	commando.PrintFields(false, 0, "NAME", "HOST NODE", "DOCKER ID", "PORTS", "LINKS")
	for _, con := range cons {
		var portMap []string
		for _, p := range con.Ports {
			portMap = append(portMap, fmt.Sprintf("%d:%d", p.PrivatePort, p.PublicPort))
		}
		portString := strings.Join(portMap, ", ")
		commando.PrintFields(false, 0, con.Name, con.HostNode, con.Id[:12], portString, "soon...")
	}
}

func linkContainers() {
	switch val := link.Options["targets"].Value.(type) {
	case string:
		name, container, port := parseTarget(val)
		target, err := disco.GetContainer(container)
		if err != nil {
			fmt.Println(err)
			return
		}
		id, err := disco.GetNodeId()
		if err != nil {
			fmt.Println(err)
			return
		}
		if id != target.HostNode {
			// handle disparate node
		} else {
			config := buildContainerConfig(target)
			docker.CreateContainer(config)
		}
	case []string:
	}
}

func parseTarget(input string) (string, string, int64) {
	split := strings.Split(input, "=")
	name := split[0]
	container := strings.Split(split[1], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(split[1], ":")[1], 0, 0)
	return name, container, port
}

func buildContainerConfig(target disco.Container) docker.CreateContainerOptions {
	create := &docker.CreateContainerOptions{
		Name: link.Options["name"].Value,
	}
	config := *docker.Config{
