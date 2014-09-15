package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.google.com/p/go-uuid/uuid"

	"github.com/danielscottt/disco/pkg/api"
	"github.com/danielscottt/disco/pkg/poller"
)

func main() {

	nodeId := uuid.New()

	api, err := discoapi.NewDiscoAPI(nodeId)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	// Close socket on SIGKILL && SIGTERM
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("%s Signal receieved: Exiting Disco Daemon", sig)
		api.Stop()
		os.Exit(0)
	}(sigc)

	go api.Start()
	defer api.Stop()

	poller.Start(nodeId)

}
