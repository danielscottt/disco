package main

import "strings"

func (d *DiscoAPI) getContainers() {
	response := []byte{'['}
	rep, err := persist.Read(d.DataPath + "/containers")
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
	data, err := d.scanContainer(PREFIX + "/containers/" + name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply(data)
}

func (d *DiscoAPI) addContainer(p string, payload []byte) {
	name := getName(p)
	persist.Create(PREFIX+"/containers/"+name, string(payload), true)
	d.Reply([]byte("success"))
}

func (d *DiscoAPI) removeContainer(p string) {
	name := getName(p)
	_, err := persist.Delete(PREFIX + "/containers/" + name)
	if err != nil {
		d.Reply([]byte(err.Error()))
		return
	}
	d.Reply([]byte("success"))
}

func (d *DiscoAPI) scanContainer(path string) ([]byte, error) {
	var data []byte
	rep, err := persist.Read(path)
	if err != nil {
		return data, err
	}
	data = []byte(rep.Value)
	return data, nil
}

func getName(path string) string {
	pathArr := strings.Split(path, "/")
	return pathArr[len(pathArr)-1]
}
