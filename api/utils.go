package api

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func authMiddleware(next func(http.ResponseWriter, *http.Request), params ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tS := r.Header.Get("Authorization")

		if claims, err := validateJWTToken(tS); err == nil {
			permissions := claims["permissions"].([]string)
			// TODO: adjust permissions
			if len(permissions) > 0 {
				next(w, r)
				return
			}
		}

		sendErrorResponse(w, errors.New("Not authorized."), http.StatusUnauthorized)
	})
}

func loadTemplate(templateName string) (*template.Template, error) {

	// Assemble template path
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templatePath := filepath.Join(dir, "/template/"+templateName)
	return template.New(templateName).Delims("[[", "]]").ParseFiles(templatePath)
}

func emptyObject() map[string]interface{} {
	return map[string]interface{}{}
}

func emptyArray() []interface{} {
	return []interface{}{}
}

func sendErrorResponse(w http.ResponseWriter, err error, status int) {
	payload := map[string]string{"error": err.Error()}
	sendResponse(w, payload, status)
}

func sendResponse(w http.ResponseWriter, payload interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func createJWTToken(claims jwt.MapClaims) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte{123})
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func getSigningSecret() []byte {
	// TODO: add actual secret
	return []byte{123}
}

func validateJWTToken(t string) (jwt.MapClaims, error) {
	if t == "" {
		return nil, errors.New("Empty JWT.")
	}

	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return getSigningSecret(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token.")
}

func verifyPassword(password string, passwordHash []byte) bool {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	return err == nil
}

func generatePasswordHash(pw string) ([]byte, error) {
	// Use bcrypt to (automatically) salt and encrypt password string
	return bcrypt.GenerateFromPassword([]byte(pw), 10)
}
