<template>
  <div class="bg-background text-foreground py-12 px-4 min-h-screen">
    <div class="max-w-6xl mx-auto">
      <!-- 返回按钮 -->
      <RouterLink
        to="/browse"
        class="inline-flex items-center gap-2 px-4 py-2 mb-8 rounded-full bg-card/50 hover:bg-card transition-colors border border-border/50"
      >
        <span class="icon-[lucide--arrow-left]"></span>
        返回
      </RouterLink>

      <!-- 加载状态 -->
      <div v-if="artworkStore.loading" class="flex justify-center items-center py-20">
        <div class="relative w-16 h-16">
          <div class="absolute inset-0 rounded-full border-4 border-border"></div>
          <div
            class="absolute inset-0 rounded-full border-4 border-transparent border-t-primary animate-spin"
          ></div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="alert alert-error mb-6 rounded-lg">
        {{ artworkStore.error.message }}
      </div>

      <!-- 作品详情 -->
      <div v-if="artworkStore.currentArtwork" class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 左侧：大图展示 -->
        <div class="lg:col-span-2">
          <div class="rounded-2xl overflow-hidden bg-muted shadow-2xl">
            <img :src="artworkStore.currentArtwork.url" :alt="'普拉娜'" class="w-full h-auto" />
          </div>
        </div>

        <!-- 右侧：交互面板 -->
        <div class="space-y-5">
          <!-- 统计卡片 -->
          <div class="rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm p-6">
            <h3 class="font-semibold mb-4 text-lg">统计信息</h3>
            <div class="space-y-4">
              <div class="flex justify-between items-center pb-4 border-b border-border/30">
                <span class="text-muted-foreground flex items-center gap-2">
                  <span class="icon-[lucide--eye]"></span>浏览
                </span>
                <span class="font-semibold text-lg">{{ artworkStore.currentArtwork.views }}</span>
              </div>
              <div class="flex justify-between items-center pb-4 border-b border-border/30">
                <span class="text-muted-foreground flex items-center gap-2">
                  <span class="icon-[lucide--heart] text-red-500"></span>点赞
                </span>
                <span class="font-semibold text-lg">{{ artworkStore.currentArtwork.likes }}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-muted-foreground flex items-center gap-2">
                  <span class="icon-[lucide--bookmark]"></span>收藏
                </span>
                <span class="font-semibold text-lg">{{
                  artworkStore.currentArtwork.bookmarks
                }}</span>
              </div>
            </div>
          </div>

          <!-- 标签卡片 -->
          <div class="rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm p-6">
            <h3 class="font-semibold mb-4 text-lg">标签</h3>
            <div
              v-if="artworkStore.currentArtwork.tags?.length ?? 0 > 0"
              class="flex flex-wrap gap-2"
            >
              <span
                v-for="tag in artworkStore.currentArtwork.tags"
                :key="tag"
                class="px-3 py-1 rounded-full bg-primary/10 text-primary text-sm"
              >
                #{{ tag }}
              </span>
            </div>
            <p v-else class="text-muted-foreground text-sm">暂无标签</p>
          </div>

          <!-- 操作按钮 -->
          <div class="space-y-3">
            <button
              @click="handleBookmark"
              :disabled="artworkStore.loading"
              class="w-full px-4 py-3 rounded-lg bg-linear-to-r from-primary to-primary/80 text-white font-medium hover:shadow-lg hover:shadow-primary/30 transition-all disabled:opacity-50"
            >
              <span class="icon-[lucide--heart] mr-2"></span>
              {{ artworkStore.isBookmarked(artworkStore.currentArtwork.id) ? '已收藏' : '收藏' }}
            </button>

            <button
              @click="handleLike"
              :disabled="artworkStore.loading"
              class="w-full px-4 py-3 rounded-lg bg-linear-to-r from-red-500 to-red-600 text-white font-medium hover:shadow-lg hover:shadow-red-500/30 transition-all disabled:opacity-50"
            >
              <span class="icon-[lucide--thumbs-up] mr-2"></span>
              {{ artworkStore.isLiked(artworkStore.currentArtwork.id) ? '已点赞' : '点赞' }}
            </button>

            <button
              @click="handleShare"
              class="w-full px-4 py-3 rounded-lg border border-primary/50 bg-primary/10 text-primary font-medium hover:bg-primary/20 transition-all"
            >
              <span class="icon-[lucide--share-2] mr-2"></span>分享
            </button>
          </div>

          <!-- 管理操作 -->
          <div v-if="hasApiKey" class="rounded-xl border border-warning/50 bg-warning/5 p-6">
            <h3 class="font-semibold mb-4 text-warning text-lg">管理操作</h3>
            <div v-if="!editMode" class="space-y-2">
              <button
                @click="startEdit"
                class="w-full px-4 py-2 rounded-lg bg-warning/10 text-warning hover:bg-warning/20 transition-colors text-sm"
              >
                <span class="icon-[lucide--edit-2] mr-2"></span>编辑
              </button>
              <button
                @click="openDeleteModal"
                class="w-full px-4 py-2 rounded-lg bg-error/10 text-error hover:bg-error/20 transition-colors text-sm"
              >
                <span class="icon-[lucide--trash-2] mr-2"></span>删除
              </button>
            </div>

            <!-- 编辑表单 -->
            <div v-else class="space-y-3">
              <div>
                <label class="label">
                  <span class="label-text text-sm">标签</span>
                </label>
                <input
                  v-model="editData.tagsInput"
                  type="text"
                  placeholder="用逗号分隔"
                  class="input input-bordered w-full input-sm"
                />
              </div>
              <div class="flex gap-2">
                <button
                  @click="saveChanges"
                  :disabled="updating"
                  class="flex-1 px-3 py-2 rounded-lg bg-success/10 text-success hover:bg-success/20 transition-colors text-sm font-medium"
                >
                  <span v-if="!updating" class="icon-[lucide--save] mr-1"></span>
                  保存
                </button>
                <button
                  @click="cancelEdit"
                  class="flex-1 px-3 py-2 rounded-lg bg-card hover:bg-card/80 transition-colors text-sm"
                >
                  取消
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <dialog ref="deleteModalRef" class="modal">
      <div class="modal-box">
        <h3 class="font-bold text-lg">删除确认</h3>
        <p class="py-4">确定要删除这个作品吗？此操作不可恢复。</p>
        <div class="modal-action">
          <form method="dialog">
            <button class="btn">取消</button>
          </form>
          <button @click="handleDelete" :disabled="deleting" class="btn btn-destructive">
            <span v-if="deleting" class="loading loading-spinner loading-sm"></span>
            删除
          </button>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useArtworkStore } from '@/stores/artwork'
