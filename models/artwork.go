package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Artwork struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"index" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	Artist       string         `gorm:"index" json:"artist"`
	AvatarURL    string         `json:"avatar_url"`
	URL          string         `gorm:"index" json:"url"` // 图片来源
	LocalPath    string         `json:"-"`                // 本地备份路径
	CDNPath      string         `json:"-"`                // CDN路径
	CDNURL       string         `json:"cdn_url"`
	ThumbnailURL string         `json:"thumbnail_url"` // 缩略图URL
	Hash         string         `json:"hash"`
	Views        int            `gorm:"default:0" json:"views"`
	Likes        int            `gorm:"default:0" json:"likes"`
	Bookmarks    int            `gorm:"default:0" json:"bookmarks"`
	Tags         string         `gorm:"type:text" json:"tags"` // JSON 字符串格式存储
	Category     string         `gorm:"index" json:"category"`
	IsPublished  bool           `gorm:"default:true;index" json:"is_published"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Artwork) TableName() string {
	return "artworks"
}

// ArtworkCreateRequest 创建请求
type ArtworkCreateRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	Artist       string   `json:"artist" binding:"required"`
	AvatarURL    string   `json:"avatar_url"`
	URL          string   `json:"url" binding:"required"`
	LocalPath    string   `json:"-"`
	CDNPath      string   `json:"-"` // CDN路径
	CDNURL       string   `json:"cdn_url"`
	Hash         string   `json:"hash"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Tags         []string `json:"tags"`
	Category     string   `json:"category"`
	IsPublished  bool     `json:"is_published"`
}

// ArtworkUpdateRequest 更新请求
type ArtworkUpdateRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Artist       string   `json:"artist"`
	AvatarURL    string   `json:"avatar_url"`
	URL          string   `json:"url"`
	LocalPath    string   `json:"local_path"`
	CDNURL       string   `json:"cdn_url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Tags         []string `json:"tags"`
	Category     string   `json:"category"`
	IsPublished  *bool    `json:"is_published"`
}

// ArtworkResponse 返回响应
type ArtworkResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Artist       string    `json:"artist"`
	AvatarURL    string    `json:"avatar_url"`
	URL          string    `json:"url"`
	LocalPath    string    `json:"local_path"`
	CDNPath      string    `json:"-"` // CDN路径
	CDNURL       string    `json:"cdn_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Views        int       `json:"views"`
	Likes        int       `json:"likes"`
	Bookmarks    int       `json:"bookmarks"`
	Tags         []string  `json:"tags"`
	Category     string    `json:"category"`
	IsPublished  bool      `json:"is_published"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (a *Artwork) ToResponse() ArtworkResponse {
	var tags []string
	if a.Tags != "" {
		json.Unmarshal([]byte(a.Tags), &tags)
	}
	// 确保 tags 不为 nil，至少是空数组
	if tags == nil {
		tags = []string{}
	}

	return ArtworkResponse{
		ID:           a.ID,
		Title:        a.Title,
		Description:  a.Description,
		Artist:       a.Artist,
		AvatarURL:    a.AvatarURL,
		URL:          a.URL,
		LocalPath:    a.LocalPath,
		CDNURL:       a.CDNURL,
		CDNPath:      a.CDNPath,
		ThumbnailURL: a.ThumbnailURL,
		Views:        a.Views,
		Likes:        a.Likes,
		Bookmarks:    a.Bookmarks,
		Tags:         tags,
		Category:     a.Category,
		IsPublished:  a.IsPublished,
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
