package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	klusterv1alpha1 "github.com/rishabhsvats/golang-learning/kluster/pkg/apis/rishabhsvats.dev/v1alpha1"
	"github.com/spf13/pflag"
	admv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/globalflag"
)

type Options struct {
	SecureServingOptions options.SecureServingOptions
}

func (o *Options) AddFlagSet(fs *pflag.FlagSet) {
	o.SecureServingOptions.AddFlags(fs)
}

type Config struct {
	SecureServingInfo *server.SecureServingInfo
}

func (o *Options) Config() *Config {
	if err := o.SecureServingOptions.MaybeDefaultWithSelfSignedCerts("0.0.0.0", nil, nil); err != nil {
		panic(err)
	}
	c := Config{}
	o.SecureServingOptions.ApplyTo(&c.SecureServingInfo)
	return &c
}

const (
	valKon = "val-kontroller"
)

func NewDefaultOptions() *Options {
	o := &Options{
		SecureServingOptions: *options.NewSecureServingOptions(),
	}
	o.SecureServingOptions.BindPort = 8443
	o.SecureServingOptions.ServerCert.PairName = valKon

	return o
}

func main() {
	options := NewDefaultOptions()
	fs := pflag.NewFlagSet(valKon, pflag.ExitOnError)
	globalflag.AddGlobalFlags(fs, valKon)

	options.AddFlagSet(fs)

	if err := fs.Parse(os.Args); err != nil {
		panic(err)
	}

	c := options.Config()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(ServeKlusterValidation))

	stopCh := server.SetupSignalHandler()
	_, ch, err := c.SecureServingInfo.Serve(mux, 30*time.Second, stopCh)
	if err != nil {
		panic(err)
	} else {
		<-ch
	}
}

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func ServeKlusterValidation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serverkluster validation was called")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
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
}
