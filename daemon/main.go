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

	p "github.com/danielscottt/disco/pkg/persist"
)

var (
	persist *p.Controller
	config  *DaemonConfig
)

func createTree() {
	for _, p := range []string{"nodes", "containers", "links"} {
		_, err := persist.Create("/disco/data/"+p, "", 0)
		if err != nil {
			log.Fatalf(err)
		}
	}
}

func registerNode(nodeId string) {
	addrs, _ := net.InterfaceAddrs()
	addrsStrings := make([]string, 0)
	for _, a := range addrs {
		if match, _ := regexp.MatchString("::", a.String()); !match {
			addrsStrings = append(addrsStrings, a.String())
		}
	}
	_, err := persist.Create("/disco/data/nodes/"+nodeId, []byte(strings.Join(addrsStrings, ",")), 0)
	if err != nil {
		log.Fatalf(err)
	}
}

func ensureCleanShutdown() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("%s Signal receieved: Exiting Disco Daemon", sig)
		api.Stop()
		// TODO
		// os.Remove(discoDataPath + "/nodes/" + nodeId)
		os.Exit(0)
	}(sigc)
}

func main() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf(err)
	}

	o := &p.ControllerOptions{
		Nodes: config.Persist.Nodes,
		Type:  config.Persist.Type,
	}

	persist, err = p.NewController(o)

	nodeId := uuid.New()

	createTree()

	registerNode(nodeId)

	api, err := NewDiscoAPI(nodeId)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	ensureCleanShutdown()

	go api.Start()
	defer api.Stop()

	StartPoller(nodeId)
}
