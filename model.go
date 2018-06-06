package main

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const UserCollection = "user"

const (
	UserStatusCreated = iota
	UserStatusActivated
	UserStatusDeleted
)

type User struct {
	Id          bson.ObjectId `bson:"_id"`
	Email       string        `bson:"email"`
	Password    string        `bson:"password"`
	FirstName   string        `bson:"first_name"`
	LastName    string        `bson:"last_name"`
	Status      int           `bson:"status,omitempty"`
	CreatedAt   time.Time     `bson:"created_at,omitempty"`
	UpdatedAt   time.Time     `bson:"updated_at,omitempty"`
	ActivatedAt time.Time     `bson:"activated_at,omitempty"`
	DeletedAt   time.Time     `bson:"deleted_at,omitempty"`
}

// Save user in database
func (u *User) Save(db *mgo.Database) error {
	_, err := db.C(UserCollection).UpsertId(u.Id, u)
	return err
}

// Encrypt password
func (u *User) EncryptPassword(password string) error {
	if hash, err := HashPassword(password); err != nil {
		return err
	} else {
		u.Password = hash
	}

	return nil
}
