package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"strings"

	klusterv1alpha1 "github.com/rishabhsvats/golang-learning/kluster/pkg/apis/rishabhsvats.dev/v1alpha1"

	"github.com/digitalocean/godo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ValidateKlusterVersion(kspec klusterv1alpha1.KlusterSpec) (bool, error) {
	client := initClient(kspec.TokenSecret)
	if client == nil {
		return false, errors.New("failed to create DO client")
	}
	options, _, err := client.Kubernetes.GetOptions(context.Background())
	if err != nil {
		return false, err
	}

	for _, version := range options.Versions {
		if kspec.Version == version.Slug {
			return true, nil
		}
	}
	return false, errors.New("The version is not supported")
}

func initClient(tokenSecret string) *godo.Client {
	token, err := getToken(tokenSecret)
	if err != nil {
		fmt.Printf("Error %s getttign token", err.Error())
		return nil
	}

	client := godo.NewFromToken(token)
	return client
}

func getToken(sec string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("error %s, getting inclusterconfig", err.Error())
	}
	// }
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		// handle error
		fmt.Printf("error %s, creating clientset\n", err.Error())
	}

	namespace := strings.Split(sec, "/")[0]
	name := strings.Split(sec, "/")[1]
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(s.Data["token"]), nil
}
