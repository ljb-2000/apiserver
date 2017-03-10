package k8sclient

import (
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

func initK8Sclient(kubeaddr, kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags(kubeaddr, kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
