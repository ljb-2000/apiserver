package config

import (
	"encoding/json"
	"io/ioutil"

	"apiserver/pkg/util/logger"
)

type Config struct {
	Driver     string `json:"driver"`
	Dsn        string `json:"dsn"`
	Server     string `json:"server"`
	K8sServer  string `json:"k8sserver"`
	Kubeconfig string `json:"kubeconfig"`
}

var (
	log         = logger.New("")
	GloabConfig = &Config{}
)

func Parse(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("read config file fail, the reason is %v ", err)
	}
	err = json.Unmarshal(data, GloabConfig)
	if err != nil {
		log.Fatalf("unmarshal config data to config struct fail, the reason is %v ", err)
	}
}
