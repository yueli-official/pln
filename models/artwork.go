package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	FileID       string         `gorm:"uniqueIndex:idx_file_id;not null" json:"-"` // 本地备份文件ID
	URL          string         `json:"url"`                                       // CDN访问链接
	ThumbnailURL string         `json:"thumbnail_url"`                             // 缩略图URL
	Hash         string         `gorm:"index:idx_hash;not null" json:"hash"`       // 文件哈希
	PHash        int64          `gorm:"index:idx_hash;column:phash;" json:"phash"`
	Views        int            `gorm:"default:0" json:"views"`
	Likes        int            `gorm:"default:0" json:"likes"`
	Bookmarks    int            `gorm:"default:0" json:"bookmarks"`
	Tags         string         `gorm:"type:text" json:"tags"` // JSON 字符串格式存储
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Artwork) TableName() string {
	return "artworks"
}

// 创建请求
type ArtworkCreateRequest struct {
	FileID       string   `json:"file_id" binding:"required"`
	URL          string   `json:"url"`
	PHash        int64    `json:"phash"`
	Hash         string   `json:"hash"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Tags         []string `json:"tags"`
}

// ArtworkUpdateRequest 更新请求
type ArtworkUpdateRequest struct {
	URL          string   `json:"url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Tags         []string `json:"tags"`
}

// 返回响应
type ArtworkResponse struct {
	ID           uint      `json:"id"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Views        int       `json:"views"`
	Likes        int       `json:"likes"`
	Bookmarks    int       `json:"bookmarks"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 转换为响应格式
func (a *Artwork) ToResponse() ArtworkResponse {
	var tags []string
	if a.Tags != "" {
		_ = json.Unmarshal([]byte(a.Tags), &tags)
	}
	if tags == nil {
		tags = []string{}
	}

	return ArtworkResponse{
		ID:           a.ID,
		URL:          a.URL,
		ThumbnailURL: a.ThumbnailURL,
		Views:        a.Views,
		Likes:        a.Likes,
		Bookmarks:    a.Bookmarks,
		Tags:         tags,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// SetTags 设置 tags（将 []string 转换为 JSON 字符串）
func (a *Artwork) SetTags(tags []string) error {
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	a.Tags = string(data)
	return nil
}
