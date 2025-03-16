package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/automatedtomato/grpc-auth-service/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Analyze command line parameter
	address := flag.String("address", "localhost:50051", "gRPC server address")
	useTLS := flag.Bool("tls", false, "Use TLS")
	certFile := flag.String("cert", "certs/server.crt", "TLS certificate file")
	flag.Parse()

	// Configure connection setting
	var opts []grpc.DialOption
	if *useTLS {
		creds, err := credentials.NewClientTLSFromFile(*certFile, "")
		if err != nil {
			log.Fatalf("Failed to load credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to gRPC server
	conn, err := grpc.Dial(*address, opts...)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	client := proto.NewAuthServiceClient(conn)

	// Create context
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute test process
	testAuthService(cxt, client)
}

func testAuthService(ctx context.Context, client proto.AuthServiceClient) {
	// User registration test
	registerResp, err := client.Register(ctx, &proto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Failed to register: %v", err)
	}
	log.Printf("Register response: %v", registerResp)

	// Login test
	loginResp, err := client.Login(ctx, &proto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	log.Printf("Login response: %v", loginResp)

	// User info acquisition test
	if loginResp.Success {
		userInfoResp, err := client.GetUserInfo(ctx, &proto.UserInfoRequest{
			SessionToken: loginResp.SessionToken,
		})
		if err != nil {
			log.Fatalf("Failed to get user info: %v", err)
		}
		log.Printf("UserInfo response: %v", userInfoResp)
	}

	// Reset password request test
	resetReqResp, err := client.RequestPasswordReset(ctx, &proto.PasswordResetRequest{
		Email: "test@example.com",
	})
	if err != nil {
		log.Fatalf("Reset password request failed: %v", err)
	}
	log.Printf("Password reset request response: %v", resetReqResp)

	// Reset password test
	if resetReqResp.Success {
		resetResp, err := client.ResetPassword(ctx, &proto.NewPasswordRequest{
			ResetToken:  resetReqResp.SessionToken,
			NewPassword: "newpassword456",
		})
		if err != nil {
			log.Fatalf("Reset password failed: %v", err)
		}
		log.Printf("Reset password response: %v", resetResp)
	}

	// Login test with new password
	newLoginResp, err := client.Login(ctx, &proto.LoginRequest{
		Username: "testuser",
		Password: "newpassword456",
	})
	if err != nil {
		log.Fatalf("Failed to login with new password: %v", err)
	}
	log.Printf("New login response: %v", newLoginResp)
}
