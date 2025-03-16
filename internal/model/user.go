package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                string
	Username          string
	Email             string
	PasswordHash      string
	CreatedAt         time.Time
	ResetToken        string
	ResetTokenExpires time.Time
}

// method to create new User instance
func NewUser(username, email, password string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           generateID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}, nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + RandomString(6)
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}

// Validate password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) SetResetToken() string {
	token := generateToken()
	u.ResetToken = token
	u.ResetTokenExpires = time.Now().Add(24 * time.Hour) // Valid for 24 hours
	return token
}

// Utility function: TODO Implement later
func generateToken() string {
	return time.Now().Format("20060102150405") + RandomString(6)
}
