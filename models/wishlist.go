package models

import "time"

type WishlistItem struct {
	UserID    uint64    `json:"userId" gorm:"primarykey"`
	AppID     uint64    `json:"appId" gorm:"primarykey"`
	SortOrder int64     `json:"sortOrder" gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

type WishlistRequestDto struct {
	AppID uint64 `json:"appId"`
}

type WishlistItemDto struct {
	AppID        uint64  `json:"appId"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	ImageURL     string  `json:"imageUrl"`
	PositiveRate int     `json:"positiveRate"`
	SortOrder    int64   `json:"sortOrder"`
}

// 实现排序时要用事务，要么全部成功，要么全部失败
type SortWishlistRequestDto struct {
	Items []SortItem `json:"items" binding:"required"`
}

type SortItem struct {
	AppID     uint64 `json:"appId" binding:"required"`
	SortOrder int    `json:"sortOrder" binding:"required"`
}
