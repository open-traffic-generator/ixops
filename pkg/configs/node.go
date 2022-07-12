package configs

import "fmt"

type Node struct {
	Host   *string `yaml:"host"`
	Port   *int    `yaml:"port"`
	User   *string `yaml:"user"`
	Master *bool   `yaml:"master"`
}

func (n *Node) DockerHost() string {
	return fmt.Sprintf("ssh://%s@%s:%d", *n.User, *n.Host, *n.Port)
}

func (v *Node) SetDefaults() {
	SetDefaultString(&v.Host, "localhost")
	SetDefaultInt(&v.Port, 22)
	SetDefaultString(&v.User, "admin")
	SetDefaultBool(&v.Master, true)
}
