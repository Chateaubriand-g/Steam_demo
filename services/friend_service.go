package services

import (
	"errors"
	"steam-backend/models"
	"steam-backend/repositories"

	"gorm.io/gorm"
)

type FriendService interface {
	SendInvitation(senderID, receiverID uint64) error
	AcceptInvitation(invitationID, userID uint64) error
	RefuseInvitation(invitationID, userID uint64) error
	GetReceivedInvitations(userID uint64) ([]models.Invitation, error)
	GetSentInvitations(userID uint64) ([]models.Invitation, error)
	RemoveFriend(userID, friendID uint64) error
	CheckFriendShip(userID, friendID uint64) (bool, error)
	GetFriendCount(userID int64) (int64, error)
	GetFriendList(userID int64) ([]models.User, error)
}

type friendService struct {
	friendRepo repositories.FriendRepository
}

func NewFriendService(repo repositories.FriendRepository) FriendService {
	return &friendService{friendRepo: repo}
}

func (s *friendService) SendInvitation(senderID, receiverID uint64) error {
	if senderID == receiverID {
		return errors.New("can not invited self")
	}
	isFriend, err := s.friendRepo.IsFriends(senderID, receiverID)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already was friend")
	}
	_, err_findInv := s.friendRepo.FindInvitation(senderID, receiverID)
	if err_findInv == nil {
		return errors.New("already send invitation")
	}
	if !errors.Is(err_findInv, gorm.ErrRecordNotFound) {
		return err_findInv
	}

	invitation := models.Invitation{
		SenderID:   senderID,
		ReceiverID: receiverID,
	}
	return s.friendRepo.CreateInvitation(&invitation)
}

func (s *friendService) AcceptInvitation(invitationID, userID uint64) error {
	return s.friendRepo.AcceptInvitation(invitationID, userID)
}

func (s *friendService) RefuseInvitation(invitationID, userID uint64) error {
	return s.friendRepo.RefuseInvitation(invitationID, userID)
}

func (s *friendService) GetReceivedInvitations(userID uint64) ([]models.Invitation, error) {
	return s.friendRepo.GetInvitationByReceiver(userID, "pending")
}

func (s *friendService) GetSentInvitations(userID uint64) ([]models.Invitation, error) {
	return s.friendRepo.GetInvitationBySender(userID, "pending")
}

func (s *friendService) RemoveFriend(userID, friendID uint64) error {
	isFriend, err := s.CheckFriendShip(userID, friendID)
	if err != nil {
		return err
	}
	if !isFriend {
		return errors.New("not friend yet")
	}
	return s.friendRepo.DeleteFriendship(userID, friendID)
}

func (s *friendService) CheckFriendShip(userID, friendID uint64) (bool, error) {
	return s.friendRepo.IsFriends(userID, friendID)
}

func (s *friendService) GetFriendCount(userID int64) (int64, error) {
	return s.friendRepo.GetFriendCount(uint64(userID))
}

func (s *friendService) GetFriendList(userID int64) ([]models.User, error) {
	return s.friendRepo.GetFriendList(uint64(userID))
}
