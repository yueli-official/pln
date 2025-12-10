<template>
  <div class="bg-background text-foreground py-8 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="alert alert-error mb-6">
        {{ artworkStore.error.message }}
      </div>

      <!-- 加载中 -->
      <div
        v-if="artworkStore.loading && artworkStore.artworks.length === 0"
        class="flex justify-center items-center py-12"
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

                <!-- 点赞按钮 -->
                <button
                  @click.stop="handleLike(artwork.id)"
                  :disabled="artworkStore.loading"
                  class="flex items-center justify-center transition-colors"
                  :class="
                    artworkStore.isLiked(artwork.id)
                      ? 'text-red-500'
                      : 'text-white hover:text-red-500'
                  "
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                  >
                    <g fill="none">
                      <path
                        fill="currentColor"
                        d="m15 10l-.74-.123a.75.75 0 0 0 .74.873zM4 10v-.75a.75.75 0 0 0-.75.75zm16.522 2.392l.735.147zM6 20.75h11.36v-1.5H6zm12.56-11.5H15v1.5h3.56zm-2.82.873l.806-4.835l-1.48-.247l-.806 4.836zm-.92-6.873h-.214v1.5h.213zm-3.335 1.67L8.97 8.693l1.248.832l2.515-3.773zM7.93 9.25H4v1.5h3.93zM3.25 10v8h1.5v-8zm16.807 8.54l1.2-6l-1.47-.295l-1.2 6zM8.97 8.692a1.25 1.25 0 0 1-1.04.557v1.5c.92 0 1.778-.46 2.288-1.225zm7.576-3.405A1.75 1.75 0 0 0 14.82 3.25v1.5a.25.25 0 0 1 .246.291zm2.014 5.462c.79 0 1.38.722 1.226 1.495l1.471.294A2.75 2.75 0 0 0 18.56 9.25zm-1.2 10a2.75 2.75 0 0 0 2.697-2.21l-1.47-.295a1.25 1.25 0 0 1-1.227 1.005zm-2.754-17.5a3.75 3.75 0 0 0-3.12 1.67l1.247.832a2.25 2.25 0 0 1 1.873-1.002zM6 19.25c-.69 0-1.25-.56-1.25-1.25h-1.5A2.75 2.75 0 0 0 6 20.75z"
                      />
                      <path stroke="currentColor" stroke-width="1.5" d="M8 10v10" />
                    </g>
                  </svg>
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
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 36 36"
                  xmlns="http://www.w3.org/2000/svg"
                  class="video-like-icon video-toolbar-item-icon"
                >
                  <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M9.77234 30.8573V11.7471H7.54573C5.50932 11.7471 3.85742 13.3931 3.85742 15.425V27.1794C3.85742 29.2112 5.50932 30.8573 7.54573 30.8573H9.77234ZM11.9902 30.8573V11.7054C14.9897 10.627 16.6942 7.8853 17.1055 3.33591C17.2666 1.55463 18.9633 0.814421 20.5803 1.59505C22.1847 2.36964 23.243 4.32583 23.243 6.93947C23.243 8.50265 23.0478 10.1054 22.6582 11.7471H29.7324C31.7739 11.7471 33.4289 13.402 33.4289 15.4435C33.4289 15.7416 33.3928 16.0386 33.3215 16.328L30.9883 25.7957C30.2558 28.7683 27.5894 30.8573 24.528 30.8573H11.9911H11.9902Z"
                    fill="currentColor"
                  ></path>
                </svg>
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
        class="text-center py-12"
      >
        <p class="text-muted-foreground text-lg">暂无作品</p>
      </div>

      <!-- 加载更多按钮 -->
      <div v-if="artworkStore.hasMoreArtworks" class="flex justify-center mt-12">
        <button
          @click="loadMore"
          :disabled="artworkStore.loading"
          class="btn btn-primary btn-lg rounded-full"
        >
          {{ artworkStore.loading ? '加载中...' : '加载更多作品' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useArtworkStore } from '@/stores/artwork'
import { toast } from '@yuelioi/toast'

const router = useRouter()
const artworkStore = useArtworkStore()

// 默认头像
const defaultAvatar = artworkStore.defaultAvatarUrl

// 初始化加载作品
const loadArtworks = async () => {
  await artworkStore.fetchArtworks(1, artworkStore.pageSize)
}

// 加载更多
const loadMore = async () => {
  await artworkStore.loadMoreArtworks()
}

// 导航到详情页
const navigateToDetail = (id: number) => {
  router.push(`/artwork/${id}`)
}

// 收藏操作
const handleLike = async (id: number) => {
  try {
    await artworkStore.toggleLike(id)
    if (artworkStore.isLiked(id)) {
      toast.success('已点赞')
    } else {
      toast.info('已取消点赞')
    }
  } catch (err) {
    console.error('点赞失败:', err)
    toast.error('操作失败')
  }
}

// 初始化
loadArtworks()
</script>

<style scoped>
.loading {
  border: 4px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
