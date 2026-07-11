package endpoints

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImVtYWlsIjoidGVzdGVAdGVzdGUuY29tIn0.pN9FoK4c5YQRLsqErfKbX9nRloJiBVcQ6ju5VjkDCtc"
var validEmail string = "teste@teste.com"

func Test_Auth_WhenAuthorizationIsMissing_ReturnError(t *testing.T) {
	assert := assert.New(t)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})
	handlerFunc := Auth(nextHandler)
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusUnauthorized, response.Code)
	assert.Contains(response.Body.String(), "request does not contain an authorization header")
}

func Test_Auth_WhenAuthorizationIsInvalid_ReturnError(t *testing.T) {
	assert := assert.New(t)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})
	ValidateToken = func(token string, ctx context.Context) error {
		return errors.New("invalid token")
	}
	handlerFunc := Auth(nextHandler)
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "Bearer invalid-token")
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusUnauthorized, response.Code)
	assert.Contains(response.Body.String(), "invalid token")
}

func Test_Auth_WhenAuthorizationIsValid_CallNextHandler(t *testing.T) {
	assert := assert.New(t)
	var email string
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email = r.Context().Value("email").(string)
	})
	ValidateToken = func(token string, ctx context.Context) error {
		return nil
	}
	handlerFunc := Auth(nextHandler)
	request, _ := http.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", validToken)
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusOK, response.Code)
	assert.Equal(validEmail, email)
}
