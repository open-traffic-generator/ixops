package configs

import "fmt"

type Node struct {
	Host   string `yaml:"host" default:"localhost"`
	Port   uint   `yaml:"port" default:"22"`
	User   string `yaml:"user" default:"admin"`
	Master bool   `yaml:"master" default:"true"`
}

func (n *Node) DockerHost() string {
	return fmt.Sprintf("ssh://%s@%s:%d", n.User, n.Host, n.Port)
}
