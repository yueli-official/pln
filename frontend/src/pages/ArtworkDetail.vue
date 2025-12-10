<template>
  <div class="bg-background text-foreground py-8 px-4 min-h-screen">
    <div class="max-w-6xl mx-auto">
      <!-- 返回按钮 -->
      <RouterLink to="/" class="btn btn-ghost btn-sm mb-6 gap-2">
        <span class="icon-[lucide--arrow-left] size-4"></span>
        返回首页
      </RouterLink>

      <!-- 加载状态 -->
      <div v-if="artworkStore.loading" class="flex justify-center items-center py-12">
        <div class="loading loading-spinner loading-lg text-primary"></div>
      </div>

      <!-- 错误提示 -->
      <div v-if="artworkStore.error" class="alert alert-error mb-6">
        <div>
          <span>{{ artworkStore.error.message }}</span>
        </div>
        <button @click="artworkStore.error = null" class="btn btn-sm btn-ghost">关闭</button>
      </div>

      <!-- 作品详情 -->
      <div v-if="artworkStore.currentArtwork" class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 左侧：大图展示 -->
        <div class="lg:col-span-2">
          <div class="rounded-lg overflow-hidden bg-muted shadow-lg">
            <img
              :src="artworkStore.currentArtwork.cdn_url"
              :alt="artworkStore.currentArtwork.title"
              class="w-full h-auto"
            />
          </div>

          <!-- 图片信息 -->
          <div class="mt-6 card bg-card border border-border">
            <div class="card-body">
              <h1 class="text-3xl font-bold">{{ artworkStore.currentArtwork.title }}</h1>

              <!-- 作者信息 -->
              <div v-if="artworkStore.currentArtwork.artist" class="flex items-center gap-3 mt-4">
                <img
                  :src="artworkStore.currentArtwork.avatar_url || defaultAvatar"
                  :alt="artworkStore.currentArtwork.artist"
                  class="avatar w-12 h-12 rounded-full"
                  referrerpolicy="no-referrer"
                />
                <div>
                  <p class="font-semibold">{{ artworkStore.currentArtwork.artist }}</p>
                  <p class="text-sm text-muted-foreground">
                    {{ formatDate(artworkStore.currentArtwork.created_at) }}
                  </p>
                </div>
              </div>

              <!-- 描述 -->
              <div v-if="artworkStore.currentArtwork.description" class="mt-6">
                <h3 class="font-semibold mb-2">描述</h3>
                <p class="text-muted-foreground whitespace-pre-wrap">
                  {{ artworkStore.currentArtwork.description }}
                </p>
              </div>

              <!-- 标签 -->
              <div
                v-if="
                  artworkStore.currentArtwork.tags && artworkStore.currentArtwork.tags.length > 0
                "
                class="mt-6"
              >
                <h3 class="font-semibold mb-3">标签</h3>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="tag in artworkStore.currentArtwork.tags"
                    :key="tag"
                    class="badge badge-primary"
                  >
                    #{{ tag }}
                  </span>
                </div>
              </div>

              <!-- 分类 -->
              <div v-if="artworkStore.currentArtwork.category" class="mt-6">
                <p class="text-sm text-muted-foreground">
                  <span class="font-semibold text-foreground">分类：</span>
                  {{ artworkStore.currentArtwork.category }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧：交互面板 -->
        <div>
          <!-- 统计信息 -->
          <div class="card bg-card border border-border mb-6">
            <div class="card-body">
              <h3 class="font-semibold mb-4">统计信息</h3>
              <div class="space-y-3">
                <div class="flex justify-between items-center">
                  <span class="text-muted-foreground flex items-center gap-2">
                    <span class="icon-[lucide--eye] size-4"></span>
                    浏览
                  </span>
                  <span class="font-semibold">{{ artworkStore.currentArtwork.views }}</span>
                </div>
                <div class="divider my-2"></div>
                <div class="flex justify-between items-center">
                  <span class="text-muted-foreground flex items-center gap-2">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="18"
                      height="18"
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
                    点赞
                  </span>
                  <span class="font-semibold">{{ artworkStore.currentArtwork.likes }}</span>
                </div>
                <div class="divider my-2"></div>

                <div class="flex justify-between items-center">
                  <span class="text-muted-foreground flex items-center gap-2">
                    <span class="icon-[lucide--heart]"></span>
                    收藏
                  </span>
                  <span class="font-semibold">{{ artworkStore.currentArtwork.bookmarks }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="space-y-2 mb-6 flex gap-3">
            <!-- 收藏按钮 -->
            <button
              @click="handleBookmark"
              :disabled="artworkStore.loading"
              class="btn btn-primary w-full gap-2"
            >
              <span class="icon-[lucide--heart] size-5 fill-current"></span>
              {{ artworkStore.isBookmarked(artworkStore.currentArtwork.id) ? '已收藏' : '收藏' }}
            </button>

            <!-- 点赞按钮 -->
            <button
              @click="handleLike"
              :disabled="artworkStore.loading"
              class="btn btn-destructive w-full gap-2"
            >
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
              {{ artworkStore.isLiked(artworkStore.currentArtwork.id) ? '已点赞' : '点赞' }}
            </button>

            <!-- 分享按钮 -->
            <button @click="handleShare" class="btn btn-ghost w-full gap-2">
              <span class="icon-[lucide--share-2] size-5"></span>
              分享
            </button>
          </div>

          <!-- 管理员操作（仅当有 API Key 时显示） -->
          <div v-if="hasApiKey" class="card bg-card border border-warning">
            <div class="card-body">
              <h3 class="font-semibold mb-4 text-warning">管理操作</h3>

              <!-- 编辑表单 -->
              <div v-if="editMode" class="space-y-3 mb-4">
                <div>
                  <label class="label">
                    <span class="label-text">标题</span>
                  </label>
                  <input
                    v-model="editData.title"
                    type="text"
                    class="input input-bordered w-full input-sm"
                  />
                </div>
                <div>
                  <label class="label">
                    <span class="label-text">描述</span>
                  </label>
                  <textarea
                    v-model="editData.description"
                    class="textarea textarea-bordered w-full textarea-sm"
                    rows="3"
                  ></textarea>
                </div>
                <div>
                  <label class="label">
                    <span class="label-text">标签</span>
                  </label>
                  <input
                    v-model="editData.tagsInput"
                    type="text"
                    placeholder="用空格分隔"
                    class="input input-bordered w-full input-sm"
                  />
                </div>

                <!-- 编辑操作按钮 -->
                <div class="flex gap-2 mt-4">
                  <button
                    @click="saveChanges"
                    :disabled="updating"
                    class="btn btn-sm btn-success flex-1"
                  >
                    <span v-if="!updating" class="icon-[lucide--save] size-4"></span>
                    <span v-if="updating" class="loading"></span>
                    保存
                  </button>
                  <button @click="cancelEdit" class="btn btn-sm btn-ghost flex-1">取消</button>
                </div>
              </div>

              <!-- 管理按钮 -->
              <div v-if="!editMode" class="space-y-2">
                <button @click="startEdit" class="btn btn-warning btn-sm w-full gap-2">
                  <span class="icon-[lucide--edit-2] size-4"></span>
                  编辑
                </button>
                <button
                  @click="openDeleteModal"
                  :disabled="deleting"
                  class="btn btn-error btn-sm w-full gap-2"
                >
                  <span v-if="!deleting" class="icon-[lucide--trash-2] size-4"></span>
                  <span v-if="deleting" class="loading"></span>
                  删除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认对话框 -->
    <dialog ref="deleteModalRef" class="modal">
      <div class="modal-content modal-content-md">
        <div class="modal-header">
          <h3 class="modal-title">删除确认</h3>
          <span class="modal-close" @click="closeDeleteModal">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="size-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </span>
        </div>
        <div class="modal-body">
          <p>确定要删除这个作品吗？</p>
          <p class="mt-2 text-sm text-muted-foreground">此操作不可恢复，请谨慎选择。</p>
        </div>
        <div class="modal-footer">
          <button @click="closeDeleteModal" class="btn btn-ghost">取消</button>
          <button @click="handleDelete" :disabled="deleting" class="btn btn-destructive">
            <span v-if="deleting" class="loading"></span>
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

// 默认头像
const defaultAvatar = artworkStore.defaultAvatarUrl

// 状态
const updating = ref(false)
const deleting = ref(false)
const editMode = ref(false)
const deleteModalRef = ref<HTMLDialogElement | null>(null)

// 编辑数据
const editData = ref({
  title: '',
  description: '',
  tagsInput: '',
  tags: [] as string[],
})

// 计算属性
const hasApiKey = computed(() => {
  return localStorage.getItem('api_key') !== null
})

// 加载作品
const loadArtwork = async () => {
  try {
    const id = parseInt(route.params.id as string)
    await artworkStore.fetchArtwork(id)
  } catch (err) {
    console.error('加载作品失败:', err)
  }
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// 收藏
const handleBookmark = async () => {
  try {
    if (artworkStore.currentArtwork) {
      await artworkStore.toggleBookmark(artworkStore.currentArtwork.id)
      if (artworkStore.isBookmarked(artworkStore.currentArtwork.id)) {
        toast.success('已收藏')
      } else {
        toast.info('已取消收藏')
      }
    }
  } catch (err) {
    console.error('收藏失败:', err)
    toast.error('操作失败')
  }
}
const handleLike = async () => {
  try {
    if (artworkStore.currentArtwork) {
      await artworkStore.toggleLike(artworkStore.currentArtwork.id)
      if (artworkStore.isLiked(artworkStore.currentArtwork.id)) {
        toast.success('已点赞')
      } else {
        toast.info('已取消点赞')
      }
    }
  } catch (err) {
    console.error('点赞失败:', err)
    toast.error('操作失败')
  }
}

// 分享
const handleShare = async () => {
  const url = window.location.href
  const text = `《${artworkStore.currentArtwork?.title}》- ${artworkStore.currentArtwork?.artist}`

  if (navigator.share) {
    try {
      await navigator.share({
        title: artworkStore.currentArtwork?.title,
        text: text,
        url: url,
      })
    } catch (err) {
      if (err instanceof Error && err.name !== 'AbortError') {
        console.error('分享失败:', err.message)
      }
    }
  } else {
    // 降级方案：复制到剪贴板
    try {
      await navigator.clipboard.writeText(url)
      toast.info('链接已复制到剪贴板')
    } catch (err) {
      console.error('复制失败:', err)
      toast.warning('分享失败，请手动复制链接')
    }
  }
}

// 开始编辑
const startEdit = () => {
  if (artworkStore.currentArtwork) {
    editData.value = {
      title: artworkStore.currentArtwork.title,
      description: artworkStore.currentArtwork.description,
      tagsInput: artworkStore.currentArtwork.tags.join(' '),
      tags: [...artworkStore.currentArtwork.tags],
    }
    editMode.value = true
  }
}

// 取消编辑
const cancelEdit = () => {
  editMode.value = false
  editData.value = {
    title: '',
    description: '',
    tagsInput: '',
    tags: [],
  }
}

// 保存更改
const saveChanges = async () => {
  updating.value = true
  try {
    if (!artworkStore.currentArtwork) return

    const tags = editData.value.tagsInput
      .split(/[\s,]+/)
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0)

    await artworkStore.update(artworkStore.currentArtwork.id, {
      title: editData.value.title,
      description: editData.value.description,
      tags: tags,
    })

    editMode.value = false
    toast.success('更新成功')
  } catch (err) {
    console.error('更新失败:', err)
  } finally {
    updating.value = false
  }
}

// 打开删除对话框
const openDeleteModal = () => {
  deleteModalRef.value?.showModal()
}

// 关闭删除对话框
const closeDeleteModal = () => {
  deleteModalRef.value?.close()
}

// 删除作品
const handleDelete = async () => {
  deleting.value = true
  try {
    if (!artworkStore.currentArtwork) return

    await artworkStore.remove(artworkStore.currentArtwork.id)
    closeDeleteModal()
    toast.info('删除成功')
    router.push('/')
  } catch (err) {
    console.error('删除失败:', err)
  } finally {
    deleting.value = false
  }
}

// 生命周期
onMounted(() => {
  loadArtwork()
})
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
