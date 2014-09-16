package discocli

import (
	"fmt"

	"github.com/danielscottt/commando"
	"github.com/danielscottt/disco/pkg/discoclient"
)

var link *commando.Command

func linkContainers() {
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
	link.AddOption("source", "The source container [the one that is being created]", true, "-s", "--source")
	link.AddOption("port", "The name and the port to map to linked container [in NAME:port format]", true, "-p", "--port")
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
