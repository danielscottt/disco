package disco

import (
	"encoding/json"
	"net"
	"regexp"

	"code.google.com/p/go-uuid/uuid"
)

type Node struct {
	Id    string
	Addrs map[string]string
}

func NewNode() *Node {
	return &Node{
		Id:    uuid.New(),
		Addrs: getInterfaces(),
	}
}

func (n *Node) Marshal() ([]byte, error) {
	nj, err := json.Marshal(n)
	if err != nil {
		return nj, err
	}
	return nj, nil
}

func getInterfaces() map[string]string {
	addrs, _ := net.InterfaceAddrs()
	addrsStrings := []string{}
	for _, a := range addrs {
		if match, _ := regexp.MatchString("::", a.String()); !match {
			addrsStrings = append(addrsStrings, a.String())
		}
	}
	ints, _ := net.Interfaces()
	intMap := make(map[string]string)
	for i, a := range ints {
		match1, _ := regexp.MatchString("veth", a.Name)
		match2, _ := regexp.MatchString("docker", a.Name)
		if !match1 && !match2 {
			intMap[a.Name] = addrsStrings[i]
		}
	}
	return intMap
}
