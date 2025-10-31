package models

type App struct {
	AppId        uint64  `json:"appId gorm:"primarykey"`
	Name         string  `json:"name" gorm:"size:255;not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Price        float64 `json:"price" gorm:"type:decimal(10,2)"`
	Discount     float64 `json:"discount" gorm:"type:decimal(5,2)"`
	ReleaseDate  string  `josn:"releaseDate" gorm:"size:50"`
	Developer    string  `json:"developer" gorm:"size:255"`
	Publisher    string  `json:"publisher" gorm:"size:255"`
	ImageURL     string  `json:"imageURL" gorm:"size:500"`
	Tags         string  `json:"tags" gorm:"type:text"`
	PositiveRate int     `json:"positiveRate" gorm:"default:0"`
}

type AppDto struct {
	AppId        uint64  `json:"appId"`
	Name         string  `json:"name"`
	Price        float64 `json:"price" gorm:"type:decimal(10,2)"`
	Discount     float64 `json:"discount" gorm:"type:decimal(5,2)"`
	ImageURL     string  `json:"imageURL" gorm:"size:500"`
	PositiveRate int     `json:"positiveRate" gorm:"default:0"`
}

type RecommendationDto struct {
	AppId    uint64  `json:"appId"`
	Name     string  `json:"name"`
	Price    float64 `json:"price" gorm:"type:decimal(10,2)"`
	ImageURL string  `json:"imageURL" gorm:"size:500"`
}

type SpecialDto struct {
	AppId         uint64  `json:"appId"`
	Name          string  `json:"name"`
	ImageURL      string  `json:"imageURL" gorm:"size:500"`
	OriginalPrice float64 `json:"originalPrice"`
	CurrentPrice  float64 `json:"currentPrice"`
	Discount      int     `json:"discount"`
}
