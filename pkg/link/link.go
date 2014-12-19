package link

type Link struct {
	Id     string
	Source container
	Target container
	Ports  []Port
}

type Port struct {
}

type container interface {
}
