package persist

import "github.com/coreos/go-etcd/etcd"

type EtcdController struct {
	Nodes []string

	client *etcd.Client
}

func NewEtcdController(options ControllerOptions) (*EtcdController, error) {
	ec := &EtcdController{
		Nodes: options.Nodes,
	}
	client, err := etcd.NewClient(ec.Nodes)
	if err != nil {
		return ec, err
	}
	ec.client = client
	return ec, nil
}

func (e *EtcdController) Create(path string, value string, ttl int) (*Reply, error) {
	var r Reply
	resp, err := e.client.Create(path, value, ttl)
	if err != nil {
		return r, err
	}
	r = &Reply{
		Location:    resp.Node.Key,
		Value:       resp.Node.Value,
		TransID:     resp.Node.ModifiedIndex,
		HasChildren: resp.Node.Dir,
	}
	return r, nil
}

func (e *EtcdController) Read(path string) (*Reply, error) {
	var r Reply
	resp, err := e.client.Get(path, false, false)
	if err != nil {
		return r, err
	}
	r = &Reply{
		Location:    resp.Node.Key,
		Value:       resp.Node.Value,
		TransID:     resp.Node.ModifiedIndex,
		HasChildren: resp.Node.Dir,
	}
	return r, nil
}

func (e *EtcdController) Delete(path string) (*Reply, error) {
	var r Reply
	resp, err := e.client.Delete(path, false)
	if err != nil {
		return r, err
	}
	r = &Reply{
		Location:    resp.Node.Key,
		Value:       resp.Node.Value,
		TransID:     resp.Node.ModifiedIndex,
		HasChildren: resp.Node.Dir,
	}
	return r, nil
}
