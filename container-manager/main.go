package main

import (
	"context"
	"fmt"
	"log"
	"os"

	containerd "github.com/containerd/containerd/v2/client"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	containerdSocket := homedir + "/" + ".containerd.sock"

	targetPlatform := "io.containerd.runc.v2"

	if err := startContainer(context.Background(), containerdSocket, targetPlatform); err != nil {
		log.Fatalf("Fatal error: %s", err)
	}

}

func startContainer(ctx context.Context, containerdSocket, targetPlatform string) error {
	client, err := containerd.New(containerdSocket, containerd.WithDefaultRuntime(targetPlatform))
	if err != nil {
		return fmt.Errorf("containerd error at New: %s", err)
	}
}
