package session

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var tokenSecret = []byte("super secret")

func Check(r *http.Request) (*Session, error) {
	tokenBearerString := r.Header.Get("Authorization")
	tokenStringSplit := strings.Split(tokenBearerString, "Bearer ")
	if len(tokenStringSplit) != 2 {
		return nil, ErrNoAuth
	}
	tokenString := tokenStringSplit[1]

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return tokenSecret, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserJWTClaims{}, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("jwt validation error")
	}

	payload, ok := token.Claims.(*UserJWTClaims)
	if !ok {
		return nil, fmt.Errorf("no payload")
	}

	userID := payload.User.UserID
	username := payload.User.Username
	sess := &Session{UserID: userID, Username: username}

	return sess, nil
}

type UserJWT struct {
	Username string `json:"username"`
	UserID   string `json:"id"`
}

type UserJWTClaims struct {
	User UserJWT `json:"user"`
	jwt.StandardClaims
}

func CreateJWT(userID string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserJWTClaims{User: UserJWT{
		UserID:   userID,
		Username: username,
	}})
	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
