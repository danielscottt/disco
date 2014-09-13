package discocli

import (
	"fmt"

	"github.com/danielscottt/commando"
	"github.com/danielscottt/disco/pkg/discoclient"
)

var link *commando.Command

func linkContainers() {
	//target := discoclient.GetContainer(link.Options["target"].Value)
	//source := discoclient.GetContainer(link.Options["source"].Value)
}

func Parse() {

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
	link.AddOption("container", "The container to start", true, "-c", "--container")
	disco.AddSubCommand(link)

	nodeId := &commando.Command{
		Name:        "node-id",
		Description: "Get Disco Node Id",
		Execute: func() {
			c := discoclient.NewClient("/var/run/disco.sock")
			id, err := c.GetNodeId()
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println(id)
		},
	}
	disco.AddSubCommand(nodeId)

	disco.Parse()
}
