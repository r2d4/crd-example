package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/golang/glog"
	r2d4client "github.com/r2d4/crd/pkg/client/clientset/versioned"
	"github.com/r2d4/crd/pkg/controller"
	"github.com/r2d4/crd/pkg/util"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
)

func main() {
	flag.Parse()
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %v", err)
	}

	r2d4Clientset, err := r2d4client.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %v", err)
	}

	apiextensionsclientset, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("getting normal clientset")
	}
	crd, err := util.CreateCustomResourceDefinition(apiextensionsclientset)
	if err != nil {
		glog.Fatalf("creating crd %s", err)
	}

	controller := controller.R2d4Controller{
		R2d4Client: r2d4Clientset.R2d4V1().RESTClient(),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			apiextensionsclientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
			os.Exit(0)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	controller.Run(ctx)
}
