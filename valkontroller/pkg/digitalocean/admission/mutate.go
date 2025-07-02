package admission

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	klusterv1alpha1 "github.com/rishabhsvats/golang-learning/kluster/pkg/apis/rishabhsvats.dev/v1alpha1"
	kdo "github.com/rishabhsvats/golang-learning/valkontroller/pkg/digitalocean"
	"gomodules.xyz/jsonpatch/v2"
	admv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func ServeKlusterMutation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serverkluster validation was called")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responsewriters.InternalError(w, r, err)

		fmt.Printf("error %s, reading the body", err.Error())
	}

	//get group version kind
	gvk := admv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
	var admissionReview admv1beta1.AdmissionReview
	_, _, err = codecs.UniversalDeserializer().Decode(body, &gvk, &admissionReview)
	if err != nil {
		fmt.Printf("error %s, converting request body to admission review type", err.Error())
	}
	//get kluster spec from admissionreview object
	gvkKluster := klusterv1alpha1.SchemeGroupVersion.WithKind("Kluster")
	var kluster klusterv1alpha1.Kluster
	_, _, err = codecs.UniversalDeserializer().Decode(admissionReview.Request.Object.Raw, &gvkKluster, &kluster)
	if err != nil {
		fmt.Printf("error %s, while getting kluster type from admissionreview", err.Error())
	}
	// this is kluster
	// apiVersion: rishabhsvats.dev/v1alpha1
	// kind: Kluster
	// metadata:
	// name: kluster-1
	// spec:
	// name: kluster-1
	// region: "nyc1"
	// tokenSecret: "default/dosecret"
	// nodePools:
	// 	- count: 3
	// 	name: "dummy-nodepool"
	// 	size: "sizes-2vcpu-2gb"
	newKluster := kluster.DeepCopy()
	if kluster.Spec.Version == "" {
		newKluster.Spec.Version = kdo.LatestKubeVersion(newKluster.Spec)
	}

	// this is k
	// apiVersion: rishabhsvats.dev/v1alpha1
	// kind: Kluster
	// metadata:
	// name: kluster-1
	// spec:
	// name: kluster-1
	// region: "nyc1"
	// tokenSecret: "default/dosecret"
	// nodePools:
	// 	- count: 3
	// 	name: "dummy-nodepool"
	// 	size: "sizes-2vcpu-2gb"
	jsonKluster, err := json.Marshal(newKluster)
	if err != nil {
		fmt.Printf("error %s, converting new kluster resource to json", err.Error())
	}
	ops, err := jsonpatch.CreatePatch(admissionReview.Request.Object.Raw, jsonKluster)
	if err != nil {
		fmt.Printf("error %s, creating patch", err.Error())
	}

	patch, err := json.Marshal(ops)
	if err != nil {
		fmt.Printf("error %s converting operations to slice byte")
	}

	fmt.Printf("patch that we have is %s", patch)

	jsonPatchType := admv1beta1.PatchTypeJSONPatch
	response := admv1beta1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		PatchType: &jsonPatchType,
		Patch:     patch,
	}
	admissionReview.Response = &response

	fmt.Printf("response that we are trying to return is %+v\n", response)
	res, err := json.Marshal(admissionReview)
	if err != nil {
		fmt.Printf("error %s, while converting response to byte slice", err.Error())
	}

	_, err = w.Write(res)
	if err != nil {
		fmt.Printf("error %s, writing respnse to responsewriter", err.Error())
	}
}
