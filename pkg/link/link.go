package link

import (
	"github.com/danielscottt/disco/pkg/discoclient"
)

type Link struct {
	Id      string
	Source  *discoclient.Container
	Target  *discoclient.Container
	PortMap map[string]Port
}
