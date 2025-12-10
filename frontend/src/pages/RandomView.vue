<template>
  <div class="bg-background text-foreground py-8 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题和刷新按钮 -->
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-4xl font-bold mb-2">随机作品</h1>
          <p class="text-muted-foreground">发现新的艺术作品</p>
        </div>
        <button
          @click="refreshArtworks"
          :disabled="artworkStore.loading"
          class="btn btn-primary btn-lg gap-2"
        >
          <span class="icon-[lucide--shuffle] text-lg"></span>
          {{ artworkStore.loading ? '加载中...' : '换一批' }}
        </button>
      </div>

      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="alert alert-error mb-6">
        {{ artworkStore.error.message }}
      </div>

      <!-- 加载中 -->
      <div
        v-if="artworkStore.loading && artworkStore.artworks.length === 0"
        class="flex justify-center items-center py-20"
      >
        <div class="loading loading-spinner loading-lg text-primary"></div>
      </div>

      <!-- 作品网格 -->
      <div
        v-if="artworkStore.artworks.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4 auto-rows-max"
      >
        <div
          v-for="artwork in artworkStore.artworks"
          :key="artwork.id"
          class="group relative overflow-hidden rounded-lg bg-card shadow-sm hover:shadow-md transition-all duration-300"
        >
          <!-- 图片容器 -->
          <div
            class="relative w-full overflow-hidden bg-muted aspect-square block cursor-pointer"
            @click="navigateToDetail(artwork.id)"
          >
            <img
              :src="artwork.thumbnail_url"
              :alt="artwork.title"
              class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
            />

            <!-- 悬停遮罩层 -->
            <div
              class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors duration-300 opacity-0 group-hover:opacity-100 flex flex-col justify-between p-3"
            >
              <!-- 顶部：标题和艺术家 -->
              <div>
                <h3 class="font-semibold text-sm truncate text-white">{{ artwork.title }}</h3>
                <p class="text-xs text-gray-300 mt-1">{{ artwork.artist }}</p>
              </div>

              <!-- 底部：用户信息和操作 -->
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <img
                    :src="artwork.avatar_url || defaultAvatar"
                    :alt="artwork.artist"
                    class="w-6 h-6 rounded-full border border-white"
                  />
                </div>

                <!-- 收藏按钮 -->
                <button
                  @click.stop="handleBookmark(artwork.id)"
                  :disabled="artworkStore.loading"
                  class="flex items-center justify-center transition-colors"
                  :class="
                    artworkStore.isBookmarked(artwork.id)
                      ? 'text-red-500'
                      : 'text-white hover:text-red-500'
                  "
                >
                  <span class="icon-[lucide--heart] text-2xl fill-current"></span>
                </button>
              </div>
            </div>

            <!-- 右上角统计信息 -->
            <div
              class="absolute top-2 right-2 flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity text-white text-xs"
            >
              <div class="bg-black/60 px-2 py-1 rounded flex items-center gap-1">
                <span class="icon-[lucide--eye] text-sm"></span>
                <span>{{ artwork.views }}</span>
              </div>
              <div class="bg-black/60 px-2 py-1 rounded flex items-center gap-1">
                <span class="icon-[lucide--heart] text-sm"></span>
                <span>{{ artwork.likes }}</span>
              </div>
            </div>
          </div>

          <!-- 卡片下部信息 -->
          <div class="p-3">
            <div
              class="font-semibold text-sm group-hover:text-primary transition-colors line-clamp-2 cursor-pointer"
              @click="navigateToDetail(artwork.id)"
            >
              {{ artwork.title }}
            </div>
            <p class="text-xs text-muted-foreground mt-1">{{ artwork.artist }}</p>

            <!-- 标签 -->
            <div class="flex flex-wrap gap-1 mt-3">
              <span
                v-for="tag in artwork.tags.slice(0, 2)"
                :key="tag"
                class="text-xs bg-muted text-muted-foreground px-2 py-1 rounded hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                #{{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div
        v-if="!artworkStore.loading && artworkStore.artworks.length === 0"
        class="text-center py-20"
      >
        <p class="text-muted-foreground text-lg mb-4">暂无作品</p>
        <button @click="refreshArtworks" class="btn btn-primary">重新加载</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useArtworkStore } from '@/stores/artwork'

const router = useRouter()
const artworkStore = useArtworkStore()

// 默认头像
const defaultAvatar = artworkStore.defaultAvatarUrl

// 刷新随机作品
const refreshArtworks = async () => {
  try {
    await artworkStore.fetchRandomArtworks(10)
  } catch (err) {
    console.error('加载随机作品失败:', err)
  }
}

// 导航到详情页
const navigateToDetail = (id: number) => {
  router.push(`/artwork/${id}`)
}

// 收藏操作
const handleBookmark = async (id: number) => {
  try {
    await artworkStore.toggleBookmark(id)
  } catch (err) {
    console.error('收藏失败:', err)
  }
}

// 初始化加载
refreshArtworks()
</script>
