package driver

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/digitalocean/godo"
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

	//storage ec2.Client
	storage godo.StorageService
	ready   bool
	// http server, health check
	// storage clients

	csi.UnimplementedNodeServer
	csi.UnimplementedControllerServer
	csi.UnimplementedIdentityServer
}
type InputParams struct {
	Name     string
	Endpoint string
	Token    string
	Region   string
}

func NewDriver(params InputParams) (*Driver, error) {
	if params.Token == "" {
		return nil, errors.New("token must be specified")
	}
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	panic("unable to load SDK config: " + err.Error())
	// }

	//ec2Client := ec2.NewFromConfig(cfg)
	client := godo.NewFromToken(params.Token)

	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
		region:   params.Region,
		storage:  client.Storage,
		//storage: ec2Client
	}, nil
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

	if err := os.Remove(grpcAddress); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing listen address failed %s\n", err.Error())
	}

	listener, err := net.Listen(url.Scheme, grpcAddress)
	if err != nil {
		return fmt.Errorf("listen failed %s\n", err.Error())
	}
	fmt.Println(listener)
	drv.srv = grpc.NewServer()

	csi.RegisterNodeServer(drv.srv, drv)
	csi.RegisterControllerServer(drv.srv, drv)
	csi.RegisterIdentityServer(drv.srv, drv)

	drv.ready = true
	return drv.srv.Serve(listener)
}
