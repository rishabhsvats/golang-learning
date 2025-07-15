package main

import (
	"flag"
	"fmt"

	"github.com/rishabhsvats/golang-learning/bsos/pkg/driver"
)

func main() {

	var (
		endpoint = flag.String("endpoint", "defaultValue", "Endpoing our gRPC server will be running at")
		token    = flag.String("token", "defaultValue", "token of the storage provider")
		region   = flag.String("region", "ams3", "region where the volumes are going to be provisioned")
	)
	flag.Parse()
	fmt.Println(*endpoint, *token, *region)

	// create a driver instance
	drv := driver.NewDriver(driver.InputParams{
		Name: driver.DefaultName,
		// unix:///var/lib/csi/sockets/csi.sock
		Endpoint: *endpoint,
		Region:   *region,
		Token:    *token,
	})
	//run on that driver instance, it would start the gRPC server

	if err := drv.Run(); err != nil {
		fmt.Printf("error %s, running the driver", err.Error())
	}
}
