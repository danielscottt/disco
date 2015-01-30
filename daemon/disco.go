package main

import (
	"encoding/json"
	"strings"

	"github.com/danielscottt/disco/pkg/disco"
)

func (d *DiscoAPI) getContainers() {
	response := []byte{'['}
	rep, err := d.Persist.Read(d.DataPath + "/containers/nodes/" + d.NodeId)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	for i, f := range rep.Children {
		data, err := d.scanContainer(f)
		if err != nil {
			d.Reply([]byte(err.Error()))
			return
		}
		response = append(response, data...)
		if i != (len(rep.Children) - 1) {
			response = append(response, ',')
		}
	}
	response = append(response, ']')
	d.Reply(response)
}

func (d *DiscoAPI) getContainer(p string) {
	name := getName(p)
	data, err := d.scanContainer(PREFIX + "/containers/nodes/" + d.NodeId + "/" + name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(data)
}

func (d *DiscoAPI) addContainer(p string, payload []byte) {
	name := getName(p)
	d.Persist.Create(PREFIX+"/containers/nodes/"+d.NodeId+"/"+name, string(payload), true)
	d.Persist.Create(PREFIX+"/containers/master/"+name, string(payload), true)
	d.Reply([]byte("success"))
}

func (d *DiscoAPI) removeContainer(p string) {
	err := d.clearOutContainer(p)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply([]byte("success"))
}

func (d *DiscoAPI) scanContainer(path string) ([]byte, error) {
	var data []byte
	rep, err := d.Persist.Read(path)
	if err != nil {
		return data, err
	}
	data = []byte(rep.Value)
	return data, nil
}

func (d *DiscoAPI) clearOutContainer(p string) error {
	name := getName(p)
	_, err := d.Persist.Delete(PREFIX+"/containers/nodes/"+d.NodeId+"/"+name, false)
	_, err = d.Persist.Delete(PREFIX+"/containers/master/"+name, false)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiscoAPI) createLink(payload []byte) {
	var l disco.Link
	err := json.Unmarshal(payload, &l)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	_, err = d.Persist.Create("/disco/links/"+l.Id, string(payload), false)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply([]byte("success"))
}

func getName(path string) string {
	pathArr := strings.Split(path, "/")
	return pathArr[len(pathArr)-1]
}
