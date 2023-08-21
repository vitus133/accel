package main

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func clusterVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "config.openshift.io",
		Version:  "v1",
		Resource: "clusterversions",
	}
}
func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// Dynamic client - for applying and monitoring arbitrary manifests.
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	//get clusterversion
	cv, err := dynamicClient.Resource(clusterVersionResource()).Get(ctx, "version", v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	stat := cv.Object["status"].(map[string]interface{})

	data, err := json.Marshal(stat["conditions"])
	if err != nil {
		panic(err.Error())
	}
	var cond []v1.Condition
	err = json.Unmarshal(data, &cond)
	if err != nil {
		panic(err.Error())
	}
	if meta.IsStatusConditionPresentAndEqual(cond, "Progressing", v1.ConditionFalse) {
		fmt.Println("not progressing")
	}
}
