package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"code.google.com/p/go-uuid/uuid"
)

func registerNode(dataPath, nodeId string) {
	addrs, _ := net.InterfaceAddrs()
	addrsStrings := make([]string, 0)
	for _, a := range addrs {
		if match, _ := regexp.MatchString("::", a.String()); !match {
			addrsStrings = append(addrsStrings, a.String())
		}
	}
	nodeFilePath := dataPath + "/nodes/" + nodeId
	ioutil.WriteFile(nodeFilePath, []byte(strings.Join(addrsStrings, ",")), 644)
}

func createTree(dataPath string) {
	for _, p := range []string{"nodes", "containers", "links"} {
		os.MkdirAll(dataPath+"/"+p, 644)
	}
}

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
	createTree(discoDataPath)
	registerNode(discoDataPath, nodeId)

	api, err := NewDiscoAPI(nodeId, discoDataPath, discoSocketPath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}

	// Close socket on SIGINT && SIGTERM
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("%s Signal receieved: Exiting Disco Daemon", sig)
		api.Stop()
		os.Remove(discoDataPath + "/nodes/" + nodeId)
		os.Exit(0)
	}(sigc)

	go api.Start()
	defer api.Stop()

	StartPoller(nodeId, discoSocketPath)

}
