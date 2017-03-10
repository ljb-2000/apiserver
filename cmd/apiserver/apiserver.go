package apiserver

import (
	"flag"

	"apiserver/pkg/util/config"
)

var (
	kubeconfig = flag.String("kubeconfig", "config", "--kubeconfig=kubeconfig file path")
	miniconfig = flag.String("config", "config.json", "--config=config file path")
)

func main() {
	flag.Parse()
	config.Parse(miniconfig)
}
