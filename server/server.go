package server

import "context"

type Server interface {
	Start() error
	Shutdown(context.Context) error
}
