package persist

import "github.com/coreos/go-etcd/etcd"

type EtcdController struct {
	Nodes []string

	client *etcd.Client
}

func NewEtcdController(options *ControllerOptions) (*EtcdController, error) {
	ec := &EtcdController{
		Nodes: options.Nodes,
	}
	client := etcd.NewClient(ec.Nodes)
	ec.client = client
	return ec, nil
}

func (e *EtcdController) Create(path string, value string, update bool) (*Reply, error) {
	var (
		r    *Reply
		resp *etcd.Response
		err  error
	)
	if update {
		resp, err = e.client.Set(path, value, 0)
		if err != nil {
			return r, err
		}
	} else {
		resp, err = e.client.Create(path, value, 0)
		if err != nil {
			return r, err
		}
	}
	r = &Reply{
		Location: resp.Node.Key,
		Value:    resp.Node.Value,
		TransID:  resp.Node.ModifiedIndex,
	}
	return r, nil
}

func (e *EtcdController) CreatePath(path string) (*Reply, error) {
	var r *Reply
	resp, err := e.client.CreateDir(path, 0)
	if err != nil {
		return r, err
	}
	r = &Reply{
		Location: resp.Node.Key,
		Value:    resp.Node.Value,
		TransID:  resp.Node.ModifiedIndex,
	}
	return r, nil
}

func (e *EtcdController) Read(path string) (*Reply, error) {
	var r *Reply
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
	if r.HasChildren {
		children := []string{}
		for _, c := range (*resp).Node.Nodes {
			children = append(children, c.Key)
		}
		r.Children = children
	}
	return r, nil
}

func (e *EtcdController) Delete(path string, recursive bool) (*Reply, error) {
	var r *Reply
	resp, err := e.client.Delete(path, recursive)
	if err != nil {
		return r, err
	}
	r = &Reply{
		Location: resp.Node.Key,
		Value:    resp.Node.Value,
		TransID:  resp.Node.ModifiedIndex,
	}
	return r, nil
}

func (e *EtcdController) Exists(path string) (bool, error) {
	_, err := e.client.Get(path, false, false)
	if err != nil {
		etcdError := err.(*etcd.EtcdError)
		if etcdError.ErrorCode == 100 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
