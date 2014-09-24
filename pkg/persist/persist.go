package persist

import "errors"

type Controller interface {
	Create() Reply
	Delete() Reply
	Read() Reply
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
}

func NewController(options ControllerOptions) (Controller, err) {
	var c Controller
	var err error
	if ControllerOptions.Type && ControllerOptions.Type == "etcd" {
		c, err = NewEtcdController(options)
		if err != nil {
			return c, err
		}
		return c, nil
	} else {
		return c, errors.New("Unknown Controller type: %s", options.Type)
	}
}
