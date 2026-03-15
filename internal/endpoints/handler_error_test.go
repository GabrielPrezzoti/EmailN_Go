package endpoints

import (
	internalerrors "emailn/internal/internalErrors"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandlerError_when_endpoint_returns_intern_error(t *testing.T) {
	assert := assert.New(t)
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, 0, internalerrors.ErrInternal
	}
	handlerFunc := HandlerError(endpoint)
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusInternalServerError, response.Code)
	assert.Contains(response.Body.String(), internalerrors.ErrInternal.Error())
}

func Test_HandlerError_when_endpoint_returns_domain_error(t *testing.T) {
	assert := assert.New(t)
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, 0, errors.New("domain error")
	}
	handlerFunc := HandlerError(endpoint)
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusBadRequest, response.Code)
	assert.Contains(response.Body.String(), "domain error")
}

func Test_HandlerError_when_endpoint_returns_obj_and_status(t *testing.T) {
	assert := assert.New(t)
	type BodyForTest struct {
		Id int
	}

	objExpected := BodyForTest{Id: 2}
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return objExpected, 201, nil
	}

	handlerFunc := HandlerError(endpoint)
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	handlerFunc.ServeHTTP(response, request)

	assert.Equal(http.StatusCreated, response.Code)
	objReturned := BodyForTest{}
	json.Unmarshal(response.Body.Bytes(), &objReturned)
	assert.Equal(objExpected, objReturned)
}
