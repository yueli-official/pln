<template>
  <div class="bg-background text-foreground py-10 px-4">
    <div class="max-w-6xl mx-auto">

      <!-- 返回按钮 -->
      <button
        @click="router.back()"
        class="inline-flex items-center gap-1.5 px-3 py-1.5 mb-8 rounded-lg text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
      >
        <span class="icon-[lucide--arrow-left] size-4"></span>
        返回
      </button>

      <!-- 加载骨架 -->
      <div v-if="artworkStore.loading && !artworkStore.currentArtwork" class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div class="lg:col-span-2 aspect-[4/3] rounded-2xl bg-muted animate-pulse"></div>
        <div class="space-y-4">
          <div class="h-32 rounded-xl bg-muted animate-pulse"></div>
          <div class="h-24 rounded-xl bg-muted animate-pulse"></div>
          <div class="h-28 rounded-xl bg-muted animate-pulse"></div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="flex items-center gap-3 p-4 rounded-xl border border-error/20 bg-error/5 text-error text-sm mb-6">
        <span class="icon-[lucide--circle-alert] size-5 shrink-0"></span>
        {{ artworkStore.error.message }}
      </div>

      <!-- 作品详情 -->
      <div v-if="artworkStore.currentArtwork" class="grid grid-cols-1 lg:grid-cols-3 gap-8">

        <!-- 左侧：图片 -->
        <div class="lg:col-span-2">
          <div class="group relative rounded-2xl overflow-hidden bg-muted shadow-2xl shadow-black/20">
            <img
              :src="artworkStore.currentArtwork.preview_url || artworkStore.currentArtwork.url"
              alt="普拉娜"
              class="w-full h-auto block"
            />
            <!-- 查看原图遮罩 -->
            <a
              :href="artworkStore.currentArtwork.url"
              target="_blank"
              rel="noopener"
              class="absolute inset-0 flex items-center justify-center bg-black/0 group-hover:bg-black/30 transition-colors duration-300 opacity-0 group-hover:opacity-100"
              @click.stop
            >
              <span class="flex items-center gap-2 px-4 py-2 rounded-xl bg-white/15 backdrop-blur-sm text-white text-sm font-medium border border-white/20">
                <span class="icon-[lucide--external-link] size-4"></span>
                查看原图
              </span>
            </a>
          </div>
        </div>

        <!-- 右侧：信息面板 -->
        <div class="space-y-4">

          <!-- 统计数据 -->
          <div class="grid grid-cols-3 gap-3">
            <div class="rounded-xl border border-border/50 bg-card/60 p-4 text-center">
              <span class="icon-[lucide--eye] size-4 text-muted-foreground mx-auto mb-1.5 block"></span>
              <div class="text-xl font-bold">{{ formatCount(artworkStore.currentArtwork.views) }}</div>
              <div class="text-xs text-muted-foreground mt-0.5">浏览</div>
            </div>
            <div class="rounded-xl border border-border/50 bg-card/60 p-4 text-center">
              <span class="icon-[lucide--heart] size-4 text-red-400 mx-auto mb-1.5 block"></span>
              <div class="text-xl font-bold">{{ formatCount(artworkStore.currentArtwork.likes) }}</div>
              <div class="text-xs text-muted-foreground mt-0.5">点赞</div>
            </div>
            <div class="rounded-xl border border-border/50 bg-card/60 p-4 text-center">
              <span class="icon-[lucide--bookmark] size-4 text-primary mx-auto mb-1.5 block"></span>
              <div class="text-xl font-bold">{{ formatCount(artworkStore.currentArtwork.bookmarks) }}</div>
              <div class="text-xs text-muted-foreground mt-0.5">收藏</div>
            </div>
          </div>

          <!-- 标签 -->
          <div class="rounded-xl border border-border/50 bg-card/60 p-5">
            <h3 class="text-sm font-medium text-muted-foreground mb-3">标签</h3>
            <div v-if="artworkStore.currentArtwork.tags?.length" class="flex flex-wrap gap-1.5">
              <span
                v-for="tag in artworkStore.currentArtwork.tags"
                :key="tag"
                class="px-2.5 py-1 rounded-lg bg-primary/10 text-primary text-xs font-medium border border-primary/15"
              >
                #{{ tag }}
              </span>
            </div>
            <p v-else class="text-muted-foreground text-sm">暂无标签</p>
          </div>

          <!-- 操作按钮 -->
          <div class="grid grid-cols-2 gap-3">
            <button
              @click="handleLike"
              :disabled="artworkStore.loading"
              class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 disabled:opacity-50"
              :class="artworkStore.isLiked(artworkStore.currentArtwork.id)
                ? 'bg-red-500/15 text-red-400 border border-red-500/25 hover:bg-red-500/25'
                : 'bg-card/60 border border-border/50 text-foreground hover:bg-accent'"
            >
              <span class="icon-[lucide--heart] size-4"></span>
              {{ artworkStore.isLiked(artworkStore.currentArtwork.id) ? '已点赞' : '点赞' }}
            </button>

            <button
              @click="handleBookmark"
              :disabled="artworkStore.loading"
              class="flex items-center justify-center gap-2 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 disabled:opacity-50"
              :class="artworkStore.isBookmarked(artworkStore.currentArtwork.id)
                ? 'bg-primary/15 text-primary border border-primary/25 hover:bg-primary/25'
                : 'bg-card/60 border border-border/50 text-foreground hover:bg-accent'"
            >
              <span class="icon-[lucide--bookmark] size-4"></span>
              {{ artworkStore.isBookmarked(artworkStore.currentArtwork.id) ? '已收藏' : '收藏' }}
            </button>

            <button
              @click="handleShare"
              class="col-span-2 flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl border border-border/50 bg-card/60 text-sm text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <span class="icon-[lucide--share-2] size-4"></span>
              分享链接
            </button>
          </div>

          <!-- 管理操作 -->
          <div v-if="hasApiKey" class="rounded-xl border border-warning/30 bg-warning/5 p-4">
            <h3 class="text-sm font-medium text-warning mb-3 flex items-center gap-1.5">
              <span class="icon-[lucide--shield] size-4"></span>
              管理操作
            </h3>
            <div v-if="!editMode" class="flex gap-2">
              <button
                @click="startEdit"
                class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg bg-warning/10 text-warning hover:bg-warning/20 transition-colors text-xs font-medium"
              >
                <span class="icon-[lucide--pencil] size-3.5"></span>编辑
              </button>
              <button
                @click="openDeleteModal"
                class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg bg-error/10 text-error hover:bg-error/20 transition-colors text-xs font-medium"
              >
                <span class="icon-[lucide--trash-2] size-3.5"></span>删除
              </button>
            </div>

            <!-- 编辑表单 -->
            <div v-else class="space-y-3">
              <div>
                <label class="text-xs text-muted-foreground mb-1.5 block">标签（逗号分隔）</label>
                <input
                  v-model="editData.tagsInput"
                  type="text"
                  placeholder="tag1, tag2, tag3"
                  class="w-full px-3 py-2 rounded-lg border border-border/60 bg-background text-sm focus:outline-none focus:border-primary/50 transition-colors"
                />
              </div>
              <div class="flex gap-2">
                <button
                  @click="saveChanges"
                  :disabled="updating"
                  class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg bg-success/10 text-success hover:bg-success/20 transition-colors text-xs font-medium disabled:opacity-50"
                >
                  <span v-if="updating" class="icon-[lucide--loader-2] size-3.5 animate-spin"></span>
                  <span v-else class="icon-[lucide--check] size-3.5"></span>
                  保存
                </button>
                <button
                  @click="cancelEdit"
                  class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg bg-muted text-muted-foreground hover:bg-accent transition-colors text-xs"
                >
                  <span class="icon-[lucide--x] size-3.5"></span>取消
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <dialog ref="deleteModalRef" class="modal">
      <div class="modal-content modal-content-sm">
        <div class="modal-header">
          <h3 class="modal-title flex items-center gap-2">
            <span class="icon-[lucide--triangle-alert] size-5 text-error"></span>
            确认删除
          </h3>
        </div>
        <div class="modal-body py-3">
          <p class="text-sm text-muted-foreground">此操作不可恢复，确定要删除这个作品吗？</p>
        </div>
        <div class="modal-footer flex gap-2 justify-end">
          <form method="dialog">
            <button class="px-4 py-2 rounded-lg bg-muted text-muted-foreground hover:bg-accent transition-colors text-sm">取消</button>
          </form>
          <button
            @click="handleDelete"
            :disabled="deleting"
            class="flex items-center gap-1.5 px-4 py-2 rounded-lg bg-error text-error-foreground hover:opacity-90 transition-opacity text-sm font-medium disabled:opacity-50"
          >
            <span v-if="deleting" class="icon-[lucide--loader-2] size-4 animate-spin"></span>
            <span v-else class="icon-[lucide--trash-2] size-4"></span>
            删除
          </button>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop"><button>close</button></form>
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

