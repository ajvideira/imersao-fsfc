package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// User represents a model user
type User struct {
	Base `valid:"required"`
	Name string `json:"name" valid:"notnull"`
	Email string `json:"email" valid:"notnull"`
}

//isValid perform validation of a user
func (user *User) isValid() error {
	_,err := govalidator.ValidateStruct(user);
	if err != nil {
		return err;
	}
	return nil;
}

// NewUser return a new instance of a User
func NewUser(name string, email string) (*User, error) {
	user := User{
		Name: name,
		Email: email,
	}

	err := user.isValid()
	if (err != nil) {
		return nil, err;
	}

	user.ID = uuid.NewV4().String()
	user.CreatedAt = time.Now()

	return &user, nil;
}