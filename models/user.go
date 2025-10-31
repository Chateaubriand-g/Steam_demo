package models

import "time"

type User struct {
	UserID    uint64    `json:"userId" gorm:"primarykey;autoIncrement"`
	Email     string    `json:"email" gorm:"size:255;not null;uniqueIndex"`
	UserName  string    `json:"userName" gorm:"size:50;uniqueIndex"`
	PassWord  string    `json:"-" gorm:"size:18"` //json:"-"指定序列化时忽略。存储到数据库时需加密处理
	NickName  string    `json:"nickName" gorm:"size:50;not null"`
	Avatar    string    `json:"avatar" gorm:"size:255"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdateAt  time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}

type UserDto struct {
	UserID   uint64 `json:"userId"`
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

type UserInfoResponseDto struct {
	UserID   uint64 `json:"userId"`
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

type LoginRequestDto struct {
	UserName string `json:"userName" binding:"required"` //bindg用于gin框架参数校验，缺少时直接返回殂错误
	Password string `json:"password" binding:"required"`
}

type JoinRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"userName" binding:"required"`
	PassWord string `json:"passWord" binding:"required"`
}

type LoginResponseDto struct {
	Token  string `json:"token"`
	UserID uint64 `json:"userId"`
}

type InvitationRequestDto struct {
	ReceiverID uint64 `json:"receiverId" binding:"required,min=8,max=16"`
}

type FriendRequestDto struct {
	FriendID uint64 `json:"friendId" binding:"required"`
}

type WishListRequestDto struct {
	AppID uint64 `json:"appId binding:"required"`
}
