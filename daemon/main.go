package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielscottt/disco/pkg/disco"
	"github.com/danielscottt/disco/pkg/discoclient"
	p "github.com/danielscottt/disco/pkg/persist"
)

const (
	PREFIX = "/disco/data"
)

var (
	node *disco.Node
	dc   *discoclient.Client
	api  *DiscoAPI
)

func createTree(persist p.Controller) {
	for _, p := range []string{"nodes", "containers/nodes", "containers/master", "links"} {
		exists, err := persist.Exists(PREFIX + "/" + p)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !exists {
			_, err := persist.CreatePath(PREFIX + "/" + p)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func registerNode(persist p.Controller) {
	node = disco.NewNode()
	nj, err := node.Marshal()
	if err != nil {
		log.Fatalf(err.Error())
	}
	_, err = persist.Create(PREFIX+"/nodes/"+node.Id, string(nj), false)
	_, err = persist.CreatePath(PREFIX + "/containers/nodes/" + node.Id)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func init() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	o := &p.ControllerOptions{
		Nodes: config.Persist.Nodes,
		Type:  config.Persist.Type,
	}
	persist, err := p.NewController(o)
	if err != nil {
		log.Fatalf(err.Error())
	}
	createTree(persist)
	registerNode(persist)
	api, err = NewDiscoAPI(&ApiConfig{
		Id:        node.Id,
		DockerUri: config.Disco.DockerSocket,
		Persist:   persist,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	dc = discoclient.NewClient("/var/run/disco.sock")
	if err != nil {
		log.Fatalf(err.Error())
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("%s Signal receieved: Exiting Disco Daemon", sig)
		api.Stop()
		os.Exit(0)
	}(sigc)
}

func main() {
	go api.Start()
	defer api.Stop()
	StartPoller()
}
