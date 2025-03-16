package server

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/automatedtomato/grpc-auth-service/api/proto"
	"github.com/automatedtomato/grpc-auth-service/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	server    *grpc.Server
	userStore storage.UserStore
}

func NewGRPCServer(useTLS bool, certFile, keyFile string) (*GRPCServer, error) {
	var opts []grpc.ServerOption

	if useTLS {
		// TLS configuration
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewServerTLSFromCert(&cert)
		opts = append(opts, grpc.Creds(creds))
	}

	// Create gRPC server
	server := grpc.NewServer(opts...)
	return &GRPCServer{
		server:    server,
		userStore: storage.NewInMemoryUserStore(),
	}, nil
}

func (s *GRPCServer) Start(address string) error {
	//  Register authentication service
	authServer := NewAuthServer(s.userStore)
	proto.RegisterAuthServiceServer(s.server, authServer)

	// Create listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("Starting gRPC server %s", address)
	return s.server.Serve(listener)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
