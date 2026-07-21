package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	Common
	UserID    string         `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	AnimeID   string         `gorm:"column:anime_id;type:text;index;not null" json:"animeId"`
	Rating    int            `gorm:"column:rating;not null" json:"rating"`
	Title     string         `gorm:"column:title;type:text" json:"title"`
	Content   string         `gorm:"column:content;type:text" json:"content"`
	Spoiler   bool           `gorm:"column:spoiler;default:false" json:"spoiler"`
	Helpful   int            `gorm:"column:helpful;default:0" json:"helpful"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Review) TableName() string { return "reviews" }

type Comment struct {
	Common
	UserID    string         `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	AnimeID   *string        `gorm:"column:anime_id;type:text;index" json:"animeId,omitempty"`
	EpisodeID *string        `gorm:"column:episode_id;type:text;index" json:"episodeId,omitempty"`
	ReviewID  *string        `gorm:"column:review_id;type:text;index" json:"reviewId,omitempty"`
	ParentID  *string        `gorm:"column:parent_id;type:text;index" json:"parentId,omitempty"`
	Content   string         `gorm:"column:content;type:text;not null" json:"content"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Comment) TableName() string { return "comments" }

type Report struct {
	Common
	UserID      string     `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	TargetType  string     `gorm:"column:target_type;type:text;not null" json:"targetType"`
	TargetID    string     `gorm:"column:target_id;type:text;not null" json:"targetId"`
	Reason      string     `gorm:"column:reason;type:text;not null" json:"reason"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Status      string     `gorm:"column:status;type:text;not null;default:'pending'" json:"status"`
	ResolvedBy  *string    `gorm:"column:resolved_by;type:text" json:"resolvedBy,omitempty"`
	ResolvedAt  *time.Time `gorm:"column:resolved_at" json:"resolvedAt,omitempty"`
}

func (Report) TableName() string { return "reports" }

type Watchlist struct {
	Common
	UserID string `gorm:"column:user_id;type:text;index;not null" json:"userId"`
	Name   string `gorm:"column:name;type:text;not null" json:"name"`
	Type   string `gorm:"column:type;type:text;not null;default:'custom'" json:"type"`
}

func (Watchlist) TableName() string { return "watchlists" }

type WatchlistItem struct {
	Common
	WatchlistID string `gorm:"column:watchlist_id;type:text;index;not null" json:"watchlistId"`
	AnimeID     string `gorm:"column:anime_id;type:text;index;not null" json:"animeId"`
}

func (WatchlistItem) TableName() string { return "watchlist_items" }
