package main

import (
	"github.com/danielscottt/commando"

	"github.com/danielscottt/disco/pkg/discoclient"
)

var disco *discoclient.Client

func main() {
	disco = discoclient.NewClient("/var/run/disco.sock")

	root := &commando.Command{
		Name:        "disco",
		Description: "A Container Network Discovery tool",
	}

	link = &commando.Command{
		Name:        "link",
		Description: "Link containers together",
		Execute:     linkContainers,
	}
	link.AddOption("targets", "The target container(s) [NAME=container:port]", true, "-t", "--target")
	link.AddOption("image", "Image to create linked container from", true, "-i", "--image")
	link.AddOption("name", "Name to give linked container", true, "-n", "--name")
	root.AddSubCommand(link)

	nodeId := &commando.Command{
		Name:        "node-id",
		Description: "Get Disco Node Id",
		Execute:     getNodeId,
	}
	root.AddSubCommand(nodeId)

	list := &commando.Command{
		Name:        "list",
		Description: "List Disco-managed Containers",
		Execute:     listContainers,
	}
	root.AddSubCommand(list)

	root.Parse()
}
