package discoclient

func (c *Client) CollectDockerContainers() ([]disco.Container, error) {
	var cons []disco.Container
	reply, err := c.do("/disco/api/docker/collect")
	if err != nil {
		return cons, err
	}
	if err := json.Unmarshal(reply, &cons); err != nil {
		return cons, err
	}
	return cons, nil
}

func (c *Client) CreateDockerContainer(c *Container) (*Container, error) {
	var updated Container
	cj, err := c.Marshal()
	reply, err := c.do("/disco/api/docker/create/" + c.Name)
	if err != nil {
		return &updated, err
	}
	replySplit := strings.Split(string(reply), "\n")
	if string(replySplit[0]) != "success" {
		return &updated, errors.New(string(reply))
	}
	err := json.Unmarshal(&updated)
	if err != nil {
		return &updated, err
	}
	return &updated, nil
}
