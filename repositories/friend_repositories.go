package repositories

import (
	"errors"
	"steam-backend/models"

	"gorm.io/gorm"
)

type FriendRepository interface {
	CreateFriendship(user1ID, user2ID uint64) error
	DeleteFriendship(user1ID, user2ID uint64) error
	IsFriends(user1ID, user2ID uint64) (bool, error)
	GetFriendCount(userID uint64) (int64, error)
	GetFriendList(userID uint64) ([]models.User, error)
	CreateInvitation(invitation *models.Invitation) error
	UpdateInvitationStatus(id uint64, status string) error
	GetInvitationByReceiver(receiverID uint64, status string) ([]models.Invitation, error)
	GetInvitationBySender(receiverID uint64, status string) ([]models.Invitation, error)
	FindInvitation(senderID, receiverID uint64) (*models.Invitation, error)
	GetInvitationByID(id uint64) (*models.Invitation, error)
	AcceptInvitation(invitationID, userID uint64) error
	RefuseInvitation(invitationID, userID uint64) error
}

type friendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) (FriendRepository, error) {
	if db == nil {
		return nil, errors.New("gorm.DB is nil")
	}
	return &friendRepository{db: db}, nil
}

func (r *friendRepository) CreateFriendship(userID1, userID2 uint64) error {
	if userID1 > userID2 {
		userID1, userID2 = userID2, userID1
	}
	return r.db.Create(&models.Friend{UserId1: userID1, UserId2: userID2}).Error
}

func (r *friendRepository) DeleteFriendship(user1ID, user2ID uint64) error {
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}
	return r.db.Where("userId1 = ? and userId2 = ?", user1ID, user2ID).Delete(&models.Friend{}).Error
}

func (r *friendRepository) IsFriends(user1ID, user2ID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Friend{}).Where(
		"userId1 = ? and userId2 = ?", user1ID, user2ID).Count(&count).Error
	return count > 0, err
}

func (r *friendRepository) GetFriendCount(userID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&models.Friend{}).Where(
		"userId1 = ? or userId2 = ?", userID, userID).Count(&count).Error
	return count, err
}

func (r *friendRepository) GetFriendList(userID uint64) ([]models.User, error) {
	var res []models.User
	err := r.db.Table("friends").Select("users.*").Joins(
		"Join users On users.userId = friends.userId2").Where(
		"(userId1 = ? or userId2 = ?) and userId != ?", userID, userID, userID).Find(&res).Error
	return res, err
}

func (r *friendRepository) CreateInvitation(invitation *models.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *friendRepository) UpdateInvitationStatus(id uint64, status string) error {
	return r.db.Model(&models.Invitation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *friendRepository) GetInvitationByReceiver(receiverID uint64, status string) ([]models.
	Invitation, error) {
	var res []models.Invitation

	query := r.db.Where("reciverId = ?", receiverID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("createdAt DESC").Find(&res).Error
	return res, err
}

func (r *friendRepository) GetInvitationBySender(senderID uint64, status string) ([]models.
	Invitation, error) {
	var res []models.Invitation

	query := r.db.Where("senderId = ?", senderID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("createdAt DESC").Find(&res).Error
	return res, err
}

func (r *friendRepository) FindInvitation(senderID, receiverID uint64) (*models.Invitation, error) {
	var res models.Invitation

	query := r.db.Where("senderId = ? and receiverId = ?", senderID, receiverID)
	err := query.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *friendRepository) GetInvitationByID(id uint64) (*models.Invitation, error) {
	var res models.Invitation

	//查询符合条件的第一条记录
	//按主键查询： r.db.First(&res,id) -- equal select * from invitations where primarykey=id limit 1
	//非主键查询:  r.db.Where("target = ?",target).First(&res)
	err := r.db.First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *friendRepository) AcceptInvitation(invitationID, userID uint64) error {
	//r.db.Begin()开启的事务需要手动控制回滚和提交
	//r.db.Transaction开启的事务，当回调函数返回err时自动回滚，返回nil时自动提交
	return r.db.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation
		if err := tx.First(&invitation, invitationID).Error; err != nil {
			return err
		}

		if invitation.ReceiverID != userID {
			return errors.New("no authorize")
		}

		if invitation.Status != "pending" {
			return errors.New("invitaion was done")
		}

		if err := tx.Model(&invitation).Update("status", "accepted").Error; err != nil {
			return err
		}

		friendShip := models.Friend{
			UserId1: invitation.SenderID,
			UserId2: invitation.ReceiverID,
		}

		if err := tx.Create(&friendShip).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *friendRepository) RefuseInvitation(invitationID, userID uint64) error {
	var invitation models.Invitation

	err := r.db.First(&invitation, invitationID).Error
	if err != nil {
		return err
	}

	if invitation.ReceiverID != userID {
		return errors.New("no authorize")
	}

	if invitation.Status != "pending" {
		return errors.New("invitaion was done")
	}

	return r.db.Model(&invitation).Update("status", "refused").Error
}
