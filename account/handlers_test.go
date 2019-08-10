package account

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleRegister_POST(t *testing.T) {
	router := httprouter.New()
	router.POST("/", HandleRegister_POST(TestDB))

	var registerReq RegisterRequest
	registerReq.Username = "John123"
	registerReq.Password = "password"
	registerReq.Name = "John"

	registerReqJson, err := json.Marshal(registerReq)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(registerReqJson))
	if err != nil {
		t.Errorf("Unexpected error, got %s", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	respBody, err := ioutil.ReadAll(rr.Body)
	if (rr.Code != http.StatusCreated || string(respBody) != "Account Created") {
		t.Errorf("Expected status code to be 201 and response to be Account Created, got %s, %s", string(rr.Code), string(respBody))
	}
}

func TestHandleLogin_POST(t *testing.T) {
	router := httprouter.New()
	router.POST("/", HandleLogin_POST(TestDB))

	var loginReq LoginRequest
	loginReq.Username = "John123"
	loginReq.Password = "password"

	_ = TestDB.DB.Put([]byte(loginReq.Username), []byte(loginReq.Password), nil)

	loginReqJson, err := json.Marshal(loginReq)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(loginReqJson))
	if err != nil {
		t.Errorf("Unexpected error, got %s", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if (rr.Code != http.StatusOK) {
		t.Errorf("Expected status code to be 200, got %s", string(rr.Code))
	}
}

func TestHandleProtectedEndpoint_GET(t *testing.T) {
	router := httprouter.New()
	router.GET("/", HandleProtectedEndpoint_GET())

	var loginReq LoginRequest
	loginReq.Username = "John123"
	loginReq.Password = "password"

	jwt, _ := CreateJwt(loginReq.Username, time.Time{})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err)
	}
	req.Header.Set("jwt", jwt)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if (rr.Code != http.StatusOK) {
		t.Errorf("Expected status code to be 200, got %s", string(rr.Code))
	}

	// test case when jwt invalid
	t.Log("Some errors might be printed as result of testing with invalid jwt")
	jwt = "invalid jwt"

	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Unexpected error, got %s", err)
	}
	req.Header.Set("jwt", jwt)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if (rr.Code != http.StatusUnauthorized) {
		t.Errorf("Expected status code to be 401, got %s", string(rr.Code))
	}
}