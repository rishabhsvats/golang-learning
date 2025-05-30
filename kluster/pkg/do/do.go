package do

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Create(c kubernetes.Interface, sec string) error {
	token, err := getToken(c, sec)
	if err != nil {
		return err
	}
	client := godo.NewFromToken(token)
	fmt.Println(client)
	return nil
}
func getToken(client kubernetes.Interface, sec string) (string, error) {
	namespace := strings.Split(sec, "/")[0]
	name := strings.Split(sec, "/")[1]
	s, err := client.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(s.Data["token"]), nil
}
