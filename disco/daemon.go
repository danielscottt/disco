package main

import (
	"log"
	"os"
	"time"

	"code.google.com/p/go-uuid/uuid"

	api "github.com/danielscottt/disco/pkg/api"
	"github.com/danielscottt/disco/pkg/poller"
)

func main() {

	nodeId := uuid.New()

	var dur string
	if os.Getenv("DISCO_LOOP_TIME") != "" {
		dur = os.Getenv("DISCO_LOOP_TIME")
	} else {
		dur = "2s"
	}
	duration, err := time.ParseDuration(dur)
	if err != nil {
		log.Fatalf("Invalid Loop Time given")
	}

	log.Println("Docker API Path:", os.Getenv("DOCKER_API_PATH"))
	log.Println("Data Path:", os.Getenv("DISCO_DATA_PATH"))
	log.Println("START:", dur, "loop time")

	go api.StartListener(nodeId)

	poller.Start(duration, nodeId)

}
