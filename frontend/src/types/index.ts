// ==================== 类型定义 ====================

// ==================== 文件服务相关 ====================

export interface FileMetadata {
  hash?: string
  width?: number
  height?: number
  orientation?: string
  color_space?: string
  duration?: number
  frame_rate?: number
  video_codec?: string
  resolution?: string
  bit_rate?: number
  sample_rate?: number
  mime_type?: string
  custom?: Record<string, any>
}

export interface ResourceVariant {
  type: string // thumbnail, preview, compressed 等
  path: string
  access_url: string
  size?: number
  created_at: string
  extra: Record<string, any>
}

export interface FileInfo {
  name: string
  original_name: string
  size: number
  path: string
  folder_path: string
  access_url: string
  storage_type: string
  file_type: 'image' | 'video' | 'document' | 'audio' | 'other'
  metadata?: FileMetadata
  extra: Record<string, any>
  variants?: ResourceVariant[]
  upload_time: number
  modified_at: number
  space?: string
}

export interface DeleteFileResponse {
  path: string
  deleted: boolean
}

// ==================== 作品相关 ====================

export interface Artwork {
  id: number
  title: string
  description: string
  artist: string
  avatar_url: string
  url: string // 图片URL
  local_path: string
  cdn_url: string
  thumbnail_url: string // 缩略图URL
  views: number
  likes: number
  bookmarks: number
  tags: string[]
  category: string
  is_published: boolean
  created_at: string
  updated_at: string
}

export interface ArtworkCreateRequest {
  title: string
  description?: string
  artist: string
  avatar_url?: string
  url: string
  local_path?: string
  cdn_url?: string
  thumbnail_url?: string
  tags?: string[]
  category?: string
  is_published?: boolean
}

export interface ArtworkUpdateRequest {
  title?: string
  description?: string
  artist?: string
  avatar_url?: string
  url?: string
  local_path?: string
  cdn_url?: string
  thumbnail_url?: string
  tags?: string[]
  category?: string
  is_published?: boolean
  is_bookmarked?: boolean
}

export interface ArtworkResponse {
  id: number
  title: string
  description: string
  artist: string
  avatar_url: string
  url: string // 对外返回时使用 image 字段名
  local_path: string
  cdn_url: string
  thumbnail_url: string
  views: number
  likes: number
  bookmarks: number
  tags: string[]
  category: string
  is_published: boolean
  is_bookmarked: boolean
  created_at: string
  updated_at: string
}

// ==================== API 响应 ====================

export interface ApiResponse<T> {
  code?: number
  message?: string
  data: T
  request_id?: string
  timestamp?: number
}

export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
}

export interface UploadResponse {
  file: FileInfo
  artwork?: ArtworkResponse
}

// ==================== 分页相关 ====================

export interface PageData<T> {
  list: T
  total: number
  page: number
  pageSize: number
}

export interface PaginationParams {
  page?: number
  page_size?: number
  sort?: string // 排序字段，如 "-created_at" 表示按创建时间倒序
}

export interface PaginationMeta {
  total: number
  page: number
  page_size: number
  total_pages: number
}
