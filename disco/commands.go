package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danielscottt/commando"
	dockerclient "github.com/fsouza/go-dockerclient"
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
	// make links
	// create containers
	// wait for discovery
	// write links to /links, containers
	links := make([]*disco.Link, 0)
	switch val := link.Options["targets"].Value.(type) {
	case string:
		link, err := linkContainer(val)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("\rcreating link " + link.Name + " [" + link.Id + "]")
		links = append(links, link)
	case []string:
		for _, v := range val {
			link, err := linkContainer(val)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("\rcreating link [" + link.Id + "]")
			links = append(links, link)
		}
	}
	fmt.Printf("\rstarting container [" + source.Name + "]")
	s := disco.CreateDockerContainer(source)
	commando.PrintFields(false, 0, "SOURCE", "TARGET", "LINK ID")
	for _, link := range links {
		commando.PrintFields(false, 0, link.Source.Name, link.Target, link.Id)
	}
}

func parseTarget(input string) (string, string, int64) {
	split := strings.Split(input, "=")
	name := split[0]
	container := strings.Split(split[1], ":")[0]
	port, _ := strconv.ParseInt(strings.Split(split[1], ":")[1], 0, 0)
	return name, container, port
}

func linkContainer(v string) *disco.Link {
	var link *disco.Link
	name, container, port := parseTarget(v)
	target, err := disco.GetContainer(container)
	if err != nil {
		return link, err
	}
	sourceName := link.Options["name"].Value.(string)
	sourceImage := link.Options["image"].Value.(string)
	source := makeContainer(sourceName, sourceImage)
	link = disco.LinkContainers(source, target)
	if err != nil {
		return link, err
	}
	return link, nil
}

func makeContainer(name, image string) *disco.Container {
	c := &disco.Container{
		Name:  name,
		Image: image,
	}
}
