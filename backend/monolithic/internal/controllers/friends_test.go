package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"MyChatApp/monolithic/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestSendFriendRequestAndReject(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, bobCookie := testutil.RegisterUser(t, "bob")

	// Alice sends friend request to Bob
	resp, _ := testutil.DoRequest(t, "POST", "/api/friend-requests/"+bob.UserID, nil, aliceCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Bob sees the incoming request
	resp, respBody := testutil.DoRequest(t, "GET", "/api/friend-requests", nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var requests []testutil.FriendRequestItem
	json.Unmarshal(respBody, &requests)
	assert.Len(t, requests, 1)
	assert.Equal(t, alice.UserID, requests[0].SenderID)

	// Bob rejects the request
	resp, _ = testutil.DoRequest(t, "DELETE", "/api/friend-requests/incoming/"+alice.UserID, nil, bobCookie)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Bob no longer sees any requests
	resp, respBody = testutil.DoRequest(t, "GET", "/api/friend-requests", nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	json.Unmarshal(respBody, &requests)
	assert.Len(t, requests, 0)

	// Neither is friends
	resp, respBody = testutil.DoRequest(t, "GET", "/api/friends", nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var friends []testutil.User
	json.Unmarshal(respBody, &friends)
	assert.Len(t, friends, 0)
}

func TestSendFriendRequestAndAccept(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, bobCookie := testutil.RegisterUser(t, "bob")

	// Alice sends friend request to Bob
	resp, _ := testutil.DoRequest(t, "POST", "/api/friend-requests/"+bob.UserID, nil, aliceCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Bob accepts the request
	resp, _ = testutil.DoRequest(t, "PUT", "/api/friend-requests/"+alice.UserID, nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bob sees Alice in his friends list
	resp, respBody := testutil.DoRequest(t, "GET", "/api/friends", nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var friends []testutil.User
	json.Unmarshal(respBody, &friends)
	assert.Len(t, friends, 1)
	assert.Equal(t, alice.UserID, friends[0].UserID)

	// Alice sees Bob in her friends list
	resp, respBody = testutil.DoRequest(t, "GET", "/api/friends", nil, aliceCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	json.Unmarshal(respBody, &friends)
	assert.Len(t, friends, 1)
	assert.Equal(t, bob.UserID, friends[0].UserID)
}
