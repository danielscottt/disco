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

func poll(nodeId string) {

	dataDir := os.Getenv("DISCO_DATA_PATH")
	if dataDir == "" {
		log.Fatalf("DISCO_DATA_PATH is blank, cannot continue")
	}

	ls, err := ioutil.ReadDir(dataDir)
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

		cJson, err := marshalContainer(c)
		if err != nil {
			log.Print("Error marshalling container [", c, "], skipping")
			continue
		}

		filePath := fmt.Sprintf("%s/%s:%s", dataDir, nodeId, c.Id)

		if _, present := (*lMap)[c.Id]; !present {
			addContainer(&c, &filePath, &cJson)
		} else {
			updateContainer(&c, &filePath, &cJson)
		}

	}

	removeStaleContainers(lMap, cMap, &dataDir)
}

func addContainer(c *dockerclient.Container, filePath *string, cJson *[]byte) {
	log.Print("New container [", (*c).Id, "] discovered")
	ioutil.WriteFile(*filePath, *cJson, 644)
}

func updateContainer(c *dockerclient.Container, filePath *string, cJson *[]byte) {
	file, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Println(err)
		return
	}
	fileHash := md5.Sum(file)
	portsHash := md5.Sum(*cJson)
	if fileHash != portsHash {
		log.Print("Container [", (*c).Id, "] exists but has been updated")
		ioutil.WriteFile(*filePath, *cJson, 644)
	}
}

func removeStaleContainers(lMap *map[string]os.FileInfo, cMap *map[string]dockerclient.Container, dataDir *string) {
	for _, l := range *lMap {
		if _, present := (*cMap)[l.Name()]; !present {
			log.Print("Container [", l.Name(), "] has been removed")
			os.Remove(fmt.Sprintf("%s/%s", *dataDir, l.Name()))
		}
	}
}

func marshalContainer(c dockerclient.Container) ([]byte, error) {
	cd := discoclient.NewRegisteredContainer(c.Names, c.Id, c.Ports)
	cJson, _ := cd.Marshal()
	return cJson, nil
}

func getContainers() (*map[string]dockerclient.Container, error) {

	cMap := make(map[string]dockerclient.Container)

	client, err := dockerclient.NewClient(os.Getenv("DOCKER_API_PATH"))
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
		(*mapPointer)[c.Id] = c
	}
}

func mapList(ls []os.FileInfo) *map[string]os.FileInfo {
	lMap := make(map[string]os.FileInfo)
	for _, l := range ls {
		lMap[l.Name()] = l
	}
	return &lMap
}
func Start(d time.Duration, nodeId string) {
	for {
		poll(nodeId)
		time.Sleep(d)
	}
}
