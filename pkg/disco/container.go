package disco

import (
	"crypto/md5"
	"encoding/json"
)

type Container struct {
	Name      string
	HostNode  string
	Id        string
	Links     map[string][]string
	Ports     []Port
	IPAddress string
	Env       []string
	Image     string
}

type Port struct {
	Private int
	Public  int
}

func (c *Container) Marshal() ([]byte, error) {
	cJson, err := json.Marshal(c)
	if err != nil {
		return cJson, err
	}
	return cJson, nil
}

func (c *Container) Hash() ([md5.Size]byte, error) {
	cj, err := c.Marshal()
	if err != nil {
		return [md5.Size]byte{}, err
	}
	return md5.Sum(cj), nil
}

func (c *Container) HasLinks() bool {
	if c.Links == nil {
		return false
	}
	return true
}
