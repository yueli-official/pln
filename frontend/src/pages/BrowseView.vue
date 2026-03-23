<template>
  <div class="bg-background text-foreground py-12 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-10">
        <h1 class="text-4xl font-bold mb-2">浏览作品</h1>
        <p class="text-muted-foreground">共 {{ artworkStore.total }} 个作品</p>
      </div>

      <!-- 错误提示 -->
      <div
        v-if="artworkStore.error"
        class="alert alert-error mb-6 rounded-lg border border-error/20"
      >
        <svg class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4v2m0 0a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <span>{{ artworkStore.error.message }}</span>
      </div>

      <!-- 加载中 -->
      <div
        v-if="artworkStore.loading && artworkStore.artworks.length === 0"
        class="flex justify-center items-center py-20"
      >
        <div class="relative w-16 h-16">
          <div class="absolute inset-0 rounded-full border-4 border-border"></div>
          <div
            class="absolute inset-0 rounded-full border-4 border-transparent border-t-primary animate-spin"
          ></div>
        </div>
      </div>

      <!-- 作品网格 -->
      <div
        v-if="artworkStore.artworks.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4 mb-12"
      >
        <ArtworkCard
          v-for="artwork in artworkStore.artworks"
          :key="artwork.id"
          :artwork="artwork"
          @navigate="navigateToDetail"
          @like="handleLike"
        />
      </div>

      <!-- 空状态 -->
      <div
        v-if="!artworkStore.loading && artworkStore.artworks.length === 0"
        class="flex flex-col items-center justify-center py-20"
      >
        <span class="icon-[lucide--inbox] text-6xl text-muted-foreground/50 mb-4"></span>
        <p class="text-muted-foreground text-lg">暂无作品</p>
      </div>

      <!-- 分页器 -->
      <div v-if="artworkStore.total > 0" class="flex justify-center mb-8">
        <PageNavigator
          v-model:currentPage="currentPage"
          :total="artworkStore.total"
          :page-size="pageSize"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useArtworkStore } from '@/stores/artwork'
import ArtworkCard from '@/components/ArtworkCard.vue'
import PageNavigator from '@/components/PageNavigator.vue'
import { toast } from '@yuelioi/toast'

const router = useRouter()
const route = useRoute()
const artworkStore = useArtworkStore()

const pageSize = 24
const currentPage = ref<number>(parseInt((route.query.page as string) || '1') || 1)

const loadArtworks = async (page: number): Promise<void> => {
  await artworkStore.fetchArtworks(page, pageSize)
}

// 翻页时更新 URL
watch(currentPage, (page) => {
  router.replace({ query: { ...route.query, page: page === 1 ? undefined : String(page) } })
  loadArtworks(page)
  window.scrollTo({ top: 0, behavior: 'smooth' })
})

// 浏览器前进/后退时同步页码
watch(
  () => route.query.page,
  (page) => {
    const p = parseInt((page as string) || '1') || 1
    if (p !== currentPage.value) {
      currentPage.value = p
    }
  },
)

const navigateToDetail = (id: number): void => {
  router.push(`/artwork/${id}`)
}

const handleLike = async (id: number): Promise<void> => {
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

loadArtworks(currentPage.value)
</script>
