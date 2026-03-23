<template>
  <div
    class="group relative overflow-hidden rounded-xl bg-muted cursor-pointer aspect-square"
    @click="emit('navigate', artwork.id)"
  >
    <!-- 图片 -->
    <img
      :src="artwork.thumbnail_url || artwork.url"
      :alt="'普拉娜 碧蓝档案 壁纸'"
      class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
      loading="lazy"
    />

    <!-- 底部渐变遮罩 + 操作区 -->
    <div
      class="absolute inset-0 bg-linear-to-t from-black/70 via-black/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex flex-col justify-end p-3"
    >
      <div class="flex items-center justify-between">
        <!-- 统计 -->
        <div class="flex items-center gap-2 text-white/80 text-xs">
          <span class="flex items-center gap-1">
            <span class="icon-[lucide--eye] size-3"></span>
            {{ formatCount(artwork.views) }}
          </span>
          <span class="flex items-center gap-1">
            <span class="icon-[lucide--heart] size-3 text-red-400"></span>
            {{ formatCount(artwork.likes) }}
          </span>
        </div>

        <!-- 操作按钮 -->
        <div class="flex items-center gap-1.5">
          <button
            v-if="showLike"
            @click.stop="emit('like', artwork.id)"
            :disabled="disabled"
            class="p-1.5 rounded-lg bg-white/10 backdrop-blur-sm hover:bg-white/25 transition-colors"
            title="点赞"
          >
            <span
              class="icon-[lucide--heart] size-3.5 transition-colors"
              :class="isLiked ? 'text-red-400' : 'text-white'"
            ></span>
          </button>

          <button
            v-if="showBookmark"
            @click.stop="emit('bookmark', artwork.id)"
            :disabled="disabled"
            class="p-1.5 rounded-lg bg-white/10 backdrop-blur-sm hover:bg-white/25 transition-colors"
            title="收藏"
          >
            <span
              class="icon-[lucide--bookmark] size-3.5 transition-colors"
              :class="isBookmarked ? 'text-primary' : 'text-white'"
            ></span>
          </button>
        </div>
      </div>
    </div>

    <!-- 已点赞角标（常显） -->
    <div
      v-if="isLiked"
      class="absolute top-2 left-2 w-2 h-2 rounded-full bg-red-400 shadow-sm opacity-80"
    ></div>
    <!-- 已收藏角标（常显） -->
    <div
      v-if="isBookmarked"
      class="absolute top-2 right-2 w-2 h-2 rounded-full bg-primary shadow-sm opacity-80"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Artwork } from '@/types'
import { useArtworkStore } from '@/stores/artwork'

interface Props {
  artwork: Artwork
  disabled?: boolean
  showLike?: boolean
  showBookmark?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  showLike: true,
  showBookmark: true,
})

const emit = defineEmits<{
  navigate: [id: number]
  like: [id: number]
  bookmark: [id: number]
}>()

const artworkStore = useArtworkStore()
const isLiked = computed(() => artworkStore.isLiked(props.artwork.id))
const isBookmarked = computed(() => artworkStore.isBookmarked(props.artwork.id))

function formatCount(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}
</script>
