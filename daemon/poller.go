package main

import (
	"crypto/md5"
	"log"
	"os"
	"time"

	"github.com/danielscottt/disco/pkg/discoclient"
	"github.com/danielscottt/disco/pkg/dockerclient"
)

func poll(nodeId, discoPath, dockerPath string) {

	dClient := discoclient.NewClient(discoPath)

	ls, err := dClient.GetContainers()
	if err != nil {
		log.Println(err)
		return
	}
	lMap := mapList(ls)

	cMap, err := getContainers(dockerPath)
	if err != nil {
		log.Println(err)
		return
	}

	for _, c := range *cMap {

		name := c.Names[0][1:]

		if _, present := (*lMap)[name]; !present {
			log.Print("New container [", name, "] discovered")
			if err := dClient.RegisterContainer(&c); err != nil {
				log.Print("Error: ", err)
				continue
			}
		} else {
			updateContainer(dClient, &c, (*lMap)[name], nodeId)
		}

	}

	removeStaleContainers(lMap, cMap, dClient)
}

func updateContainer(dClient *discoclient.Client, c *dockerclient.Container, dc discoclient.Container, nodeId string) {
	current := &discoclient.Container{
		HostNode: nodeId,
		Name:     (*c).Names[0][1:],
		Id:       (*c).Id,
		Ports:    (*c).Ports,
	}
	currentJson, err := current.Marshal()
	existingJson, err := dc.Marshal()
	if err != nil {
		log.Print(err)
		return
	}
	existingHash := md5.Sum(existingJson)
	currentHash := md5.Sum(currentJson)
	if existingHash != currentHash {
		log.Print("Container [", dc.Name, "] exists but has been updated")
		if err := (*dClient).RegisterContainer(c); err != nil {
			log.Print("Error: ", err)
		}
	}
}

func removeStaleContainers(lMap *map[string]discoclient.Container, cMap *map[string]dockerclient.Container, dClient *discoclient.Client) {
	for _, l := range *lMap {
		if _, present := (*cMap)[l.Name]; !present {
			log.Print("Removing Container [", l.Name, "]")
			if err := (*dClient).RemoveContainer(l.Name); err != nil {
				log.Print("Error removing container [", l.Name, "]")
			}
		}
	}
}

func getContainers(dockerPath string) (*map[string]dockerclient.Container, error) {

	cMap := make(map[string]dockerclient.Container)

	client, err := dockerclient.NewClient(dockerPath)
	if err != nil {
		log.Println(err)
		return &cMap, nil
	}

	cs, err := client.GetContainers()
	if err != nil {
		log.Println(err)
		return &cMap, nil
	}

	mapContainers(&cMap, cs)
	return &cMap, nil
}

func mapContainers(mapPointer *map[string]dockerclient.Container, cs []dockerclient.Container) {
	for _, c := range cs {
		(*mapPointer)[c.Names[0][1:]] = c
	}
}

func mapList(ls []discoclient.Container) *map[string]discoclient.Container {
	lMap := make(map[string]discoclient.Container)
	for _, l := range ls {
		lMap[l.Name] = l
	}
	return &lMap
}

func Start(nodeId, discoPath string) {

	var dockerPath string
	if os.Getenv("DOCKER_API_PATH") != "" {
		dockerPath = os.Getenv("DOCKER_API_PATH")
	} else {
		log.Fatalf("Docker api path cannot be blank. Cannot start")
	}

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

	log.Println("START Poller:", dur, "loop time")

	for {
		poll(nodeId, discoPath, dockerPath)
		time.Sleep(duration)
	}
}
