# gRPC Authentication Service

A simple authentication microservice implemented using gRPC and Protocol Buffers in Go. This project demonstrates how to build a secure authentication service with TLS encryption and basic authentication features.

## Features

- **gRPC API with Protocol Buffers**: Strongly-typed API definition using protobuf
- **Basic Authentication**: User registration, login, and session management
- **Password Reset Flow**: Complete password reset functionality
- **TLS Encryption**: Secure communication with TLS certificates
- **In-Memory Storage**: Simple storage implementation for user data
- **Web Interface**: Simple frontend using gRPC-Web for browser access
- **Concurrent Request Handling**: Leveraging Go's concurrency model

## Installation

### Requirements

- Go 1.16+
- Protocol Buffers compiler (`protoc`)
- Git

### Build and Install

```bash
# Clone the repository
git clone https://github.com/yourusername/auth-service.git
cd auth-service

# Install Protocol Buffers compiler (macOS)
brew install protobuf
# For other platforms, download from: https://github.com/protocolbuffers/protobuf/releases

# Install Go Protocol Buffers plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install gRPC-Web proxy
go install github.com/improbable-eng/grpc-web/go/grpcwebproxy@latest

# Generate Protocol Buffers code
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  api/proto/auth.proto

# Generate TLS certificates
cd certs
./generate_certs.sh  # Or run the OpenSSL commands manually
cd ..

# Install dependencies
go mod tidy
```

## Usage

### Running the Server

```bash
# Start the server without TLS
go run cmd/server/main.go

# Start the server with TLS
go run cmd/server/main.go --tls=true
```

### Running the CLI Client (for testing)

```bash
# Test without TLS
go run cmd/client/main.go

# Test with TLS
go run cmd/client/main.go --tls=true
```

### Running the Web Interface

```bash
# Start the gRPC-Web proxy
cd web/proxy
go run main.go

# Access the web interface
open http://localhost:8080
```

### Available API Methods

- `Register`: Create a new user account
- `Login`: Authenticate and obtain a session token
- `RequestPasswordReset`: Request a password reset token
- `ResetPassword`: Reset password using a token
- `GetUserInfo`: Retrieve user information using a session token

## Project Structure

```
auth-service/
│
├── api/
│   └── proto/
│       ├── auth.proto      # Protocol Buffers definition file
│       └── auth.pb.go      # Auto-generated Go code
│
├── cmd/
│   ├── server/
│   │   └── main.go         # Server entry point
│   └── client/
│       └── main.go         # CLI client for testing
│
├── internal/
│   ├── server/
│   │   ├── server.go       # gRPC server implementation
│   │   └── auth.go         # Authentication logic
│   ├── storage/
│   │   └── user_store.go   # User data storage
│   └── model/
│       └── user.go         # User model
│
├── web/
│   ├── public/
│   │   └── index.html      # Web client interface & client logic
│   └── proxy/
│       └── main.go         # gRPC-Web proxy
│
├── certs/                  # TLS certificates
│   ├── server.key          
│   └── server.crt
│
├── go.mod
└── go.sum
```

## Learning Objectives

This project was created to learn the following Golang concepts/features:

- **Protocol Buffers**: Defining message types and services using protobuf
- **gRPC Server/Client Implementation**: Building RPC services in Go
- **TLS Configuration**: Setting up secure communication with certificates
- **Subject Alternative Names (SANs)**: Modern TLS certificate requirements
- **Authentication Flow**: Implementing registration, login, and password reset
- **Concurrent Request Handling**: Using Go's goroutines and synchronization
- **Error Handling**: Proper error handling in distributed systems
- **Web Integration**: Connecting browser clients to gRPC services

## Possible Extensions

Ideas for extending this project:

- Database integration (PostgreSQL, MySQL, Redis)
- JWT authentication implementation
- More robust error handling and validation
- Unit and integration tests
- CI/CD pipeline setup
- Logging and monitoring
- Role-based access control
- OAuth/OpenID Connect integration
- Containerization with Docker
- Kubernetes deployment configuration
- Rate limiting and request throttling
- Advanced security features (MFA, brute force protection)

## License

MIT License

## Contributing

While this project was created for learning purposes, improvement suggestions and bug reports are welcome. Feel free to create an Issue or Pull Request.

## Author

Hikaru Tomizawa