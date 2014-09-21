package main

import (
	"fmt"
	"strconv"
	"strings"
)

func linkContainers() {
	switch val := link.Options["targets"].Value.(type) {
	case string:
		name, container, port := parseTarget(val)
		// determine location / IP
		target, err := c.GetContainer(container)
		if err != nil {
			fmt.Println(err)
			return
		}
		if c.GetNodeId() != target.HostNode {
			// handle disparate node
		} else {
			docker.CreateContainer()
		}
		// add link to disco
		// create ENV array
		// create container
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
