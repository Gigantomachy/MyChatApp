package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"MyChatApp/monolithic/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestCreatePublicChannelAndJoin(t *testing.T) {
	testutil.TruncateAll(t)

	alice, aliceCookie := testutil.RegisterUser(t, "alice")
	_, bobCookie := testutil.RegisterUser(t, "bob")
	_, carolCookie := testutil.RegisterUser(t, "carol")

	// Alice creates a public channel
	body := map[string]interface{}{
		"channel_name": "#general",
		"channel_type": "public",
	}
	resp, respBody := testutil.DoRequest(t, "POST", "/api/channels", body, aliceCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var ch testutil.Channel
	json.Unmarshal(respBody, &ch)
	assert.Equal(t, "#general", ch.Name)
	assert.Equal(t, "public", ch.Type)
	assert.NotEmpty(t, ch.ChannelID)

	// Bob (non-friend) joins the channel
	resp, _ = testutil.DoRequest(t, "POST", "/api/channels/"+ch.ChannelID, nil, bobCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Carol (non-friend) joins the channel
	resp, _ = testutil.DoRequest(t, "POST", "/api/channels/"+ch.ChannelID, nil, carolCookie)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// All three see the channel in their sidebar
	for _, cookie := range []string{aliceCookie, bobCookie, carolCookie} {
		resp, respBody = testutil.DoRequest(t, "GET", "/api/channels", nil, cookie)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var memberships []testutil.ChannelMembership
		json.Unmarshal(respBody, &memberships)
		assert.Len(t, memberships, 1)
		assert.Equal(t, "#general", memberships[0].ChannelName)
		assert.Equal(t, "public", memberships[0].ChannelType)
	}

	// Channel has 3 members
	resp, respBody = testutil.DoRequest(t, "GET", "/api/channels/"+ch.ChannelID, nil, aliceCookie)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var details struct {
		ChannelID string             `json:"channel_id"`
		Users     []testutil.User    `json:"users"`
	}
	json.Unmarshal(respBody, &details)
	assert.Len(t, details.Users, 3)

	_ = alice
}
