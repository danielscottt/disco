package persist

import "errors"

type Controller interface {
	Create(string, string, bool) (*Reply, error)
	Delete(string) (*Reply, error)
	Read(string) (*Reply, error)
}

type ControllerOptions struct {
	Type  string
	Nodes []string
}

type Reply struct {
	Location    string
	Value       string
	TransID     uint64
	HasChildren bool
	Children    []string
}

func NewController(options ControllerOptions) (Controller, error) {
	var c Controller
	var err error
	if options.Type == "etcd" {
		c, err = NewEtcdController(options)
		if err != nil {
			return c, err
		}
		return c, nil
	} else {
		return c, errors.New("Unknown Controller type: " + options.Type)
	}
}
