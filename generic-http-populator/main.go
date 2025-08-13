package main

import (
	"flag"
	"fmt"
	"log"

	populatorMachinery "github.com/kubernetes-csi/lib-volume-populator/populator-machinery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	prefix    = "k8s.rishabhsvats.dev"
	mountPath = "/mnt/vol"
)

func main() {
	var image string
	var namespace string
	var mode string
	var uri string

	flag.StringVar(&image, "image", "", "Image for populator component")
	flag.StringVar(&namespace, "namespace", "", "Namespace for populator component")
	flag.StringVar(&mode, "mode", "", "Mode to run the application in")
	flag.StringVar(&uri, "uri", "", "URI for the content of the volume")

	flag.Parse()

	switch mode {
	case "controller":
		gk := schema.GroupKind{
			Group: "k8s.rishabhsvats.dev",
			Kind:  "GenericHTTPPopulator",
		}
		gvr := schema.GroupVersionResource{
			Group:    gk.Group,
			Version:  "v1alpha1",
			Resource: "generichttppopulators",
		}
		populatorMachinery.RunController("", "", image, "", "", namespace, prefix, gk, gvr, mountPath, "", populatorArgs)
	case "populator":
		// the code that we write here is going to get called from populator pod
		// we know that PVC is mounted at `mountPath`
		populate(uri)
	default:
		log.Printf("Mode %s is not supported", mode)
	}
}
func populate(uri string) {

}

type GenericHTTPPopulator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec GenericHTTPPopulatorSpec `json:"spec"`
}
type GenericHTTPPopulatorSpec struct {
	URI string `json:"uri"`
}

// b here specifies if the volume is in block mode
// u is the populator instance
func populatorArgs(b bool, u *unstructured.Unstructured) ([]string, error) {
	var ghp GenericHTTPPopulator
	var args []string
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &ghp)
	if err != nil {
		log.Printf("Failed converting unstructured to GHP, error %s\n", err.Error())
		return args, err
	}
	args = append(args, "--mode=populator")
	args = append(args, fmt.Sprintf("--uri=%s", ghp.Spec.URI))
	return args, nil
}
