package poller

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/danielscottt/disco/pkg/discoclient"
	"github.com/danielscottt/disco/pkg/dockerclient"
)

func poll(nodeId, dataPath, discoPath, dockerPath string) {

	dClient := discoclient.NewClient(discoPath)

	ls, err := ioutil.ReadDir(dataPath)
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

		filePath := fmt.Sprintf("%s/%s", dataPath, name)

		if _, present := (*lMap)[name]; !present {
			log.Print("New container [", name, "] discovered")
			if err := dClient.RegisterContainer(&c); err != nil {
				log.Print("Error: ", err)
				continue
			}
		} else {
			updateContainer(dClient, &c, filePath, nodeId)
		}

	}

	removeStaleContainers(lMap, cMap, dClient)
}

func updateContainer(dClient *discoclient.Client, c *dockerclient.Container, filePath, nodeId string) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	cd := &discoclient.Container{
		HostNode: nodeId,
		Name:     (*c).Names[0][1:],
		Id:       (*c).Id,
		Ports:    (*c).Ports,
	}
	cJson, err := cd.Marshal()
	if err != nil {
		log.Print(err)
		return
	}
	fileHash := md5.Sum(file)
	currentHash := md5.Sum(cJson)
	if fileHash != currentHash {
		log.Print("Container [", cd.Name, "] exists but has been updated")
		if err := (*dClient).RegisterContainer(c); err != nil {
			log.Print("Error: ", err)
		}
	}
}

func removeStaleContainers(lMap *map[string]os.FileInfo, cMap *map[string]dockerclient.Container, dClient *discoclient.Client) {
	for _, l := range *lMap {
		if _, present := (*cMap)[l.Name()]; !present {
			log.Print("Removing Container [", l.Name(), "]")
			if err := (*dClient).RemoveContainer(l.Name()); err != nil {
				log.Print("Error removing container [", l.Name(), "]")
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

func mapList(ls []os.FileInfo) *map[string]os.FileInfo {
	lMap := make(map[string]os.FileInfo)
	for _, l := range ls {
		lMap[l.Name()] = l
	}
	return &lMap
}
func Start(nodeId, dataPath, discoPath string) {

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
		poll(nodeId, dataPath, discoPath, dockerPath)
		time.Sleep(duration)
	}
}
