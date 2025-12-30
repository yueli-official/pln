// ==================== 文件服务相关 ====================

export type FileType = 'image' | 'video' | 'document' | 'audio' | 'other'

export interface FileMetadata {
  hash?: string

  // 图片
  width?: number
  height?: number
  orientation?: string
  color_space?: string

  // 视频
  duration?: number
  frame_rate?: number
  video_codec?: string
  resolution?: string

  // 音频
  bit_rate?: number
  sample_rate?: number

  // 通用
  mime_type?: string
  custom?: Record<string, any>
}

export interface ResourceVariant {
  type: string
  path: string
  access_url: string
  size?: number
  created_at: string
  extra: Record<string, any>
}

export interface FileInfo {
  // 基本信息
  name: string
  original_name: string
  size: number

  // 路径
  path: string

  // 访问信息
  access_url: string
  storage_type: string

  // 类型
  file_type: FileType

  // 扩展
  metadata?: FileMetadata
  extra: Record<string, any>
  variants?: ResourceVariant[]

  // 时间
  upload_time: number
  modified_at: number

  // 标识符
  file_id: string
  id: string
}

export interface DeleteFileResponse {
  path: string
  deleted: boolean
}
// ==================== 作品相关 ====================
// ==================== 作品相关 ====================

export interface Artwork {
  id: number
  url: string
  thumbnail_url: string
  views: number
  likes: number
  bookmarks: number
  tags: string[]
  created_at: string
  updated_at: string
}

export interface ArtworkCreateRequest {
  file_id: string
  url?: string
  hash?: string
  thumbnail_url?: string
  tags?: string[]
}

export interface ArtworkUpdateRequest {
  url?: string
  thumbnail_url?: string
  tags?: string[]
}

export type ArtworkResponse = Artwork

// ==================== API 通用 ====================

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

// ==================== 分页 ====================

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface PaginationParams {
  page?: number
  page_size?: number
  sort?: string
}

export interface PaginationMeta {
  total: number
  page: number
  page_size: number
  total_pages: number
}
