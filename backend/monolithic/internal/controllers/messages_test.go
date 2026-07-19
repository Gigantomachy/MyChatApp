package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"MyChatApp/monolithic/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestDMWithNonFriendFails(t *testing.T) {
	testutil.TruncateAll(t)

	_, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, _ := testutil.RegisterUser(t, "bob")

	// Alice tries to DM Bob (not friends)
	body := map[string]interface{}{
		"channel_type": "dm",
		"members":      []string{bob.UserID},
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var errResp testutil.ErrorResponse
	json.Unmarshal(respBody, &errResp)
	assert.Contains(t, errResp.Error, "friends")
}

func TestDMWithFriendAndSendMessage(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, bobCookie := testutil.RegisterUser(t, "bob")

	// Make them friends
	testutil.DoRequest(t, "POST", "/api/friend-requests/"+bob.UserID, nil, aliceCookie)
	testutil.DoRequest(t, "PUT", "/api/friend-requests/"+alice.UserID, nil, bobCookie)

	// Alice starts a DM with Bob
	body := map[string]interface{}{
		"channel_type": "dm",
		"members":      []string{bob.UserID},
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var ch testutil.Channel
	json.Unmarshal(respBody, &ch)
	assert.Equal(t, "dm", ch.Type)
	channelID := ch.ChannelID

	// Both see the DM in their sidebar
	resp, respBody = testutil.DoRequest(t, "GET", "/api/channels", nil, aliceCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var memberships []testutil.ChannelMembership
	json.Unmarshal(respBody, &memberships)
	assert.Len(t, memberships, 1)
	assert.Equal(t, "dm", memberships[0].ChannelType)

	resp, respBody = testutil.DoRequest(t, "GET", "/api/channels", nil, bobCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	json.Unmarshal(respBody, &memberships)
	assert.Len(t, memberships, 1)

	// Alice sends a message
	msgBody := map[string]string{"content": "hello bob!"}
	resp, respBody = testutil.DoRequest(t, "POST", "/api/channels/"+channelID+"/messages", msgBody, aliceCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var msg testutil.MessageItem
	json.Unmarshal(respBody, &msg)
	assert.Equal(t, "hello bob!", msg.Content)
	assert.Equal(t, alice.UserID, msg.AuthorID)
	assert.Equal(t, "Alice", msg.AuthorFirstName)

	// Bob sends a message
	msgBody = map[string]string{"content": "hi alice!"}
	resp, _ = testutil.DoRequest(t, "POST", "/api/channels/"+channelID+"/messages", msgBody, bobCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Alice fetches messages — sees both
	resp, respBody = testutil.DoRequest(t, "GET", "/api/channels/"+channelID+"/messages", nil, aliceCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var messages []testutil.MessageItem
	json.Unmarshal(respBody, &messages)
	assert.Len(t, messages, 2)
}

func TestNonMemberCannotSendMessages(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, bobCookie := testutil.RegisterUser(t, "bob")
	_, carolCookie := testutil.RegisterUser(t, "carol")

	// Alice and Bob become friends
	testutil.DoRequest(t, "POST", "/api/friend-requests/"+bob.UserID, nil, aliceCookie)
	testutil.DoRequest(t, "PUT", "/api/friend-requests/"+alice.UserID, nil, bobCookie)

	// Alice creates a DM with Bob
	body := map[string]interface{}{
		"channel_type": "dm",
		"members":      []string{bob.UserID},
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	var ch testutil.Channel
	json.Unmarshal(respBody, &ch)

	// Carol (not a member) tries to send a message
	msgBody := map[string]string{"content": "intruder!"}
	resp, _ = testutil.DoRequest(t, "POST", "/api/channels/"+ch.ChannelID+"/messages", msgBody, carolCookie)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestDMDedup(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	bob, bobCookie := testutil.RegisterUser(t, "bob")

	// Make them friends
	testutil.DoRequest(t, "POST", "/api/friend-requests/"+bob.UserID, nil, aliceCookie)
	testutil.DoRequest(t, "PUT", "/api/friend-requests/"+alice.UserID, nil, bobCookie)

	// Alice creates a DM with Bob
	body := map[string]interface{}{
		"channel_type": "dm",
		"members":      []string{bob.UserID},
	}
	_, respBody := testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	var ch1 testutil.Channel
	json.Unmarshal(respBody, &ch1)

	// Alice tries to create another DM with Bob — should return the same channel
	_, respBody = testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	var ch2 testutil.Channel
	json.Unmarshal(respBody, &ch2)

	assert.Equal(t, ch1.ChannelID, ch2.ChannelID)
}
