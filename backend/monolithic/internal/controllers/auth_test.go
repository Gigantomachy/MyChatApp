package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"MyChatApp/monolithic/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	testutil.TruncateAll(t)

	user, cookie := testutil.RegisterUser(t, "alice")
	assert.NotEmpty(t, user.UserID)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, "Alice", user.FirstName)
	assert.NotEmpty(t, cookie)
}

func TestRegisterDuplicateUser(t *testing.T) {
	testutil.TruncateAll(t)

	testutil.RegisterUser(t, "bob")

	body := map[string]string{
		"username":   "bob",
		"password":   "password123",
		"email":      "bob2@test.com",
		"first_name": "Bob",
		"last_name":  "Second",
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/register", body, "")

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errResp testutil.ErrorResponse
	json.Unmarshal(respBody, &errResp)
	assert.Contains(t, errResp.Error, "already taken")
}

func TestLoginNonExistentUser(t *testing.T) {
	testutil.TruncateAll(t)

	body := map[string]string{
		"username": "ghost",
		"password": "password123",
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/login", body, "")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errResp testutil.ErrorResponse
	json.Unmarshal(respBody, &errResp)
	assert.Contains(t, errResp.Error, "invalid credentials")
}

func TestLoginWrongPassword(t *testing.T) {
	testutil.TruncateAll(t)

	testutil.RegisterUser(t, "charlie")

	body := map[string]string{
		"username": "charlie",
		"password": "wrongpassword",
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/login", body, "")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errResp testutil.ErrorResponse
	json.Unmarshal(respBody, &errResp)
	assert.Contains(t, errResp.Error, "invalid credentials")
}

func TestLoginSuccess(t *testing.T) {
	testutil.TruncateAll(t)

	testutil.RegisterUser(t, "dave")

	user, cookie := testutil.LoginUser(t, "dave")
	assert.Equal(t, "dave", user.Username)
	assert.NotEmpty(t, cookie)
}

func TestMeWithCookie(t *testing.T) {
	testutil.TruncateAll(t)

	_, cookie := testutil.RegisterUser(t, "eve")

	resp, respBody := testutil.DoRequest(t, "GET", "/api/me", nil, cookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		User testutil.User `json:"user"`
	}
	json.Unmarshal(respBody, &result)
	assert.Equal(t, "eve", result.User.Username)
}

func TestMeWithoutCookie(t *testing.T) {
	testutil.TruncateAll(t)

	resp, _ := testutil.DoRequest(t, "GET", "/api/me", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
