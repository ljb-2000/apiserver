package node

import (
	"errors"

	"apiserver/pkg/util/jsonx"
)

//Node node is server node
type Node struct {
	Ip          string  `json:"ip" xorm:"pk not null varchar(255)"`
	Hostname    string  `json:"hostname" xorm:"varchar(1024)"`
	Os          string  `json:"os" xorm:"varchar(1024)"`
	CpuNum      int     `json:"cpuNum" xorm:"int(11)"`
	CpuFree     float64 `json:"cpuFree" xorm:"int(11)"`
	CpuUse      float64 `json:"cpuUse" xorm:"int(11)"`
	MemoryTotal int64   `json:"memoryTotal" xorm:"int(11)"`
	MemoryFree  int64   `json:"memoryFree" xorm:"int(11)"`
	DiskTotal   int64   `json:"diskTotal" xorm:"int(11)"`
	DiskFree    int64   `json:"diskFree" xorm:"int(11)"`
}

func (node *Node) String() string {
	nodeStr, err := jsonx.ToJson(node)
	if err != nil {
		log.Errorf("node to string err :%s", err.Error())
		return ""
	}
	return nodeStr
}

func (node *Node) Insert() error {
	_, err := engine.Insert(node)

	if err != nil {
		return err
	}

	return nil
}

func (node *Node) Delete() error {
	_, err := engine.Id(node.Ip).Delete(node)

	if err != nil {
		return err
	}

	return nil
}

func (node *Node) Update() error {
	_, err := engine.Id(node.Ip).Update(node)

	if err != nil {
		return err
	}

	return nil
}

func (node *Node) QueryOne() (*Node, error) {
	has, err := engine.Id(node.Ip).Get(node)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("the query data not exist")
	}

	return node, nil
}

func (node *Node) QuerySet() ([]*Node, error) {
	nodeSet := []*Node{}
	err := engine.Where("1 and 1 order by ip desc").Find(&nodeSet)

	if err != nil {
		return nil, err
	}

	return nodeSet, nil
}
