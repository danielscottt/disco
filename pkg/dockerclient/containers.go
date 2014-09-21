package dockerclient

import "encoding/json"

type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

type Container struct {
	Command string
	Created int
	Id      string
	Image   string
	Names   []string
	Ports   []Port
}

func (c *Client) GetContainers() ([]Container, error) {

	var containers []Container

	body, err := c.do("GET", "/containers/json")
	if err != nil {
		return containers, err
	}

	err = json.Unmarshal(body, &containers)
	if err != nil {
		return containers, err
	}

	return containers, nil
}

type containerCreate struct {
"Hostname":"",
"Domainname": "",
"User":"",
"Memory":0,
"MemorySwap":0,
"CpuShares": 512,
"Cpuset": "0,1",
 "AttachStdin":false,
"AttachStdout":true,
"AttachStderr":true,
"PortSpecs":null,
"Tty":false,
"OpenStdin":false,
"StdinOnce":false,
"Env":null,
"Cmd":[
"date"
],
"Image":"base",
"Volumes":{
"/tmp": {}
},
"WorkingDir":"",
"NetworkDisabled": false,
"ExposedPorts":{
"22/tcp": {}
},
"RestartPolicy": { "Name": "always" }
}
