package driver

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

const (
	DefaultName = "bsos.rishabhsvats.dev"
)

type Driver struct {
	name     string
	region   string
	endpoint string

	srv *grpc.Server
}
type InputParams struct {
	Name     string
	Endpoint string
	Token    string
	Region   string
}

func NewDriver(params InputParams) *Driver {
	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
		region:   params.Region,
	}
}

// starting the gRPC server as per the CSI spec
func (drv *Driver) Run() error {

	url, err := url.Parse(drv.endpoint)
	if err != nil {
		return fmt.Errorf("error while parsing the endpoint %s\n", err.Error())
	}

	if url.Scheme != "unix" {
		return fmt.Errorf("only supported scheme  is unix, but provided %s\n", url.Scheme)
	}
	grpcAddress := path.Join(url.Host, filepath.FromSlash(url.Path))
	if url.Host == "" {
		grpcAddress = filepath.FromSlash(url.Path)
	}

	listener, err := net.Listen(url.Scheme, grpcAddress)
	if err != nil {
		return fmt.Errorf("listen failed %s\n", err.Error())
	}
	fmt.Println(listener)
	drv.srv = grpc.NewServer()
	csi.RegisterNodeServer(drv.srv, drv)

	return nil
}
