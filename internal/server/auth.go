package server

import (
	"context"
	"time"

	"github.com/automatedtomato/grpc-auth-service/api/proto"
	"github.com/automatedtomato/grpc-auth-service/internal/model"
	"github.com/automatedtomato/grpc-auth-service/internal/storage"
)

// simple map to manage session token
type sessionManager struct {
	session map[string]string
}

func newSessionManager() *sessionManager {
	return &sessionManager{
		session: make(map[string]string),
	}
}

// Implementation of gRPC authentication service
type AuthServer struct {
	proto.UnimplementedAuthServiceServer
	userStore  storage.UserStore
	sessionMgr *sessionManager
}

func NewAuthServer(userStore storage.UserStore) *AuthServer {
	return &AuthServer{
		userStore:  userStore,
		sessionMgr: newSessionManager(),
	}
}

func (s *AuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &proto.RegisterResponse{
			Success: false,
			Message: "Username, email and password are required",
		}, nil
	}

	// Create user
	user, err := model.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return &proto.RegisterResponse{
			Success: false,
			Message: "Failed to create user: " + err.Error(),
		}, nil
	}

	// save user
	if err := s.userStore.Create(user); err != nil {
		return &proto.RegisterResponse{
			Success: false,
			Message: "Failed to register user:" + err.Error(),
		}, nil
	}

	return &proto.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  user.ID,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Search user by username
	user, err := s.userStore.GetByUsername(req.Username)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Validate password
	if !user.CheckPassword(req.Password) {
		return &proto.LoginResponse{
			Success: false,
			Message: "Invalid username or password",
		}, nil
	}

	// Generate session token
	token := model.RandomString(32)
	s.sessionMgr.session[token] = user.ID

	return &proto.LoginResponse{
		Success:      true,
		Message:      "Login successfully",
		SessionToken: token,
	}, nil
}

// Process password reset
func (s *AuthServer) RequestPasswordReset(ctx context.Context, req *proto.PasswordResetRequest) (*proto.PasswordResetResponse, error) {
	// Search user by email
	user, err := s.userStore.GetByEmail(req.Email)
	if err != nil {
		return &proto.PasswordResetResponse{
			Success: false,
			Message: "No account found with that email",
		}, nil
	}

	// Generate reset token
	resetToken := user.SetResetToken()

	if err := s.userStore.Update(user); err != nil {
		return &proto.PasswordResetResponse{
			Success: false,
			Message: "Failed to process reset request",
		}, nil
	}

	// NOTE: In a "real-world" application, you could implement email sending func etc.
	// In this project, for the simplification purpose, just return token
	return &proto.PasswordResetResponse{
		Success:      true,
		Message:      "Password reset link sent to your email",
		SessionToken: resetToken, // In an actual app, you don't do this!!
	}, nil
}

// Process password reset
func (s *AuthServer) ResetPassword(ctx context.Context, req *proto.NewPasswordRequest) (*proto.NewPasswordResponse, error) {

	// Find user by token
	user, err := s.userStore.GetByResetToken(req.ResetToken)
	if err != nil {
		return &proto.NewPasswordResponse{
			Success: false,
			Message: "Invalid or expired reset token",
		}, nil
	}

	// Check token is expired or not
	if user.ResetTokenExpires.Before(time.Now()) {
		return &proto.NewPasswordResponse{
			Success: false,
			Message: "Reset token has expired",
		}, nil
	}

	// Set new password
	newUser, err := model.NewUser(user.Username, user.Email, req.NewPassword)
	if err != nil {
		return &proto.NewPasswordResponse{
			Success: false,
			Message: "Failed to update password",
		}, nil
	}

	user.PasswordHash = newUser.PasswordHash
	user.ResetToken = ""

	// Update password via userStore interface
	if err := s.userStore.Update(user); err != nil {
		return &proto.NewPasswordResponse{
			Success: false,
			Message: "Failed to update password",
		}, nil
	}

	return &proto.NewPasswordResponse{
		Success: true,
		Message: "Password has been reset successfully",
	}, nil
}

func (s *AuthServer) GetUserInfo(ctx context.Context, req *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	// Validate session token
	userID, exists := s.sessionMgr.session[req.SessionToken]
	if !exists {
		return &proto.UserInfoResponse{
			Success: false,
			Message: "Invalid session token",
		}, nil
	}

	// Get user info
	user, err := s.userStore.GetByID(userID)
	if err != nil {
		return &proto.UserInfoResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	return &proto.UserInfoResponse{
		Success:  true,
		Message:  "User information received successfully",
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
