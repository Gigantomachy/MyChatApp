package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"MyChatApp/monolithic/internal/controllers"
	"MyChatApp/monolithic/internal/repositories/db"
	"MyChatApp/monolithic/internal/routers"
	"MyChatApp/monolithic/internal/services"
	"MyChatApp/monolithic/internal/ws"

	"github.com/gocql/gocql"
)

const testKeyspace = "my_chat_app_test"

var (
	testSession *gocql.Session
	testServer  *httptest.Server
)

func Setup() {
	var err error
	testSession, err = createTestDB()
	if err != nil {
		panic(err)
	}
	testServer = createServer(testSession)
}

func Teardown() {
	testServer.Close()
	testSession.Query("DROP KEYSPACE IF EXISTS " + testKeyspace).Exec()
	testSession.Close()
}

func createTestDB() (*gocql.Session, error) {
	// Connect without keyspace to create it
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS " + testKeyspace +
		" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec()
	session.Close()
	if err != nil {
		return nil, err
	}

	// Reconnect with the test keyspace
	cluster.Keyspace = testKeyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	for _, stmt := range db.SchemaStatements {
		err = session.Query(stmt).Exec()
		if err != nil {
			session.Close()
			return nil, err
		}
	}

	return session, nil
}

func createServer(session *gocql.Session) *httptest.Server {
	userRepo := db.NewUserRepository(session)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(userService)

	friendRepo := db.NewFriendshipRepository(session)
	friendService := services.NewFriendService(friendRepo, userRepo)

	hub := ws.NewHub()

	friendController := controllers.NewFriendController(friendService, hub)
	userController := controllers.NewUserController(userService)

	channelsRepo := db.NewChannelsRepository(session)
	channelsService := services.NewChannelService(channelsRepo, userRepo, friendRepo)
	channelsController := controllers.NewChannelsController(channelsService, hub)

	messagesRepo := db.NewMessagesRepository(session)
	messagesService := services.NewMessageService(messagesRepo, channelsRepo, userRepo)
	messagesController := controllers.NewMessagesController(messagesService, hub)

	r := routers.NewRouter(routers.RouterDependencies{
		Hub:               hub,
		AuthController:    authController,
		FriendController:  friendController,
		UserController:    userController,
		ChannelController: channelsController,
		MessageController: messagesController,
	})

	return httptest.NewServer(r.Setup())
}

func TruncateAll(t *testing.T) {
	t.Helper()
	for _, stmt := range db.TruncateStatements {
		if err := testSession.Query(stmt).Exec(); err != nil {
			t.Fatalf("Failed to truncate: %v", err)
		}
	}
}

func ServerURL() string {
	return testServer.URL
}

// --- HTTP helpers ---

func DoRequest(t *testing.T, method, path string, body interface{}, cookie string) (*http.Response, []byte) {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	return resp, respBody
}

func ExtractCookie(resp *http.Response) string {
	for _, c := range resp.Cookies() {
		if c.Name == "auth_token" {
			return "auth_token=" + c.Value
		}
	}
	return ""
}

// --- Types ---

type User struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type FriendRequestItem struct {
	SenderID       string `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	SenderFirstName string `json:"sender_first_name"`
	SenderLastName  string `json:"sender_last_name"`
}

type Channel struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type ChannelMembership struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	ChannelType string `json:"channel_type"`
}

type MessageItem struct {
	MessageID       string `json:"message_id"`
	AuthorID        string `json:"author_id"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
	AuthorUsername  string `json:"author_username"`
	Content         string `json:"content"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --- Seed helpers ---

func RegisterUser(t *testing.T, username string) (User, string) {
	t.Helper()
	body := map[string]string{
		"username":   username,
		"password":   "password123",
		"email":      username + "@test.com",
		"first_name": strings.ToUpper(username[:1]) + username[1:],
		"last_name":  "Test",
	}
	resp, respBody := DoRequest(t, "POST", "/api/register", body, "")
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Register %q failed: %d %s", username, resp.StatusCode, string(respBody))
	}
	var result AuthResponse
	json.Unmarshal(respBody, &result)
	return result.User, ExtractCookie(resp)
}

func LoginUser(t *testing.T, username string) (User, string) {
	t.Helper()
	body := map[string]string{
		"username": username,
		"password": "password123",
	}
	resp, respBody := DoRequest(t, "POST", "/api/login", body, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login %q failed: %d %s", username, resp.StatusCode, string(respBody))
	}
	var result AuthResponse
	json.Unmarshal(respBody, &result)
	return result.User, ExtractCookie(resp)
}
