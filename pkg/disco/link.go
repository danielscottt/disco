package disco

type Link struct {
	Id      string
	Source  *Container
	Target  *Container
	PortMap map[string]Port
}
