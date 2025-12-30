package models

import "time"

// FileType 文件类型枚举
type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeAudio    FileType = "audio"
	FileTypeOther    FileType = "other"
)

// ============ 核心模型 ============

// FileInfo 统一的文件信息（对外返回）
type FileInfo struct {
	// 基本信息
	Name         string `json:"name"`          // 文件名（不含路径）
	OriginalName string `json:"original_name"` // 原始文件名
	Size         int64  `json:"size"`          // 文件大小（字节）

	// 路径信息
	Path string `json:"path"` // 逻辑路径（用于生成访问URL）

	// 访问信息
	AccessURL   string `json:"access_url"`   // 文件访问URL
	StorageType string `json:"storage_type"` // 存储类型（local/smms/aliyun等）

	// 文件属性
	FileType FileType `json:"file_type"` // 文件类型（image/video/document/audio/other）

	// 元数据（用于扩展不同文件类型的信息）
	Metadata FileMetadata      `json:"metadata,omitempty"` // 图片宽高、视频时长等
	Extra    map[string]any    `json:"extra"`
	Variants []ResourceVariant `json:"variants,omitempty"` // 缩略图、预览等变体

	// 时间信息
	UploadTime int64 `json:"upload_time"` // 上传时间
	ModifiedAt int64 `json:"modified_at"` // 修改时间

	// ===== 内部字段（不对外返回） =====
	StoragePath          string `json:"-"`
	ThumbnailStoragePath string `json:"-"`
	// ===== 标识符 =====
	FileID string `json:"file_id"` // 文件ID（来自第三方存储服务）
	ID     string `json:"id"`      // 本地数据库ID
}

// FileMetadata 文件元数据
type FileMetadata struct {
	Hash string `json:"hash,omitempty"`

	// 图片特定
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Orientation string `json:"orientation,omitempty"`
	ColorSpace  string `json:"color_space,omitempty"`

	// 视频特定
	Duration   int     `json:"duration,omitempty"` // 秒
	FrameRate  float64 `json:"frame_rate,omitempty"`
	VideoCodec string  `json:"video_codec,omitempty"`
	Resolution string  `json:"resolution,omitempty"`

	// 音频特定
	BitRate    int `json:"bit_rate,omitempty"`
	SampleRate int `json:"sample_rate,omitempty"`

	// 通用
	MimeType string         `json:"mime_type,omitempty"`
	Custom   map[string]any `json:"custom,omitempty"`
}

// ResourceVariant 资源变体（缩略图、预览等）
type ResourceVariant struct {
	Type      string         `json:"type"`       // thumbnail, preview, compressed 等
	Path      string         `json:"path"`       // 逻辑路径
	AccessURL string         `json:"access_url"` // 访问URL
	Size      int64          `json:"size,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	Extra     map[string]any `json:"extra"`
}

type DeleteFileResponse struct {
	Path    string `json:"path"`    // 被删除文件路径
	Deleted bool   `json:"deleted"` // 是否成功删除
}
