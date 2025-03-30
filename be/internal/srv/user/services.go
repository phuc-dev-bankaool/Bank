package user

import (
	"app/be/internal/models"
	"app/be/internal/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (ins *AuthService) Register(ctx context.Context,
	username, email,
	password, firstName, lastName, phoneNumber,
	address, city, zipCode, idCardFront, idCardBack string) (
	id string, err error) {

	// Add validation
	if username == "" || email == "" || password == "" {
		return "", fmt.Errorf("username, email and password are required")
	}

	// Log the registration attempt
	log.Printf("Attempting to register user: %s, email: %s\n", username, email)

	// Check if username exists
	_, err = ins.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", fmt.Errorf("error checking existing username: %v", err)
		}
	}
	_, err = ins.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", fmt.Errorf("error checking existing email: %v", err)
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumber:  phoneNumber,
		Address:      address,
		City:         city,
		ZipCode:      zipCode,
		IDCardFront:  idCardFront,
		IDCardBack:   idCardBack,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Log the user object (remove sensitive data in production)
	log.Printf("Creating user: %+v\n", user)

	// Create user in database
	id, err = ins.repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return "", fmt.Errorf("error creating user: %v", err)
	}
	return id, nil
}

func (ins *AuthService) Login(ctx context.Context, username, password string) (id string, err error) {
	user, err := ins.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return user.ID.Hex(), nil
}
