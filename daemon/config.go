package main

import "code.google.com/p/gcfg"

type DaemonConfig struct {
	Persist struct {
		Nodes []string
		Type  string
	}
	Disco struct {
		DiscoSocket  string
		DockerSocket string
	}
}

func LoadConfig() error {
	err := gcfg.ReadFileInto(&config, "/etc/disco/disco.conf")
	if err != nil {
		return err
	}
	return nil
}
