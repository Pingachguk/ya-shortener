package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pingachguk/ya-shortener/config"
	"github.com/pingachguk/ya-shortener/internal/models"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

const CookieName = "auth_access"

type TokenData struct {
	jwt.RegisteredClaims
	UserID string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CheckAuth(r) {
			next.ServeHTTP(w, r)
		}

		_, err := Authenticate(w)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			log.Error().Err(err).Msgf("bad Authenticate")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetTokenData(tokenString string) (*TokenData, error) {
	tokenData := &TokenData{}
	_, err := jwt.ParseWithClaims(tokenString, tokenData, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Config.JWTKey), nil
	})

	return tokenData, err
}

func GetCurrentUser(r http.Request) (*models.User, error) {
	token, err := r.Cookie(CookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return &models.User{}, nil
	} else if err != nil {
		return nil, err
	}

	tokenData, err := GetTokenData(token.Value)
	if err != nil {
		return nil, err
	}

	return &models.User{
		UUID: tokenData.UserID,
	}, nil
}

func CheckAuth(r *http.Request) bool {
	token, err := r.Cookie(CookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return false
	}

	return verifyToken(token.Value)
}

func Authenticate(w http.ResponseWriter) (*models.User, error) {
	tokenString, err := generateToken()
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:    CookieName,
		Value:   tokenString,
		Path:    "/",
		Expires: time.Now().Add(time.Second * time.Duration(config.Config.JWTExpireSeconds)),
	}
	http.SetCookie(w, cookie)

	tokenData, err := GetTokenData(tokenString)
	if err != nil {
		return nil, err
	}

	return &models.User{
		UUID: tokenData.UserID,
	}, nil
}

func verifyToken(tokenString string) bool {
	tokenData := &TokenData{}
	token, err := jwt.ParseWithClaims(tokenString, tokenData, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Config.JWTKey), nil
	})

	return (err == nil) && token.Valid && tokenData.UserID != ""
}

func generateToken() (string, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenData{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(config.Config.JWTExpireSeconds))),
		},
		UserID: userID.String(),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWTKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
