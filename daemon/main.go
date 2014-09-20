package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.google.com/p/go-uuid/uuid"

	"github.com/danielscottt/disco/pkg/api"
)

func main() {

	nodeId := uuid.New()

	var discoSocketPath, discoDataPath string

	if os.Getenv("DISCO_SOCKET") != "" {
		discoSocketPath = os.Getenv("DISCO_SOCKET")
	} else {
		log.Fatalf("Disco socket path not set. Cannot start.")
	}

	if os.Getenv("DISCO_DATA_PATH") != "" {
		discoDataPath = os.Getenv("DISCO_DATA_PATH")
	} else {
		log.Fatalf("Disco data path not set. Cannot start.")
	}

	api, err := discoapi.NewDiscoAPI(nodeId, discoDataPath, discoSocketPath)
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

	Start(nodeId, discoSocketPath)

}
