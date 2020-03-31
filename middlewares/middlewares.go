package middlewares

import (
	"net/http"
	"time"

	"github.com/covid-19-ppe-tracker/hospital-module-kit/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
)

type contextKey int

const (
	//SerializerContextKey is used to set data to be serialized in request context
	SerializerContextKey contextKey = iota
	DeserializerContextKey
	UserContextKey
)

const (
	JwtForKey               = "jk"
	JwtForAuth              = "auth"
	JwtForEmailVerification = "emailVerification"
	ApiV1                   = "api/v1"
)

type pattern string

var routes = make(map[string]http.Handler)

//noinspection GoExportedFuncWithUnexportedType
func Handle(p string, h http.Handler) pattern {
	routes[p] = h
	return pattern(p)
}

func GetRoutes() map[string]http.Handler {
	return routes
}

func getJWTTokenForAuth(user models.User, secret string, d time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(d).Unix()
	claims[JwtForKey] = JwtForAuth
	return token.SignedString([]byte(secret))
}

func GetJWTTokensForAuth(user models.User, secret string) (string, string, error) {
	// Create the accessToken for auth
	at, err := getJWTTokenForAuth(user, secret, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	// Create the refreshToken for auth
	rt, err := getJWTTokenForAuth(user, secret, time.Hour*24)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func GetJWTTokenForEmailVerification(_id primitive.ObjectID, secret string) (string, error) {
	// Create the token for email verification
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = _id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Second * 3600 * 24).Unix()
	claims[JwtForKey] = JwtForEmailVerification
	return token.SignedString([]byte(secret))
}
