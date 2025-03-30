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

	if username == "" || email == "" || password == "" {
		return "", fmt.Errorf("username, email and password are required")
	}

	existingUser, err := ins.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
		} else {
			log.Printf("Error checking existing username: %v\n", err)
			return "", fmt.Errorf("error checking existing username: %v", err)
		}
	} else if existingUser != nil {
		return "", fmt.Errorf("username already exists")
	}

	_, err = ins.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", fmt.Errorf("error checking existing email: %v", err)
		}
	}

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

	id, err = ins.repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return "", fmt.Errorf("error creating user: %v", err)
	}
	return id, nil
}

func (ins *AuthService) Login(ctx context.Context, identifier, password string) (id string, err error) {
	user, err := ins.repo.GetUserByUsername(ctx, identifier)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return user.ID.Hex(), nil
}
