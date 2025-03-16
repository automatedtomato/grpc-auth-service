package main

import (
	"flag"
	"log"

	"github.com/automatedtomato/grpc-auth-service/internal/server"
)

func main() {
	// Analyze command line parameters
	address := flag.String("address", ":50051", "gRPC server address")
	useTLS := flag.Bool("tls", false, "Use TLS")
	certFile := flag.String("cert", "certs/server.crt", "TLS certificate file")
	keyFile := flag.String("key", "certs/server.key", "TLS key file")
	flag.Parse()

	// Create server
	grpcServer, err := server.NewGRPCServer(*useTLS, *certFile, *keyFile)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	// Launch server
	if err := grpcServer.Start(*address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
