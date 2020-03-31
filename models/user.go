package models

import (
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userType int

const (
	//SerializerContextKey is used to set data to be serialized in request context
	AvailabilityUserType userType = iota
	RequirementUserType
)

//User in-mem represents a users collection document
type User struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone        string             `json:"phone,omitempty" bson:"phone,omitempty"`
	CountryCode  string             `json:"countryCode,omitempty" bson:"countryCode,omitempty"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty"`
	AccessToken  string             `json:"accessToken,omitempty" bson:",omitempty"`
	RefreshToken string             `json:"refreshToken,omitempty" bson:",omitempty"`
	ErrorMessage string             `json:"error,omitempty" bson:"-"`
}

//ComparePassword compares password in request with password in DB
func (u User) ComparePassword(candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(candidatePassword), []byte(u.Password))
	if err != nil {
		return false
	}
	return true
}

func (u User) SerializeToJSON(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(&u)
}

func (u User) DeserializeFromJSONInBody(r *http.Request) JSONBodyDeserializer {
	decoder := json.NewDecoder(r.Body)
	for {
		if err := decoder.Decode(&u); err == io.EOF || err != nil {
			break
		}
	}
	return u
}

//AuthenticateEmail authenticates user based on email
func (u User) AuthenticateEmail() (User, bool, error) {
	var candidateUser = User{}
	err := FindOne(userModel, &candidateUser, bson.D{{"email", u.Email}})
	if err != nil {
		return candidateUser, false, err
	}
	if !u.ComparePassword(candidateUser.Password) {
		return candidateUser, false, nil
	}
	return candidateUser, true, nil
}

//NewUserFromRequestBody gets a user from request payload
func NewUserFromRequestBody(r io.Reader) *User {
	var user User
	decoder := json.NewDecoder(r)
	for {
		if err := decoder.Decode(&user); err == io.EOF || err != nil {
			break
		}
	}
	return &user
}

//UserModel represents a user model
type UserModel struct {
	b2cdb
}

func (u UserModel) collectionName() string {
	return "users"
}

var userModel = &UserModel{}

//GetUserModel returns a user model
func GetUserModel() *UserModel {
	return userModel
}
