package main

import (
	"log"
	"time"

	"github.com/danielscottt/disco/pkg/disco"
)

func poll() {
	lMap, err := getDiscoContainers()
	cMap, err := collectDockerContainers()
	if err != nil {
		log.Println("Error collecting Docker Containers", err)
		return
	}
	for _, c := range *cMap {
		if _, present := (*lMap)[c.Name]; !present {
			log.Print("New container [", c.Name, "] discovered")
			if err := dc.RegisterContainer(c); err != nil {
				log.Print("Error: ", err)
				continue
			}
		} else {
			updateContainer(c, (*lMap)[c.Name])
		}
	}
	removeStaleContainers(lMap, cMap)
}

func updateContainer(fromDocker, fromDisco *disco.Container) {
	if fromDisco.HasLinks() {
		fromDocker.Links = fromDisco.Links
	}
	dockerHash, err := fromDocker.Hash()
	discoHash, err := fromDisco.Hash()
	if err != nil {
		log.Print("Update error:" + err.Error())
		return
	}
	if dockerHash != discoHash {
		log.Print("Container [", fromDisco.Name, "] exists but has been updated")
		if err := dc.RegisterContainer(fromDocker); err != nil {
			log.Print("Error: ", err)
		}
	}
}

func removeStaleContainers(lMap *map[string]*disco.Container, cMap *map[string]*disco.Container) {
	for _, l := range *lMap {
		if _, present := (*cMap)[l.Name]; !present {
			log.Print("Removing Container [", l.Name, "]")
			if err := dc.RemoveContainer(l.Name); err != nil {
				log.Print("Error removing container [", l.Name, "]: ", err.Error())
			}
		}
	}
}

func getDiscoContainers() (*map[string]*disco.Container, error) {
	lMap := make(map[string]*disco.Container)
	ls, err := dc.GetContainers()
	if err != nil {
		log.Println("Error retrieving Disco Containers:", err)
		return &lMap, err
	}
	mapContainers(&lMap, ls)
	return &lMap, nil
}

func collectDockerContainers() (*map[string]*disco.Container, error) {
	cMap := make(map[string]*disco.Container)
	cs, err := dc.CollectDockerContainers()
	if err != nil {
		log.Println(err)
		return &cMap, err
	}
	mapContainers(&cMap, cs)
	return &cMap, nil
}

func mapContainers(mapPointer *map[string]*disco.Container, cs []disco.Container) {
	for _, c := range cs {
		(*mapPointer)[c.Name] = &c
	}
}

func StartPoller() {
	var dur string
	if config.Disco.LoopTime != "" {
		dur = config.Disco.LoopTime
	} else {
		dur = "2s"
	}
	duration, err := time.ParseDuration(dur)
	if err != nil {
		log.Fatalf("Invalid Loop Time given")
	}
	log.Println("START Poller:", dur, "loop time")
	for {
		poll()
		time.Sleep(duration)
	}
}
