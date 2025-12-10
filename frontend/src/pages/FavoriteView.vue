<template>
  <div class="bg-background text-foreground py-8 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-4xl font-bold mb-2">我的收藏</h1>
        <p class="text-muted-foreground">共 {{ bookmarkedArtworks.length }} 个收藏</p>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="loading loading-spinner loading-lg text-primary"></div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="bookmarkedArtworks.length === 0" class="text-center py-20">
        <div class="mb-4">
          <span class="icon-[lucide--heart] text-6xl text-muted-foreground opacity-50"></span>
        </div>
        <p class="text-muted-foreground text-lg mb-4">暂无收藏作品</p>
        <RouterLink to="/" class="btn btn-primary"> 去浏览作品 </RouterLink>
      </div>

      <!-- 收藏列表 -->
      <div v-else class="grid gap-4 mb-8">
        <div
          v-for="artwork in paginatedArtworks"
          :key="artwork.id"
          class="group relative overflow-hidden rounded-lg bg-card shadow-sm hover:shadow-md transition-all duration-300"
        >
          <div class="flex gap-4 p-4">
            <!-- 图片 -->
            <div
              class="w-32 h-32 shrink-0 overflow-hidden rounded-lg bg-muted cursor-pointer"
              @click="navigateToDetail(artwork.id)"
            >
              <img
                :src="artwork.thumbnail_url"
                :alt="artwork.title"
                class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
              />
            </div>

            <!-- 信息 -->
            <div class="flex-1 flex flex-col justify-between min-w-0 relative">
              <div>
                <h3
                  class="font-semibold text-lg group-hover:text-primary transition-colors cursor-pointer"
                  @click="navigateToDetail(artwork.id)"
                >
                  {{ artwork.title }}
                </h3>
                <p class="text-sm text-muted-foreground mt-1">{{ artwork.artist }}</p>

                <!-- 标签 -->
                <div class="flex flex-wrap gap-2 mt-3">
                  <span
                    v-for="tag in artwork.tags.slice(0, 3)"
                    :key="tag"
                    class="text-xs bg-muted text-muted-foreground px-2 py-1 rounded"
                  >
                    #{{ tag }}
                  </span>
                </div>
              </div>

              <!-- 统计和收藏日期 -->
              <div class="flex items-center justify-between mt-3 pt-3 border-t border-base-300">
                <div class="flex items-center gap-4 text-sm text-muted-foreground">
                  <span class="flex items-center gap-1">
                    <span class="icon-[lucide--eye] text-sm"></span>
                    {{ artwork.views }}
                  </span>
                  <span class="flex items-center gap-1">
                    <span class="icon-[lucide--heart] text-sm"></span>
                    {{ artwork.likes }}
                  </span>
                </div>

                <div class="text-sm text-muted-foreground">
                  收藏于 {{ formatDate(getBookmarkTime(artwork.id)) }}
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div
              class="absolute top-4 right-4 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <button
                @click="navigateToDetail(artwork.id)"
                class="btn btn-sm btn-primary"
                title="查看详情"
              >
                <span class="icon-[lucide--eye]"></span>
              </button>
              <button
                @click="handleRemoveBookmark(artwork.id)"
                title="取消收藏"
                class="btn btn-sm btn-destructive"
                :disabled="artworkStore.loading"
              >
                <span class="icon-[lucide--heart-crack]"></span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页器 -->
      <div v-if="bookmarkedArtworks.length > 0" class="flex justify-center mb-8">
        <PageNavigator
          v-model:currentPage="currentPage"
          :total="bookmarkedArtworks.length"
          :page-size="pageSize"
        />
      </div>

      <!-- 底部操作 -->
      <div
        v-if="bookmarkedArtworks.length > 0"
        class="pt-8 border-t border-base-300 flex justify-center"
      >
        <button @click="clearAllBookmarks" class="btn btn-ghost text-error">清空所有收藏</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ref, computed, onMounted } from 'vue'
import { useArtworkStore } from '@/stores/artwork'
import PageNavigator from '@/components/PageNavigator.vue'
import type { Artwork } from '@/types'

const router = useRouter()
const artworkStore = useArtworkStore()

const loading = ref(false)
const bookmarkedArtworksList = ref<Artwork[]>([])
const currentPage = ref(1)
const pageSize = 20

// 获取已收藏的作品
const bookmarkedArtworks = computed(() => {
  return bookmarkedArtworksList.value.filter((artwork) => artworkStore.isBookmarked(artwork.id))
})

// 分页后的作品
const paginatedArtworks = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return bookmarkedArtworks.value.slice(start, end)
})

// 获取收藏时间
const getBookmarkTime = (artworkId: number): number => {
  const status = artworkStore.getBookmarkStatus(artworkId)
  return status.timestamp || Date.now()
}

// 格式化日期
const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) {
    return `今天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  }

  if (diffDays === 1) {
    return `昨天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  }

  if (diffDays < 7) {
    return `${diffDays}天前`
  }

  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// 导航到详情页
const navigateToDetail = (id: number) => {
  router.push(`/artwork/${id}`)
}

// 移除单个收藏
const handleRemoveBookmark = async (id: number) => {
  try {
    await artworkStore.toggleBookmark(id)
  } catch (err) {
    console.error('取消收藏失败:', err)
  }
}

// 清空所有收藏
const clearAllBookmarks = async () => {
  if (!confirm('确定要清空所有收藏吗？')) return
  try {
    artworkStore.clearAllBookmarkRecords()
    bookmarkedArtworksList.value = []
  } catch (err) {
    console.error('清空收藏失败:', err)
  }
}

// 获取收藏的作品ID列表
const getBookmarkedArtworkIds = (): number[] => {
  const bookmarkRecords = localStorage.getItem('artwork:bookmarks')
  if (!bookmarkRecords) return []

  try {
    const records = JSON.parse(bookmarkRecords)
    return Object.keys(records)
      .filter((key) => records[key].bookmarked)
      .map((key) => parseInt(key))
  } catch (err) {
    console.error('解析收藏记录失败:', err)
    return []
  }
}

// 初始化加载
const initLoad = async () => {
  loading.value = true
  try {
    const bookmarkedIds = getBookmarkedArtworkIds()

    if (bookmarkedIds.length === 0) {
      bookmarkedArtworksList.value = []
      return
    }

    const artworks: Artwork[] = []
    for (const id of bookmarkedIds) {
      try {
        const result = await artworkStore.fetchArtwork(id)
        if (result.data) {
          artworks.push(result.data)
        }
      } catch (err) {
        console.error(`加载作品 ${id} 失败:`, err)
      }
    }

    bookmarkedArtworksList.value = artworks
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initLoad()
})
</script>
