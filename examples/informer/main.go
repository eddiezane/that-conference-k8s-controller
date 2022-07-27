package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		panic(err)
	}

	client := kubernetes.NewForConfigOrDie(config)

	ctx := signals.SetupSignalHandler()
	watch, err := client.CoreV1().Pods("default").Watch(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("starting")

  for running := true; running; {
		select {
		case e := <-watch.ResultChan():
			fmt.Println(e)
		case <-ctx.Done():
			fmt.Println("stopping")
      running = false
		}
	}

	fmt.Println("done")
}
