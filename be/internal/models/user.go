package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      string        `bson:"username" json:"username"`
	Email         string        `bson:"email" json:"email"`
	PasswordHash  string        `bson:"passwordHash" json:"-"`
	FirstName     string        `bson:"firstName" json:"firstName"`
	LastName      string        `bson:"lastName" json:"lastName"`
	PhoneNumber   string        `bson:"phoneNumber" json:"phoneNumber"`
	Address       string        `bson:"address" json:"address"`
	City          string        `bson:"city" json:"city"`
	ZipCode       string        `bson:"zipCode" json:"zipCode"`
	IDCardFront   string        `bson:"idCardFront" json:"idCardFront"`
	IDCardBack    string        `bson:"idCardBack" json:"idCardBack"`
	Verified      bool          `bson:"verified" json:"verified"`
	Status        string        `bson:"status" json:"status"`
	TermsAccepted bool          `bson:"termsAccepted" json:"termsAccepted"`
	CreatedAt     time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time     `bson:"updatedAt" json:"updatedAt"`
}
