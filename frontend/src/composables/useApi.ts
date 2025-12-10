import axios, { type AxiosInstance, AxiosError } from 'axios'
import { ref } from 'vue'
import type {
  Artwork,
  ArtworkCreateRequest,
  ArtworkUpdateRequest,
  ArtworkResponse,
  FileInfo,
  DeleteFileResponse,
  ApiResponse,
  ApiError,
  UploadResponse,
  PageData,
} from '@/types'

// ==================== Composable ====================

export function useApi() {
  // API 配置
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:9000/api/v1'
  const API_KEY = localStorage.getItem('api_key') || ''

  // 状态
  const loading = ref(false)
  const error = ref<ApiError | null>(null)

  // 创建 axios 实例
  const instance: AxiosInstance = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // 请求拦截器
  instance.interceptors.request.use(
    (config) => {
      if (API_KEY && API_KEY.trim()) {
        config.headers['Authorization'] = `X-API-Key ${API_KEY.trim()}`
      }
      return config
    },
    (err) => Promise.reject(err),
  )

  // 响应拦截器
  instance.interceptors.response.use(
    (response) => response,
    (err: AxiosError<any>) => {
      error.value = {
        message: err.response?.data?.message || err.message || '请求失败',
        code: err.code,
      }
      return Promise.reject(err)
    },
  )

  // ==================== 作品相关 API ====================

  /**
   * 获取作品列表
   */
  const getArtworks = async (page = 1, pageSize = 20, filters?: Record<string, any>) => {
    loading.value = true
    error.value = null
    try {
      const params = {
        page,
        page_size: pageSize,
        ...filters,
      }
      const response = await instance.get<ApiResponse<PageData<ArtworkResponse[]>>>('/artworks', {
        params,
      })
      return response.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取随机作品
   */
  const getRandomArtworks = async (limit = 10, filters?: Record<string, any>) => {
    loading.value = true
    error.value = null
    try {
      const params = {
        limit,
        ...filters,
      }
      const response = await instance.get<ApiResponse<ArtworkResponse[]>>('/artworks/random', {
        params,
      })
      return response.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取单个作品
   */
  const getArtwork = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.get<ApiResponse<ArtworkResponse>>(`/artworks/${id}`)
      return response.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 按分类获取作品
   */
  const getArtworksByCategory = async (category: string, page = 1, pageSize = 20) => {
    loading.value = true
    error.value = null
    try {
      const params = { page, page_size: pageSize }
      const response = await instance.get<ApiResponse<ArtworkResponse[]>>(
        `/artworks/category/${category}`,
        { params },
      )
      return response.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 创建作品
   */
  const createArtwork = async (data: ArtworkCreateRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.post<ApiResponse<ArtworkResponse>>('/artworks', data)
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 更新作品
   */
  const updateArtwork = async (id: number, data: ArtworkUpdateRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.put<ApiResponse<ArtworkResponse>>(`/artworks/${id}`, data)
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 删除作品
   */
  const deleteArtwork = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await instance.delete(`/artworks/${id}`)
      return true
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 增加点赞
   */
  const incrementLikes = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.post<ApiResponse<{ likes: number }>>(`/artworks/${id}/like`)
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 取消点赞
   */
  const decrementLikes = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.post<ApiResponse<{ likes: number }>>(`/artworks/${id}/unlike`)
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 增加收藏
   */
  const incrementBookmarks = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.post<ApiResponse<{ bookmarks: number }>>(
        `/artworks/${id}/bookmark`,
      )
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 取消收藏
   */
  const decrementBookmarks = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.post<ApiResponse<{ bookmarks: number }>>(
        `/artworks/${id}/unbookmark`,
      )
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  // ==================== 文件相关 API ====================

  /**
   * 上传文件
   */
  const uploadFile = async (file: File, onProgress?: (progress: number) => void) => {
    loading.value = true
    error.value = null
    try {
      const formData = new FormData()
      formData.append('file', file)

      const response = await instance.post<ApiResponse<FileInfo>>('/files', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          const progress = Math.round((progressEvent.loaded / (progressEvent.total || 1)) * 100)
          onProgress?.(progress)
        },
      })
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 删除文件
   */
  const deleteFile = async (path: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.delete<ApiResponse<DeleteFileResponse>>('/files', {
        params: { path },
      })
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * 获取文件信息
   */
  const getFileInfo = async (path: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await instance.get<ApiResponse<FileInfo>>('/files/info', {
        params: { path },
      })
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  // ==================== 上传和创建作品相关 API ====================

  /**
   * 上传文件并创建作品
   */
  const uploadAndCreateArtwork = async (
    file: File,
    artworkData: {
      title: string
      artist: string
      description?: string
      category?: string
      avatar_url?: string
      tags?: string[]
    },
    onProgress?: (progress: number) => void,
  ) => {
    loading.value = true
    error.value = null
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('title', artworkData.title)
      formData.append('artist', artworkData.artist)

      if (artworkData.description) {
        formData.append('description', artworkData.description)
      }
      if (artworkData.category) {
        formData.append('category', artworkData.category)
      }
      if (artworkData.avatar_url) {
        formData.append('avatar_url', artworkData.avatar_url)
      }
      if (artworkData.tags?.length) {
        artworkData.tags.forEach((tag) => {
          formData.append('tags[]', tag)
        })
      }

      const response = await instance.post<ApiResponse<UploadResponse>>(
        '/artworks/upload',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          onUploadProgress: (progressEvent) => {
            const progress = Math.round((progressEvent.loaded / (progressEvent.total || 1)) * 100)
            onProgress?.(progress)
          },
        },
      )
      return response.data.data
    } catch (err) {
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    // 作品 API
    getArtworks,
    getRandomArtworks,
    getArtwork,
    getArtworksByCategory,
    createArtwork,
    updateArtwork,
    deleteArtwork,
    incrementBookmarks,
    decrementBookmarks,
    incrementLikes,
    decrementLikes,
    // 文件 API
    uploadFile,
    deleteFile,
    getFileInfo,
    // 组合 API
    uploadAndCreateArtwork,
  }
}
