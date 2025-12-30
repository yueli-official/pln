import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useStorage } from '@vueuse/core'
import { useApi } from '@/composables/useApi'
import type { Artwork, ArtworkCreateRequest, ArtworkUpdateRequest } from '@/types'

interface StoredLike {
  artworkId: number
  liked: boolean
  timestamp: number
}

interface StoredBookmark {
  artworkId: number
  bookmarked: boolean
  timestamp: number
}

export const useArtworkStore = defineStore('artwork', () => {
  // API
  const {
    getArtworks,
    getArtwork,
    createArtwork,
    updateArtwork,
    deleteArtwork,
    incrementBookmarks,
    decrementBookmarks,
    incrementLikes,
    decrementLikes,
    getRandomArtworks,
    loading,
    error,
  } = useApi()

  // ==================== 本地存储 ====================

  // 点赞记录存储在本地
  const likeRecords = useStorage<Record<number, StoredLike>>('artwork:likes', {}, localStorage, {
    mergeDefaults: true,
  })

  // 收藏记录存储在本地
  const bookmarkRecords = useStorage<Record<number, StoredBookmark>>(
    'artwork:bookmarks',
    {},
    localStorage,
    { mergeDefaults: true },
  )

  // ==================== 状态 ====================

  const artworks = ref<Artwork[]>([])
  const currentArtwork = ref<Artwork | null>(null)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const total = ref(0)

  // ==================== 计算属性 ====================

  const hasMoreArtworks = computed(() => artworks.value.length < total.value)

  // ==================== 点赞相关方法 ====================

  /**
   * 检查作品是否已点赞
   */
  const isLiked = (artworkId: number): boolean => {
    return likeRecords.value[artworkId]?.liked ?? false
  }

  /**
   * 获取点赞状态
   */
  const getLikeStatus = (artworkId: number) => {
    return likeRecords.value[artworkId] || { artworkId, liked: false, timestamp: 0 }
  }

  /**
   * 点赞作品
   */
  /**
   * 点赞作品
   */
  const toggleLike = async (artworkId: number) => {
    try {
      const currentStatus = isLiked(artworkId)

      // 调用对应的 API
      if (currentStatus) {
        // 已点赞，取消点赞
        await decrementLikes(artworkId)
      } else {
        // 未点赞，增加点赞
        await incrementLikes(artworkId)
      }

      // 更新本地记录
      likeRecords.value[artworkId] = {
        artworkId,
        liked: !currentStatus,
        timestamp: Date.now(),
      }

      // 更新作品列表中的数据
      updateArtworkLikeStatus(artworkId, !currentStatus)

      return !currentStatus
    } catch (err) {
      console.error('点赞失败:', err)
      throw err
    }
  }

  /**
   * 清除点赞记录
   */
  const clearLikeRecord = (artworkId: number) => {
    delete likeRecords.value[artworkId]
  }

  /**
   * 清除所有点赞记录
   */
  const clearAllLikeRecords = () => {
    likeRecords.value = {}
  }

  // ==================== 收藏相关方法 ====================

  /**
   * 检查作品是否已收藏
   */
  const isBookmarked = (artworkId: number): boolean => {
    return bookmarkRecords.value[artworkId]?.bookmarked ?? false
  }

  /**
   * 获取收藏状态
   */
  const getBookmarkStatus = (artworkId: number) => {
    return bookmarkRecords.value[artworkId] || { artworkId, bookmarked: false, timestamp: 0 }
  }

  /**
   * 收藏作品
   */
  const toggleBookmark = async (artworkId: number) => {
    try {
      const currentStatus = isBookmarked(artworkId)

      // 调用对应的 API
      if (currentStatus) {
        // 已收藏，取消收藏
        await decrementBookmarks(artworkId)
      } else {
        // 未收藏，增加收藏
        await incrementBookmarks(artworkId)
      }

      // 更新本地记录
      bookmarkRecords.value[artworkId] = {
        artworkId,
        bookmarked: !currentStatus,
        timestamp: Date.now(),
      }

      return !currentStatus
    } catch (err) {
      console.error('收藏失败:', err)
      throw err
    }
  }

  /**
   * 清除收藏记录
   */
  const clearBookmarkRecord = (artworkId: number) => {
    delete bookmarkRecords.value[artworkId]
  }

  /**
   * 清除所有收藏记录
   */
  const clearAllBookmarkRecords = () => {
    bookmarkRecords.value = {}
  }

  // ==================== 作品列表相关方法 ====================

  /**
   * 获取作品列表
   */
  const fetchArtworks = async (page = 1, pageSize_ = 20, filters?: Record<string, any>) => {
    try {
      const result = await getArtworks(page, pageSize_, filters)
      artworks.value = result.data.list || []
      total.value = result.data.total || 0
      currentPage.value = page
      pageSize.value = pageSize_

      return result
    } catch (err) {
      console.error('获取作品列表失败:', err)
      throw err
    }
  }

  /**
   * 获取随机作品
   */
  const fetchRandomArtworks = async (limit = 10, filters?: Record<string, any>) => {
    try {
      const result = await getRandomArtworks(limit, filters)
      artworks.value = result.data || []

      return result
    } catch (err) {
      console.error('获取随机作品失败:', err)
      throw err
    }
  }

  /**
   * 加载更多作品
   */
  const loadMoreArtworks = async (filters?: Record<string, any>) => {
    if (!hasMoreArtworks.value) return

    try {
      const nextPage = currentPage.value + 1
      const result = await getArtworks(nextPage, pageSize.value, filters)
      artworks.value.push(...(result.data.list || []))
      total.value = result.data.total || 0
      currentPage.value = nextPage

      return result
    } catch (err) {
      console.error('加载更多作品失败:', err)
      throw err
    }
  }

  /**
   * 获取单个作品
   */
  const fetchArtwork = async (id: number) => {
    try {
      const result = await getArtwork(id)
      currentArtwork.value = result.data

      return result
    } catch (err) {
      console.error('获取作品详情失败:', err)
      throw err
    }
  }

  // ==================== 作品创建/更新/删除 ====================

  /**
   * 创建作品
   */
  const create = async (data: ArtworkCreateRequest) => {
    try {
      const result = await createArtwork(data)
      return result
    } catch (err) {
      console.error('创建作品失败:', err)
      throw err
    }
  }

  /**
   * 更新作品
   */
  const update = async (id: number, data: ArtworkUpdateRequest) => {
    try {
      const result = await updateArtwork(id, data)
      currentArtwork.value = result
      return result
    } catch (err) {
      console.error('更新作品失败:', err)
      throw err
    }
  }

  /**
   * 删除作品
   */
  const remove = async (id: number) => {
    try {
      await deleteArtwork(id)
      // 从列表中移除
      artworks.value = artworks.value.filter((a) => a.id !== id)
      // 清除本地记录
      clearLikeRecord(id)
      clearBookmarkRecord(id)
      return true
    } catch (err) {
      console.error('删除作品失败:', err)
      throw err
    }
  }

  // ==================== 辅助方法 ====================

  /**
   * 更新作品的点赞状态
   */
  const updateArtworkLikeStatus = (artworkId: number, liked: boolean) => {
    const artwork = artworks.value.find((a) => a.id === artworkId)
    if (artwork) {
      if (liked) {
        artwork.likes = (artwork.likes || 0) + 1
      } else {
        artwork.likes = Math.max(0, (artwork.likes || 1) - 1)
      }
    }

    if (currentArtwork.value?.id === artworkId) {
      if (liked) {
        currentArtwork.value.likes = (currentArtwork.value.likes || 0) + 1
      } else {
        currentArtwork.value.likes = Math.max(0, (currentArtwork.value.likes || 1) - 1)
      }
    }
  }

  /**
   * 重置 store
   */
  const reset = () => {
    artworks.value = []
    currentArtwork.value = null
    currentPage.value = 1
    total.value = 0
  }

  return {
    // 状态
    artworks,
    currentArtwork,
    currentPage,
    pageSize,
    total,
    loading,
    error,

    // 计算属性
    hasMoreArtworks,

    // 点赞
    isLiked,
    getLikeStatus,
    toggleLike,
    clearLikeRecord,
    clearAllLikeRecords,

    // 收藏
    isBookmarked,
    getBookmarkStatus,
    toggleBookmark,
    clearBookmarkRecord,
    clearAllBookmarkRecords,

    // 作品列表
    fetchArtworks,
    loadMoreArtworks,
    fetchArtwork,
    fetchRandomArtworks,

    // 作品操作
    create,
    update,
    remove,

    // 工具方法
    reset,
  }
})
