package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Analyze command line parameters
	grpcAddr := flag.String("grpc-addr", "localhost:50051", "gRPC server address")
	webAddr := flag.String("web-addr", ":8080", "Web server address")
	useTLS := flag.Bool("tls", false, "Use TLS fot gRPC connection")
	certFiles := flag.String("cert", "../certs/server.crt", "TLS certification file")
	staticDir := flag.String("static", "../public", "Static file directory")
	flag.Parse()

	// Options for gRPC connection
	var opts []grpc.DialOption
	if *useTLS {
		creds, err := credentials.NewClientTLSFromFile(*certFiles, "")
		if err != nil {
			log.Fatalf("Failed to load credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to gRPC server
	conn, err := grpc.Dial(*grpcAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC wrapper
	grpcWebServer := grpcweb.WrapServer(
		grpc.NewServer(),
		grpcweb.WithOriginFunc(func(origin string) bool {
			// for developing purpose, accept all origins
			return true
		}),
	)

	// Static file handler
	fileServer := http.FileServer(http.Dir(*staticDir))

	// HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if grpcWebServer.IsGrpcWebRequest(r) || grpcWebServer.IsAcceptableGrpcCorsRequest(r) {
			grpcWebServer.ServeHTTP(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	// Initiate HTTP server
	log.Printf("Starting Web server on %s", *webAddr)
	if err := http.ListenAndServe(*webAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
