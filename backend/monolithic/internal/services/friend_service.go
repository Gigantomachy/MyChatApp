package services

import (
	"MyChatApp/monolithic/internal/models"
	"MyChatApp/monolithic/internal/repositories/db"

	"github.com/gocql/gocql"
	"time"
)

type FriendService struct {
	friendRepo *db.FriendshipRepository
	userRepo   *db.UserRepository
}

func NewFriendService(f *db.FriendshipRepository, u *db.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: f,
		userRepo:   u,
	}
}

// TODO: replace this with a DTO later
type FriendRequestItem struct {
	SenderID  gocql.UUID `json:"sender_id"`
	Username  string     `json:"sender_username"`
	FirstName string     `json:"sender_first_name"`
	LastName  string     `json:"sender_last_name"`
	CreatedAt time.Time  `json:"created_at"`
}

// manage friends and friend requests

// get my current friends
func (fs *FriendService) GetFriends(user_id string) ([]models.User, error) {
	userID, err := gocql.ParseUUID(user_id)
	if err != nil {
		return []models.User{}, err
	}

	friendShips, err := fs.friendRepo.GetFriendships(userID)
	if err != nil {
		return []models.User{}, err
	}

	var friendIDs []gocql.UUID
	for _, fsh := range friendShips {
		friendIDs = append(friendIDs, fsh.FriendID)
	}

	friends, err := fs.userRepo.FindByIDs(friendIDs)
	if err != nil {
		return []models.User{}, err
	}

	return friends, nil
}

// remove a friend
func (fs *FriendService) RemoveFriend(user_id, friend_id string) error {
	userID, err := gocql.ParseUUID(user_id)
	if err != nil {
		return err
	}

	friendID, err := gocql.ParseUUID(friend_id)
	if err != nil {
		return err
	}

	err = fs.friendRepo.DeleteFriendship(userID, friendID)
	return err
}

func (fs *FriendService) SendFriendRequest(user_id, recipient_id string) (*models.FriendRequest, error) {
	userID, err := gocql.ParseUUID(user_id)
	if err != nil {
		return nil, err
	}

	recipientID, err := gocql.ParseUUID(recipient_id)
	if err != nil {
		return nil, err
	}

	fr := models.FriendRequest{
		RecipientID: recipientID,
		SenderID:    userID,
		Status:      "PENDING",
		CreatedAt:   time.Now(),
	}

	return &fr, fs.friendRepo.CreateFriendRequest(&fr)
}

func (fs *FriendService) AcceptFriendRequest(recipient_id, sender_id string) error {
	recipientID, err := gocql.ParseUUID(recipient_id) // the current user
	if err != nil {
		return err
	}

	senderID, err := gocql.ParseUUID(sender_id) // the current user
	if err != nil {
		return err
	}

	fr := models.FriendRequest{
		RecipientID: recipientID,
		SenderID:    senderID,
		Status:      "ACCEPTED",
		CreatedAt:   time.Now(),
	}

	return fs.friendRepo.AcceptFriendRequest(&fr)
}

func (fs *FriendService) DeleteFriendRequest(sender_id, recipient_id string) error {
	senderID, err := gocql.ParseUUID(sender_id)
	if err != nil {
		return err
	}

	recipientID, err := gocql.ParseUUID(recipient_id)
	if err != nil {
		return err
	}

	return fs.friendRepo.DeleteFriendRequest(recipientID, senderID)
}

/*
	{
	  "sender_id": "uuid",
	  "sender_username": "bob",
	  "sender_first_name": "Bob",
	  "sender_last_name": "Smith"
	}
*/
func (fs *FriendService) GetFriendRequests(user_id string) ([]FriendRequestItem, error) {
	userID, err := gocql.ParseUUID(user_id)
	if err != nil {
		return []FriendRequestItem{}, err
	}

	friendrequests, err := fs.friendRepo.GetFriendRequests(userID)
	if err != nil {
		return []FriendRequestItem{}, err
	}

	if len(friendrequests) == 0 {
		return []FriendRequestItem{}, nil
	}

	var createdAtMap = make(map[gocql.UUID]time.Time)
	var senderIDs []gocql.UUID
	for _, req := range friendrequests {
		senderIDs = append(senderIDs, req.SenderID)
		createdAtMap[req.SenderID] = req.CreatedAt
	}

	var friendRequestDTOs []FriendRequestItem
	friends, err := fs.userRepo.FindByIDs(senderIDs)
	if err != nil {
		return []FriendRequestItem{}, err
	}

	for _, friend := range friends {
		dto := FriendRequestItem{
			SenderID:  friend.UserID,
			FirstName: friend.FirstName,
			LastName:  friend.LastName,
			Username:  friend.Username,
			CreatedAt: createdAtMap[friend.UserID],
		}

		friendRequestDTOs = append(friendRequestDTOs, dto)
	}

	return friendRequestDTOs, nil
}
