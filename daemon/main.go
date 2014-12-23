package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"code.google.com/p/go-uuid/uuid"

	"github.com/danielscottt/disco/pkg/discoclient"
	p "github.com/danielscottt/disco/pkg/persist"
)

const (
	PREFIX = "/disco/data"
)

var (
	nodeId  string
	persist p.Controller
	dc      *discoclient.Client
	api     *DiscoAPI
)

func createTree() {
	for _, p := range []string{"nodes", "containers", "links"} {
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

func registerNode() {
	addrs, _ := net.InterfaceAddrs()
	addrsStrings := make([]string, 0)
	for _, a := range addrs {
		if match, _ := regexp.MatchString("::", a.String()); !match {
			addrsStrings = append(addrsStrings, a.String())
		}
	}
	_, err := persist.Create(PREFIX+"/nodes/"+nodeId, strings.Join(addrsStrings, ","), false)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func init() {
	nodeId = uuid.New()
	err := LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	o := &p.ControllerOptions{
		Nodes: config.Persist.Nodes,
		Type:  config.Persist.Type,
	}
	persist, err = p.NewController(o)
	if err != nil {
		log.Fatalf(err.Error())
	}
	createTree()
	registerNode()
	api, err = NewDiscoAPI(&ApiConfig{
		Id:        nodeId,
		DockerUri: config.Disco.DockerSocket,
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
		persist.Delete(PREFIX + "/nodes/" + nodeId)
		os.Exit(0)
	}(sigc)
}

func main() {
	go api.Start()
	defer api.Stop()
	StartPoller()
}