import { toast } from '@yuelioi/toast'

const route = useRoute()
const router = useRouter()
const artworkStore = useArtworkStore()

const updating = ref<boolean>(false)
const deleting = ref<boolean>(false)
const editMode = ref<boolean>(false)
const deleteModalRef = ref<HTMLDialogElement | null>(null)

const editData = ref<{
  tagsInput: string
  tags: string[]
}>({
  tagsInput: '',
  tags: [],
})

const hasApiKey = computed<boolean>(() => {
  return localStorage.getItem('api_key') !== null
})

const loadArtwork = async (): Promise<void> => {
  try {
    const id = parseInt(route.params.id as string)
    await artworkStore.fetchArtwork(id)
  } catch (err) {
    console.error('加载作品失败:', err)
  }
}

const handleBookmark = async (): Promise<void> => {
  try {
    if (artworkStore.currentArtwork) {
      await artworkStore.toggleBookmark(artworkStore.currentArtwork.id)
      const msg = artworkStore.isBookmarked(artworkStore.currentArtwork.id)
        ? '已收藏'
        : '已取消收藏'
      toast.success(msg)
    }
  } catch (err) {
    console.error('收藏失败:', err)
    toast.error('操作失败')
  }
}

const handleLike = async (): Promise<void> => {
  try {
    if (artworkStore.currentArtwork) {
      await artworkStore.toggleLike(artworkStore.currentArtwork.id)
      const msg = artworkStore.isLiked(artworkStore.currentArtwork.id) ? '已点赞' : '已取消点赞'
      toast.success(msg)
    }
  } catch (err) {
    console.error('点赞失败:', err)
    toast.error('操作失败')
  }
}

const handleShare = async (): Promise<void> => {
  const url = window.location.href
  const text = '普拉娜'

  if (navigator.share) {
    try {
      await navigator.share({
        text: text,
        url: url,
      })
    } catch (err) {
      if (err instanceof Error && err.name !== 'AbortError') {
        console.error('分享失败:', err.message)
      }
    }
  } else {
    try {
      await navigator.clipboard.writeText(url)
      toast.info('链接已复制到剪贴板')
    } catch (err) {
      console.error('复制失败:', err)
      toast.warning('分享失败，请手动复制链接')
    }
  }
}

const startEdit = (): void => {
  if (artworkStore.currentArtwork) {
    editData.value = {
      tagsInput: artworkStore.currentArtwork.tags.join(','),
      tags: [...artworkStore.currentArtwork.tags],
    }
    editMode.value = true
  }
}

const cancelEdit = (): void => {
  editMode.value = false
  editData.value = {
    tagsInput: '',
    tags: [],
  }
}

const saveChanges = async (): Promise<void> => {
  updating.value = true
  try {
    if (!artworkStore.currentArtwork) return

    const tags = editData.value.tagsInput
      .split(/[,，]+/)
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0)

    await artworkStore.update(artworkStore.currentArtwork.id, {
      tags: tags,
    })

    editMode.value = false
    toast.success('更新成功')
  } catch (err) {
    console.error('更新失败:', err)
    toast.error('更新失败')
  } finally {
    updating.value = false
  }
}

const openDeleteModal = (): void => {
  deleteModalRef.value?.showModal()
}

const handleDelete = async (): Promise<void> => {
  deleting.value = true
  try {
    if (!artworkStore.currentArtwork) return

    await artworkStore.remove(artworkStore.currentArtwork.id)
    toast.info('删除成功')
    router.push('/')
  } catch (err) {
    console.error('删除失败:', err)
    toast.error('删除失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  artworkStore.currentArtwork = null
  loadArtwork()
})
</script>
