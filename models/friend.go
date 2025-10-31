package models

import "time"

type Friend struct {
	UserId1   uint64    `json:"userId1" gorm:"primarykey"`
	UserId2   uint64    `json:"userId2" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type Invitation struct {
	ID         uint64    `json:"id" gorm:"primarykey;autoIncrement"`
	SenderID   uint64    `json:"senderId" gorm:"index"`
	ReceiverID uint64    `json:"recriverId" gorm:"index"`
	Status     string    `json:"status" gorm:"size:20;default:'pending'"` //pending,accepted,refused
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`         //自动填充当前时间
}
