package driver

import "google.golang.org/grpc"

const (
	DefaultName = "bsos.rishabhsvats.dev"
)

type Driver struct {
	name     string
	region   string
	endpoint string

	srv grpc.Server
}
type InputParams struct {
	Name     string
	Endpoint string
	Token    string
	Region   string
}

func NewDriver(params InputParams) *Driver {
	return &Driver{}
}
