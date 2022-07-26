package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
  configFlags := genericclioptions.NewConfigFlags(true)
  discoveryClient, err := configFlags.ToDiscoveryClient()
  if err != nil {
    panic(err)
  }

  mapper, err := configFlags.ToRESTMapper()
  if err != nil {
    panic(err)
  }

  resources, err := discoveryClient.ServerPreferredResources()
  if err != nil {
    panic(err)
  }

  for _, r := range resources {
    for _, rr := range r.APIResources {
      gv, err := schema.ParseGroupVersion(r.GroupVersion)
      if err != nil {
        panic(err)
      }
      gvr := gv.WithResource(rr.Name)
      fmt.Println(gvr)

      gvk, err := mapper.KindFor(gvr)
      if err != nil {
        panic(err)
      }

      fmt.Println(gvk)
    }
  }
}
