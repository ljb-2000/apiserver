package container

import (
	"fmt"

	"errors"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"minipaas/pkg/util/jsonx"
)

//
type Port struct {
	ExposePort int
}

func (port *Port) String() string {
	return fmt.Sprintf("[ExposePort:%d]", port.ExposePort)
}

//ContainerDto is a DTO of Container
type ContainerDto struct {
	Id      string
	Image   string
	AppName string
	Ports   []*Port
	Status  string
}

//Container is struct of docker container
type Container struct {
	Id         int                `json:"id" xorm:"pk not null varchar(12)"`
	ExtranetIp string             `json:"extranetIp" xorm:"varchar(256)"`
	IntranetIp string             `json:"intranetIp" xorm:"varchar(256)"`
	Image      string             `json:"image" xorm:"varchar(1024)"`
	Name       string             `json:"name" xorm:"varchar(1024)"`
	Ports      []v1.ContainerPort `json:"ports" xorm:"varchar(1024)"`
	Mounts     []v1.VolumeMount   `json:"mounts" xorm:"varchar(1024)"`
	Status     int                `json:"status" xorm:"int(1)"`
	Envs       []v1.EnvVar        `json:"envs" xorm:"varchar(1024)"`
	Created    string             `json:"created" xorm:"varchar(1024)"`
	Cpu        string             `json:"cpu" xorm:"varchar(1024)"`
	Memory     string             `json:"memory" xorm:"varchar(1024)"`
}

func (container *Container) String() string {
	containerStr, err := jsonx.ToJson(container)
	if err != nil {
		log.Errorf("container to string err :%s", err.Error())
		return ""
	}
	return containerStr
}

func (container *Container) Insert() error {
	_, err := engine.Insert(container)

	if err != nil {
		return err
	}

	return nil
}

func (container *Container) Delete() error {
	_, err := engine.Id(container.Id).Delete(container)

	if err != nil {
		return err
	}

	return nil
}

func (container *Container) Update() error {
	_, err := engine.Id(container.Id).Update(container)
	if err != nil {
		return err
	}

	return nil
}

func (container *Container) QueryOne() (*Container, error) {
	has, err := engine.Id(container.Id).Get(container)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("the query data not exist")
	}

	return container, nil
}

func (container *Container) QuerySet() ([]*Container, error) {
	containerSet := []*Container{}
	err := engine.Where("1 and 1 order by ip desc").Find(&containerSet)

	if err != nil {
		return nil, err
	}

	return containerSet, nil
}
