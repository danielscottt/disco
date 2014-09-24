package main

import (
	"crypto/md5"
	"log"
	"time"

	"github.com/danielscottt/disco/pkg/discoclient"
	"github.com/fsouza/go-dockerclient"
)

func poll() {

	dClient := discoclient.NewClient(config.Disco.DiscoSocket)

	ls, err := dClient.GetContainers()
	if err != nil {
		log.Println(err)
		return
	}
	lMap := mapList(ls)

	cMap, err := getContainers()
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
			updateContainer(dClient, &c, (*lMap)[name])
		}

	}

	removeStaleContainers(lMap, cMap, dClient)
}

func updateContainer(dClient *discoclient.Client, c *docker.APIContainers, dc discoclient.Container) {
	current := &discoclient.Container{
		HostNode: nodeId,
		Name:     (*c).Names[0][1:],
		Id:       (*c).ID,
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

func removeStaleContainers(lMap *map[string]discoclient.Container, cMap *map[string]docker.APIContainers, dClient *discoclient.Client) {
	for _, l := range *lMap {
		if _, present := (*cMap)[l.Name]; !present {
			log.Print("Removing Container [", l.Name, "]")
			if err := (*dClient).RemoveContainer(l.Name); err != nil {
				log.Print("Error removing container [", l.Name, "]: ", err.Error())
			}
		}
	}
}

func getContainers() (*map[string]docker.APIContainers, error) {

	cMap := make(map[string]docker.APIContainers)

	client, err := docker.NewClient(config.Disco.DockerSocket)
	if err != nil {
		log.Println(err)
		return &cMap, nil
	}

	cs, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Println(err)
		return &cMap, nil
	}

	mapContainers(&cMap, cs)
	return &cMap, nil
}

func mapContainers(mapPointer *map[string]docker.APIContainers, cs []docker.APIContainers) {
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
