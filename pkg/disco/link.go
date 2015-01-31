package disco

import "encoding/json"

type Link struct {
	Id     string
	Source *Container
	Target *Container
	Name   string
}

func (l *Link) Marshal() ([]byte, error) {
	lj, err := json.Marshal(l)
	if err != nil {
		return lj, err
	}
	return lj, nil
}
