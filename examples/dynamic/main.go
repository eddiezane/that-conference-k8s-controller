package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	picturesv1 "github.com/eddiezane/that-conference-k8s-controller/pkg/apis/pictures/v1"
)

func main() {
	configFlags := genericclioptions.NewConfigFlags(true)

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		panic(err)
	}

	client := dynamic.NewForConfigOrDie(config)

	ul, err := client.Resource(schema.GroupVersionResource{Group: "kuberneddies.dev", Version: "v1", Resource: "pictures"}).
		List(context.Background(), metav1.ListOptions{})

  for _, u := range ul.Items {
    picture := new(picturesv1.Picture)
    err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, picture)
    if err != nil {
      panic(err)
    }

    fmt.Println(picture)
  }
}
