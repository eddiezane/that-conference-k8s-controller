package main

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"

	picturesv1 "github.com/eddiezane/that-conference-k8s-controller/pkg/apis/pictures/v1"
)

func init() {
  picturesv1.AddToScheme(scheme.Scheme)
}

func main() {
  configFlags := genericclioptions.NewConfigFlags(true)

  mapper, err := configFlags.ToRESTMapper()
  if err != nil {
    panic(err)
  }

  resourceMapping, err := mapper.RESTMapping(picturesv1.Kind("Picture"))
  if err != nil {
    panic(err)
  }

  object, err := scheme.Scheme.New(resourceMapping.GroupVersionKind)
  if err != nil {
    panic(err)
  }

  picture := object.(*picturesv1.Picture)

  picture.Name = "testing"

  fmt.Println(picture)
}
