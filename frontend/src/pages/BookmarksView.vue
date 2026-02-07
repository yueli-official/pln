<template>
  <div class="bg-background text-foreground py-12 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-10">
        <h1 class="text-4xl font-bold mb-2">我的收藏</h1>
        <p class="text-muted-foreground">共 {{ bookmarkedArtworks.length }} 个收藏</p>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="flex justify-center items-center py-20">
        <div class="relative w-16 h-16">
          <div class="absolute inset-0 rounded-full border-4 border-border"></div>
          <div
            class="absolute inset-0 rounded-full border-4 border-transparent border-t-primary animate-spin"
          ></div>
        </div>
      </div>

      <!-- 空状态 -->
      <div
        v-else-if="bookmarkedArtworks.length === 0"
        class="flex flex-col items-center justify-center py-20"
      >
        <span class="icon-[lucide--heart] text-6xl text-muted-foreground/50 mb-4"></span>
        <p class="text-muted-foreground text-lg mb-6">暂无收藏作品</p>
        <RouterLink
          to="/browse"
          class="px-6 py-2 rounded-full bg-primary text-white hover:shadow-lg transition-all"
        >
          去浏览作品
        </RouterLink>
      </div>

      <!-- 收藏网格 -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5 mb-12">
        <div
          v-for="artwork in paginatedArtworks"
          :key="artwork.id"
          class="group flex flex-col rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm overflow-hidden hover:border-primary/50 hover:shadow-lg transition-all duration-300"
        >
          <!-- 图片部分 -->
          <div
            class="w-full aspect-square overflow-hidden bg-muted cursor-pointer relative"
            @click="navigateToDetail(artwork.id)"
          >
            <img
              :src="artwork.thumbnail_url"
              :alt="'普拉娜'"
              class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
            />
            <!-- 悬停覆盖层 -->
            <div
              class="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors duration-300 flex items-center justify-center opacity-0 group-hover:opacity-100"
            >
              <span class="icon-[lucide--eye] text-4xl text-white"></span>
            </div>
          </div>

          <!-- 信息部分 -->
          <div class="flex-1 flex flex-col p-4">
            <!-- 标签 -->
            <div class="flex flex-wrap gap-2 mb-3">
              <span
                v-for="tag in artwork.tags.slice(0, 2)"
                :key="tag"
                class="text-xs bg-primary/10 text-primary px-2 py-1 rounded-full"
              >
                #{{ tag }}
              </span>
            </div>

            <!-- 统计信息 -->
            <div
              class="flex items-center gap-4 text-sm text-muted-foreground mb-3 py-2 border-y border-border/30"
            >
              <span class="flex items-center gap-1">
                <span class="icon-[lucide--eye]"></span>
                {{ artwork.views }}
              </span>
              <span class="flex items-center gap-1">
                <span class="icon-[lucide--heart] text-red-500"></span>
                {{ artwork.likes }}
              </span>
            </div>

            <!-- 收藏日期 -->
            <div class="text-xs text-muted-foreground mb-4">
              💾 {{ formatDate(getBookmarkTime(artwork.id)) }}
            </div>

            <!-- 操作按钮 -->
            <div class="flex gap-2 mt-auto">
              <button
                @click="navigateToDetail(artwork.id)"
                class="flex-1 px-3 py-2 rounded-lg bg-primary/10 text-primary hover:bg-primary/20 transition-colors text-sm font-medium"
              >
                <span class="icon-[lucide--eye] mr-1"></span>查看
              </button>
              <button
                @click="handleRemoveBookmark(artwork.id)"
                class="flex-1 px-3 py-2 rounded-lg bg-error/10 text-error hover:bg-error/20 transition-colors text-sm font-medium"
              >
                <span class="icon-[lucide--heart-crack] mr-1"></span>取消
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
        class="pt-8 border-t border-border/30 flex justify-center"
      >
        <button
          @click="clearAllBookmarks"
          class="px-6 py-2 rounded-full bg-error/10 text-error hover:bg-error/20 transition-colors"
        >
          清空所有收藏
        </button>
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

const loading = ref<boolean>(false)
const bookmarkedArtworksList = ref<Artwork[]>([])
const currentPage = ref<number>(1)
const pageSize: number = 20

const bookmarkedArtworks = computed<Artwork[]>(() => {
  return bookmarkedArtworksList.value.filter((artwork) => artworkStore.isBookmarked(artwork.id))
})

const paginatedArtworks = computed<Artwork[]>(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return bookmarkedArtworks.value.slice(start, end)
})

const getBookmarkTime = (artworkId: number): number => {
  const status = artworkStore.getBookmarkStatus(artworkId)
  return status.timestamp || Date.now()
}

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

const navigateToDetail = (id: number): void => {
  router.push(`/artwork/${id}`)
}

const handleRemoveBookmark = async (id: number): Promise<void> => {
  try {
    await artworkStore.toggleBookmark(id)
  } catch (err) {
    console.error('取消收藏失败:', err)
  }
}

const clearAllBookmarks = async (): Promise<void> => {
  if (!confirm('确定要清空所有收藏吗？')) return
  try {
    artworkStore.clearAllBookmarkRecords()
    bookmarkedArtworksList.value = []
  } catch (err) {
    console.error('清空收藏失败:', err)
  }
}

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

const initLoad = async (): Promise<void> => {
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
