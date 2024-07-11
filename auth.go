package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/brotherlogic/auth/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/auth/proto"
)

var (
	port         = flag.Int("port", 8080, "Server port for grpc traffic")
	metricsPort  = flag.Int("metrics_port", 8081, "Metrics port")
	internalPort = flag.Int("internal_port", 8082, "Port to serve internal grpc traffic")
)

func main() {
	s := server.NewServer()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%v", *metricsPort), nil)
		log.Fatalf("gramophile is unable to serve metrics: %v", err)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("gramophile is unable to listen on the grpc port %v: %v", *port, err)
	}
	gsInternal := grpc.NewServer()
	pb.RegisterAuthServiceServer(gsInternal, s)
	go func() {
		if err := gsInternal.Serve(lis); err != nil {
			log.Fatalf("gramophile is unable to serve internal grpc: %v", err)
		}
	}()

	lis2, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("gramophile is unable to listen on the grpc port %v: %v", *port, err)
	}
	gs := grpc.NewServer()
	pb.RegisterAuthServiceServer(gs, s)
	err = gs.Serve(lis2)
	if err != nil {
		log.Fatalf("Bad serve: %v", err)
	}
}
