package disco

type Link struct {
	Id     string
	Source *Container
	Target *Container
	Name   string
}
