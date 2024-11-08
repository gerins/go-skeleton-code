package cmd

import (
	"fmt"
	"net"

	"go-skeleton-code/pkg/log"

	"google.golang.org/grpc"

	"go-skeleton-code/config"
)

type GRPCServer struct {
	cfg    *config.Config
	Server *grpc.Server
}

// NewGRPCServer returns new GRPCServer.
func NewGRPCServer(cfg *config.Config) *GRPCServer {
	return &GRPCServer{
		cfg:    cfg,
		Server: grpc.NewServer(),
	}
}

func (s *GRPCServer) Run() chan bool {
	// Run GRPC Server
	grpcAddress := fmt.Sprintf("%s:%s", s.cfg.App.GRPC.Host, s.cfg.App.GRPC.Port)
	listen, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed init grpc tcp server, %v", err)
	}

	go func() {
		log.Infof("grpc server running at %v", grpcAddress)
		if err := s.Server.Serve(listen); err != nil {
			log.Fatalf("failed start grpc server, %v", err)
		}
	}()

	grpcServerExitSignal := make(chan bool)
	go func() {
		<-grpcServerExitSignal

		log.Info("stopping grpc server")
		s.Server.GracefulStop()
		log.Info("finished stopping grpc server")

		grpcServerExitSignal <- true // Send signal already finish the job
	}()

	return grpcServerExitSignal
}
