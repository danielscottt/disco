package main

import (
	"fmt"
	"strings"

	"github.com/danielscottt/commando"
	"github.com/danielscottt/disco/pkg/discoclient"
)

var (
	link *commando.Command
	c    *discoclient.Client
)

func linkContainers() {
	//target, err := c.GetContainer(link.Options["target"].Value.(string))
	//source, err := c.GetContainer(link.Options["source"].Value.(string))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}

func Parse() {

	c = discoclient.NewClient("/var/run/disco.sock")

	disco := &commando.Command{
		Name:        "disco",
		Description: "A Container Discovery tool",
	}

	link = &commando.Command{
		Name:        "link",
		Description: "Link containers together",
		Execute:     linkContainers,
	}
	link.AddOption("target", "The target container", true, "-t", "--target")
	link.AddOption("source", "The source container [the one that is being created]", true, "-s", "--source")
	link.AddOption("port", "The name and the port to map to linked container [in NAME:port format]", true, "-p", "--port")
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
