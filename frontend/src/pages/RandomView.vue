<template>
  <div class="bg-background text-foreground py-10 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-8 flex items-end justify-between">
        <div>
          <h1 class="text-3xl font-bold mb-1">随机发现</h1>
          <p class="text-sm text-muted-foreground">每次都有惊喜</p>
        </div>
        <button
          @click="refreshArtworks"
          :disabled="artworkStore.loading"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:opacity-90 active:scale-95 transition-all duration-150 disabled:opacity-50"
        >
          <span
            class="icon-[lucide--shuffle] size-4"
            :class="artworkStore.loading ? 'animate-spin' : ''"
          ></span>
          {{ artworkStore.loading ? '加载中' : '换一批' }}
        </button>
      </div>

      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="flex items-center gap-3 mb-6 p-4 rounded-xl border border-error/20 bg-error/5 text-error text-sm">
        <span class="icon-[lucide--circle-alert] size-5 shrink-0"></span>
        <span>{{ artworkStore.error.message }}</span>
      </div>

      <!-- 骨架屏 -->
      <div
        v-if="artworkStore.loading && artworkStore.artworks.length === 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4"
      >
        <div v-for="i in 24" :key="i" class="aspect-square rounded-xl bg-muted animate-pulse"></div>
      </div>

      <!-- 作品网格 -->
      <div
        v-else-if="artworkStore.artworks.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4"
      >
        <ArtworkCard
          v-for="artwork in artworkStore.artworks"
          :key="artwork.id"
          :artwork="artwork"
          @navigate="navigateToDetail"
          @bookmark="handleBookmark"
        />
      </div>

      <!-- 空状态 -->
      <div
        v-else-if="!artworkStore.loading"
        class="flex flex-col items-center justify-center py-24"
      >
        <span class="icon-[lucide--inbox] size-16 text-muted-foreground/30 mb-4"></span>
        <p class="text-muted-foreground mb-6">暂无作品</p>
        <button
          @click="refreshArtworks"
          class="px-5 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:opacity-90 transition-opacity"
        >
          重新加载
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useArtworkStore } from '@/stores/artwork'
import ArtworkCard from '@/components/ArtworkCard.vue'
import { toast } from '@yuelioi/toast'

const router = useRouter()
const artworkStore = useArtworkStore()

const refreshArtworks = async (): Promise<void> => {
  try {
    await artworkStore.fetchRandomArtworks(24)
  } catch (err) {
    console.error('加载随机作品失败:', err)
    toast.error('加载失败')
  }
}

const navigateToDetail = (id: number): void => {
  router.push(`/artwork/${id}`)
}

const handleBookmark = async (id: number): Promise<void> => {
  try {
    await artworkStore.toggleBookmark(id)
    if (artworkStore.isBookmarked(id)) {
      toast.success('已收藏')
    } else {
      toast.info('已取消收藏')
    }
  } catch (err) {
    console.error('收藏失败:', err)
    toast.error('操作失败')
  }
}

refreshArtworks()
</script>