const editData = ref<{ tagsInput: string; tags: string[] }>({
  tagsInput: '',
  tags: [],
})

const hasApiKey = computed<boolean>(() => localStorage.getItem('api_key') !== null)

function formatCount(n: number): string {
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

const loadArtwork = async (): Promise<void> => {
  try {
    await artworkStore.fetchArtwork(parseInt(route.params.id as string))
  } catch (err) {
    console.error('加载作品失败:', err)
  }
}

const handleBookmark = async (): Promise<void> => {
  try {
    if (!artworkStore.currentArtwork) return
    await artworkStore.toggleBookmark(artworkStore.currentArtwork.id)
    toast.success(artworkStore.isBookmarked(artworkStore.currentArtwork.id) ? '已收藏' : '已取消收藏')
  } catch {
    toast.error('操作失败')
  }
}

const handleLike = async (): Promise<void> => {
  try {
    if (!artworkStore.currentArtwork) return
    await artworkStore.toggleLike(artworkStore.currentArtwork.id)
    toast.success(artworkStore.isLiked(artworkStore.currentArtwork.id) ? '已点赞' : '已取消点赞')
  } catch {
    toast.error('操作失败')
  }
}

const handleShare = async (): Promise<void> => {
  const url = window.location.href
  if (navigator.share) {
    try {
      await navigator.share({ text: '普拉娜', url })
    } catch (err) {
      if (err instanceof Error && err.name !== 'AbortError') console.error(err)
    }
  } else {
    try {
      await navigator.clipboard.writeText(url)
      toast.info('链接已复制到剪贴板')
    } catch {
      toast.warning('请手动复制链接')
    }
  }
}

const startEdit = (): void => {
  if (!artworkStore.currentArtwork) return
  editData.value = {
    tagsInput: artworkStore.currentArtwork.tags.join(', '),
    tags: [...artworkStore.currentArtwork.tags],
  }
  editMode.value = true
}

const cancelEdit = (): void => {
  editMode.value = false
  editData.value = { tagsInput: '', tags: [] }
}

const saveChanges = async (): Promise<void> => {
  if (!artworkStore.currentArtwork) return
  updating.value = true
  try {
    const tags = editData.value.tagsInput
      .split(/[,，]+/)
      .map((t) => t.trim())
      .filter(Boolean)
    await artworkStore.update(artworkStore.currentArtwork.id, { tags })
    editMode.value = false
    toast.success('更新成功')
  } catch {
    toast.error('更新失败')
  } finally {
    updating.value = false
  }
}

const openDeleteModal = (): void => deleteModalRef.value?.showModal()

const handleDelete = async (): Promise<void> => {
  if (!artworkStore.currentArtwork) return
  deleting.value = true
  try {
    await artworkStore.remove(artworkStore.currentArtwork.id)
    toast.info('删除成功')
    router.push('/')
  } catch {
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
