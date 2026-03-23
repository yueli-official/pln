<template>
  <div class="bg-background text-foreground py-10 px-4">
    <div class="max-w-7xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-8 flex items-end justify-between">
        <div>
          <h1 class="text-3xl font-bold mb-1">浏览作品</h1>
          <p class="text-sm text-muted-foreground">
            共 <span class="text-primary font-medium">{{ artworkStore.total }}</span> 个作品
          </p>
        </div>
        <div v-if="artworkStore.loading && artworkStore.artworks.length > 0" class="flex items-center gap-2 text-sm text-muted-foreground">
          <span class="icon-[lucide--loader-2] size-4 animate-spin"></span>
          加载中
        </div>
      </div>

      <!-- 错误提示 -->
      <div
        v-if="artworkStore.error"
        class="flex items-center gap-3 mb-6 p-4 rounded-xl border border-error/20 bg-error/5 text-error text-sm"
      >
        <span class="icon-[lucide--circle-alert] size-5 shrink-0"></span>
        <span>{{ artworkStore.error.message }}</span>
      </div>

      <!-- 骨架屏 -->
      <div
        v-if="artworkStore.loading && artworkStore.artworks.length === 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4 mb-10"
      >
        <div
          v-for="i in pageSize"
          :key="i"
          class="aspect-square rounded-xl bg-muted animate-pulse"
        ></div>
      </div>

      <!-- 作品网格 -->
      <div
        v-else-if="artworkStore.artworks.length > 0"
        class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4 mb-10"
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
        v-else-if="!artworkStore.loading"
        class="flex flex-col items-center justify-center py-24"
      >
        <span class="icon-[lucide--inbox] size-16 text-muted-foreground/30 mb-4"></span>
        <p class="text-muted-foreground">暂无作品</p>
      </div>

      <!-- 分页器 -->
      <div v-if="artworkStore.total > pageSize" class="flex justify-center pb-4">
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
