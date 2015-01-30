package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/danielscottt/commando"

	d "github.com/danielscottt/disco/pkg/disco"
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
			portMap = append(portMap, fmt.Sprintf("%d:%d", p.Private, p.Public))
		}
		portString := strings.Join(portMap, ", ")
		commando.PrintFields(false, 0, con.Name, con.HostNode, con.Id[:12], portString, "soon...")
	}
}

func linkContainers() {
	links := make([]*d.Link, 0)
	switch val := link.Options["targets"].Value.(type) {
	case string:
		link, err := linkContainer(val)
		if err != nil {
			fmt.Println("Linking error: " + err.Error())
			return
		}
		fmt.Printf("creating link " + link.Name + " [" + link.Id + "]\n")
		links = append(links, link)
	case []string:
		for _, v := range val {
			link, err := linkContainer(v)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Printf("creating link [" + link.Id + "]\n")
			links = append(links, link)
		}
	}
	// each link has the same source
	source := links[0].Source
	fmt.Printf("starting container [" + source.Name + "]\n")
	_, err := disco.CreateDockerContainer(source)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	i := 0
	var exists bool
	for !exists {
		fmt.Printf("\rwaiting for discovery%s", strings.Repeat(".", i))
		time.Sleep(500 * time.Millisecond)
		exists, err = disco.ContainerExists(source.Name)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		i++
		if i > 50 {
			fmt.Println("Error: timeout waiting for discovery")
			os.Exit(1)
		}
	}
	time.Sleep(time.Second * 3)
	fmt.Println("")
	for _, link := range links {
		fmt.Printf("\r\033[Kupdating %s", link.Source.Name)
		disco.RegisterContainer(link.Source)
		fmt.Printf("\r\033[Kupdating %s", link.Target.Name)
		disco.RegisterContainer(link.Target)
	}
	fmt.Println("")
	fmt.Println("done.")
	fmt.Println("Results:")
	commando.PrintFields(true, 2, "SOURCE", "TARGET", "LINK ID")
	for _, link := range links {
		commando.PrintFields(true, 2, link.Source.Name, link.Target.Name, link.Id)
	}
}

func parseTarget(input string) (string, string, int) {
	split := strings.Split(input, "=")
	name := split[0]
	container := strings.Split(split[1], ":")[0]
	port, _ := strconv.Atoi(strings.Split(split[1], ":")[1])
	return name, container, port
}

func linkContainer(v string) (*d.Link, error) {
	var l *d.Link
	name, container, port := parseTarget(v)
	target, err := disco.GetContainer(container)
	if err != nil {
		return l, errors.New("error retrieving container: " + err.Error())
	}
	sourceName := link.Options["name"].Value.(string)
	sourceImage := link.Options["image"].Value.(string)
	source := makeContainer(sourceName, sourceImage)
	for _, p := range target.Ports {
		if p.Private == port {
			source.Env = append(source.Env, name+"_PORT="+fmt.Sprintf("%d", p.Private))
			source.Env = append(source.Env, name+"_HOST="+target.IPAddress)
		}
	}
	l, err = disco.LinkContainers(name, source, target)
	if err != nil {
		return l, err
	}
	return l, nil
}

func makeContainer(name, image string) *d.Container {
	return &d.Container{
		Name:  name,
		Image: image,
	}
}
