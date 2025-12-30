<template>
  <div
    class="group relative overflow-hidden rounded-lg bg-card shadow-sm hover:shadow-md transition-all duration-300"
  >
    <div
      class="relative w-full overflow-hidden bg-muted aspect-square block cursor-pointer"
      @click="emit('navigate', artwork.id)"
    >
      <img
        :src="artwork.thumbnail_url || artwork.url"
        :alt="'普拉娜'"
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
      />

      <!-- 悬停遮罩层 -->
      <div
        class="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-colors duration-300 opacity-0 group-hover:opacity-100 flex flex-col justify-between p-3"
      >
        <!-- 操作按钮 -->
        <div class="flex items-center justify-between">
          <button
            v-if="showLike"
            @click.stop="emit('like', artwork.id)"
            :disabled="disabled"
            class="text-white hover:text-red-500 transition-colors"
            title="点赞"
          >
            <span class="icon-[lucide--heart] text-2xl fill-current"></span>
          </button>

          <button
            v-if="showBookmark"
            @click.stop="emit('bookmark', artwork.id)"
            :disabled="disabled"
            class="text-white hover:text-primary transition-colors"
            title="收藏"
          >
            <span class="icon-[lucide--bookmark] text-2xl fill-current"></span>
          </button>
        </div>
      </div>

      <!-- 右上角统计 -->
      <div
        class="absolute top-2 right-2 flex flex-col gap-1 opacity-0 group-hover:opacity-100 transition-opacity text-white text-xs"
      >
        <div class="bg-black/60 backdrop-blur-sm px-2 py-1 rounded-md flex items-center gap-1">
          <span class="icon-[lucide--eye] text-sm"></span>
          <span>{{ artwork.views }}</span>
        </div>
        <div class="bg-black/60 backdrop-blur-sm px-2 py-1 rounded-md flex items-center gap-1">
          <span class="icon-[lucide--heart] text-sm text-red-500"></span>
          <span>{{ artwork.likes }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Artwork } from '@/types'

interface Props {
  artwork: Artwork
  disabled?: boolean
  showLike?: boolean
  showBookmark?: boolean
}

withDefaults(defineProps<Props>(), {
  disabled: false,
  showLike: true,
  showBookmark: true,
})

const emit = defineEmits<{
  navigate: [id: number]
  like: [id: number]
  bookmark: [id: number]
}>()
</script>

<style scoped></style>
