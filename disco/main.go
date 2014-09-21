package main

import (
	"fmt"
	"strings"

	"github.com/danielscottt/commando"
	"github.com/danielscottt/disco/pkg/discoclient"
	"github.com/danielscottt/disco/pkg/dockerclient"
)

var (
	link   *commando.Command
	c      *discoclient.Client
	docker *dockerclient.Client
)

func main() {
	c = discoclient.NewClient("/var/run/disco.sock")
	docker = dockerclient.NewClient("/var/run/docker.sock")

	disco := &commando.Command{
		Name:        "disco",
		Description: "A Container Network Discovery tool",
	}

	link = &commando.Command{
		Name:        "link",
		Description: "Link containers together",
		Execute:     linkContainers,
	}
	link.AddOption("targets", "The target container(s) [NAME=container:port", true, "-t", "--target")
	link.AddOption("image", "Image to create linked container from", true, "-i", "--image")
	disco.AddSubCommand(link)

	nodeId := &commando.Command{
		Name:        "node-id",
		Description: "Get Disco Node Id",
		Execute: func() {
			id, err := c.GetNodeId()
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println(id)
		},
	}
	disco.AddSubCommand(nodeId)

	list := &commando.Command{
		Name:        "list",
		Description: "List Disco-managed Containers",
		Execute: func() {
			cons, err := c.GetContainers()
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
		},
	}
	disco.AddSubCommand(list)

	disco.Parse()
}
