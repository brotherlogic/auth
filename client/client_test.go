package client

import (
	"testing"

	"google.golang.org/grpc"
)

func TestLoad(t *testing.T) {
	interceptor := &AuthInterceptor{}
	grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthIntercept))
}
