package db

import (
	"MyChatApp/monolithic/internal/models"
	"time"

	"github.com/gocql/gocql"
)

type FriendshipRepository struct {
	session *gocql.Session
}

func NewFriendshipRepository(session *gocql.Session) *FriendshipRepository {
	return &FriendshipRepository{
		session: session,
	}
}

/*
DELETE	/friend-requests/:sender_id	Reject (or cancel if I'm the sender)

GET	/friends	List my mutual friends

GET /friend-requests
{
  "incoming": [
    {
      "sender_id": "uuid",
      "sender_username": "bob",
      "sender_first_name": "Bob",
      "sender_last_name": "Smith",
      "created_at": "iso-string"
    }
  ]
}

GET /friends
[
  {
    "user_id": "uuid",
    "username": "bob",
    "first_name": "Bob",
    "last_name": "Smith"
  }
]
*/

// GET	/friend-requests	List incoming pending requests for me (recipient id == me)
func (f *FriendshipRepository) GetFriendRequests(userID gocql.UUID) ([]models.FriendRequest, error) {
	// Query() for a single operation - Iter() for multiple returned rows, Scan() directly for a single returned row
	iter := f.session.Query(
		"SELECT recipient_id, sender_id, status, created_at FROM friend_requests WHERE recipient_id = ?",
		userID,
	).Iter()

	var requests []models.FriendRequest
	var recipient_id gocql.UUID
	var sender_id gocql.UUID
	var status string
	var created_at time.Time

	// Iter() returns the first page from cassandra (5000 by default)
	// gocql transparently fetches more as we exhaust results
	// data until now is still in raw byte form, iter.Scan does the conversion
	// here, Scan() returns true until we run out of values
	for iter.Scan(&recipient_id, &sender_id, &status, &created_at) {
		requests = append(requests, models.FriendRequest{
			RecipientID: recipient_id,
			SenderID:    sender_id,
			Status:      status,
			CreatedAt:   created_at,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return requests, nil
}

// POST	/friend-requests	Send a request to someone { "recipient_id": "uuid-string" }
func (f *FriendshipRepository) CreateFriendRequest(fr *models.FriendRequest) error {
	return f.session.Query(
		"INSERT INTO friend_requests (recipient_id, sender_id, status, created_at) VALUES (?, ?, ?, ?)",
		fr.RecipientID, fr.SenderID, fr.Status, fr.CreatedAt,
	).Exec()
}

// PUT	/friend-requests/:sender_id/accept	Accept a request
func (f *FriendshipRepository) AcceptFriendRequest(fr *models.FriendRequest) error {

	/*
		TODO: we are doing cross-partition writes (3 different partitions).
		A logged batch across partitions is more expensive than a single-partition write because
		Cassandra has to use the batch log to guarantee atomicity.
		Fine for now since accepting friend requests is probably rare.
	*/

	batch := f.session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		"DELETE FROM friend_requests WHERE recipient_id = ? AND sender_id = ?",
		fr.RecipientID, fr.SenderID,
	)

	/*
		The primary key is (user_id, sender_id). Here, the user_id is the partition key, and
		Cassandra does not enforce global uniqueness across the (user_id, friend_id) tuple.
		Therefore we need to do both (user_id, sender_id) and (sender_id, user_id)
		so that everyone can look up their friends efficiently.
	*/

	batch.Query(
		"INSERT INTO friendships (user_id, friend_id, created_at) VALUES (?, ?, ?)",
		fr.RecipientID, fr.SenderID, fr.CreatedAt,
	)

	batch.Query(
		"INSERT INTO friendships (user_id, friend_id, created_at) VALUES (?, ?, ?)",
		fr.SenderID, fr.RecipientID, fr.CreatedAt,
	)

	return f.session.ExecuteBatch(batch)
}
