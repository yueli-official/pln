<template>
  <div class="bg-background text-foreground py-8 px-4 min-h-screen">
    <div class="max-w-4xl mx-auto">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold mb-2">上传作品</h1>
        <p class="text-muted-foreground text-sm">支持批量上传多个作品</p>
        <p class="text-muted-foreground text-sm">最大5M, 超过5M会被压缩</p>
      </div>

      <!-- 错误提示 -->
      <div v-if="error" class="alert alert-error mb-6">
        <div>
          <span>{{ error }}</span>
        </div>
        <button @click="error = null" class="btn btn-sm btn-ghost">关闭</button>
      </div>

      <!-- 已成功上传的文件展示 -->
      <div v-if="uploadedFiles.length > 0" class="alert alert-success mb-8">
        <span class="icon-[lucide--check-circle] size-5"></span>
        <div class="flex-1">
          <span>已成功上传 {{ uploadedFiles.length }} 个作品</span>
        </div>
        <button @click="clearUploadedFiles" class="btn btn-sm btn-ghost">清空</button>
      </div>

      <!-- 上传区域 -->
      <div
        @drop.prevent="handleDrop"
        @dragover.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        class="border-2 border-dashed rounded-lg p-8 text-center transition-colors mb-8"
        :class="[
          isDragging
            ? 'border-primary bg-primary/10'
            : 'border-border bg-muted/30 hover:border-primary/50',
        ]"
      >
        <span class="icon-[lucide--cloud-upload] size-12 mx-auto mb-4 text-muted-foreground"></span>
        <h3 class="text-lg font-semibold mb-2">拖拽,粘贴或点击上传图片</h3>
        <p class="text-sm text-muted-foreground mb-4">支持 JPG, PNG, GIF, WebP</p>
        <input type="file" multiple accept="image/*" class="hidden" @change="handleFileSelect" />
        <button @click="triggerFileInput" class="btn btn-primary">选择文件</button>
      </div>

      <!-- 待上传文件列表 -->
      <div v-if="files.length > 0" class="mb-8">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-xl font-semibold">待上传 {{ files.length }} 个文件</h3>
          <button @click="clearFailedFiles" class="btn btn-sm btn-destructive">
            <span class="icon-[lucide--trash-2] size-4"></span>
            清空
          </button>
        </div>

        <!-- 文件列表 -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div
            v-for="(file, index) in files"
            :key="index"
            class="border border-border rounded-lg p-4 flex gap-4"
          >
            <!-- 预览图 -->
            <div class="w-20 h-20 shrink-0 bg-muted rounded-lg overflow-hidden">
              <img
                v-if="previews[index]"
                :src="previews[index]"
                :alt="file.name"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center">
                <span class="icon-[lucide--image] size-8 text-muted-foreground"></span>
              </div>
            </div>

            <!-- 文件信息 -->
            <div class="flex-1 min-w-0">
              <p class="font-medium truncate">{{ file.name }}</p>
              <p class="text-sm text-muted-foreground">{{ formatFileSize(file.size) }}</p>

              <!-- 上传状态 -->
              <div v-if="uploadStatus[index]" class="mt-2">
                <div v-if="uploadStatus[index].type === 'progress'" class="space-y-1">
                  <div class="progress progress-primary h-2">
                    <div
                      class="progress-bar h-full transition-all"
                      :style="{ width: uploadStatus[index].value + '%' }"
                    ></div>
                  </div>
                  <p class="text-xs text-muted-foreground">{{ uploadStatus[index].value }}%</p>
                </div>
                <div
                  v-else-if="uploadStatus[index].type === 'error'"
                  class="text-xs text-destructive"
                >
                  {{ uploadStatus[index].message }}
                </div>
              </div>
            </div>

            <!-- 删除按钮 -->
            <button
              @click="removeFile(index)"
              :disabled="uploading"
              class="btn btn-sm btn-ghost text-destructive hover:bg-destructive/10"
            >
              <span class="icon-[lucide--x] size-4"></span>
            </button>
          </div>
        </div>

        <!-- 表单字段 -->
        <div class="card bg-card border border-border p-6 mb-6">
          <h4 class="text-lg font-semibold mb-4">作品信息</h4>

          <div class="space-y-4">
            <!-- 标题 -->
            <div>
              <label class="label">
                <span class="label-text">作品标题</span>
              </label>
              <input
                v-model="formData.title"
                type="text"
                placeholder="默认: 普拉娜"
                class="input input-bordered w-full"
              />
            </div>

            <!-- 艺术家名称 -->
            <div>
              <label class="label">
                <span class="label-text">艺术家名称</span>
              </label>
              <input
                v-model="formData.artist"
                type="text"
                placeholder="默认: 未知艺术家"
                class="input input-bordered w-full"
              />
            </div>

            <!-- 描述 -->
            <div>
              <label class="label">
                <span class="label-text">描述</span>
              </label>
              <textarea
                v-model="formData.description"
                placeholder="输入作品描述"
                class="textarea textarea-bordered w-full h-24"
              ></textarea>
            </div>

            <!-- 分类 -->
            <div>
              <label class="label">
                <span class="label-text">分类</span>
              </label>
              <select v-model="formData.category" class="select select-bordered w-full">
                <option value="普拉娜">普拉娜</option>
              </select>
            </div>

            <!-- 标签 -->
            <div>
              <label class="label">
                <span class="label-text">标签（用逗号分隔）</span>
              </label>
              <input
                v-model="formData.tagsInput"
                @keydown.enter.prevent="handleAddTag"
                @blur="handleAddTag"
                type="text"
                placeholder="例如: 白毛,单人,"
                class="input input-bordered w-full"
              />
              <div v-if="formData.tags.length > 0" class="flex flex-wrap gap-2 mt-3">
                <span
                  v-for="(tag, index) in formData.tags"
                  :key="index"
                  class="badge badge-primary gap-2"
                >
                  {{ tag }}
                  <button @click="formData.tags.splice(index, 1)" class="hover:opacity-75">
                    <span class="icon-[lucide--x] size-3"></span>
                  </button>
                </span>
              </div>
            </div>

            <!-- 头像 URL -->
            <div>
              <label class="label">
                <span class="label-text">艺术家头像 URL</span>
              </label>
              <input
                v-model="formData.avatar_url"
                type="url"
                placeholder="https://example.com/avatar.jpg"
                class="input input-bordered w-full"
              />
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-3">
          <button
            @click="uploadFiles"
            :disabled="uploading || !canUpload"
            class="btn btn-primary flex-1 gap-2"
          >
            <span v-if="!uploading" class="icon-[lucide--upload] size-4"></span>
            <span v-if="uploading" class="loading loading-spinner loading-sm"></span>
            {{ uploading ? '上传中...' : '开始上传' }}
          </button>

          <RouterLink to="/" class="btn btn-ghost flex-1">
            <span class="icon-[lucide--arrow-left] size-4"></span>
            返回首页
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useApi } from '@/composables/useApi'

const { uploadAndCreateArtwork } = useApi()

// 状态
const files = ref<File[]>([])
const previews = ref<string[]>([])
const uploadedFiles = ref<string[]>([]) // 已成功上传的文件名列表
const isDragging = ref(false)
const uploading = ref(false)
const error = ref<string | null>(null)
const uploadStatus = ref<
  Record<number, { type: 'progress' | 'error'; value?: number; message?: string }>
>({})

// 表单数据
const formData = ref({
  title: '普拉娜',
  artist: '未知艺术家',
  description: '',
  category: '普拉娜',
  avatar_url: '',
  tagsInput: '',
  tags: [] as string[],
})

// 计算属性
const canUpload = computed(() => files.value.length > 0)

// 处理文件选择
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    addFiles(Array.from(target.files))
  }
  target.value = ''
}

// 触发文件选择
const triggerFileInput = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.multiple = true
  input.accept = 'image/*'
  input.onchange = handleFileSelect
  input.click()
}

// 处理拖拽
const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const droppedFiles = event.dataTransfer?.files
  if (droppedFiles) {
    addFiles(Array.from(droppedFiles))
  }
}

// 处理粘贴
const handlePaste = (event: ClipboardEvent) => {
  const items = event.clipboardData?.items
  if (!items) return

  const pastedFiles: File[] = []
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (!item) continue
    if (item.kind === 'file' && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        pastedFiles.push(file)
      }
    }
  }

  if (pastedFiles.length > 0) {
    addFiles(pastedFiles)
  } else {
    error.value = '粘贴板中没有图片'
  }
}

// 添加文件
const addFiles = (newFiles: File[]) => {
  const validFiles = newFiles.filter((file) => file.type.startsWith('image/'))

  if (newFiles.length !== validFiles.length) {
    error.value = '只支持图片文件'
  }

  const invalidSize = validFiles.some((file) => file.size > 100 * 1024 * 1024)
  if (invalidSize) {
    error.value = '文件大小不能超过 100MB'
    return
  }

  validFiles.forEach((file) => {
    if (files.value.some((f) => f.name === file.name && f.size === file.size)) {
      return
    }

    files.value.push(file)

    const reader = new FileReader()
    reader.onload = (e) => {
      previews.value.push(e.target?.result as string)
    }
    reader.readAsDataURL(file)
  })
}

// 移除文件
const removeFile = (index: number) => {
  files.value.splice(index, 1)
  previews.value.splice(index, 1)

  const newStatus: Record<number, any> = {}
  Object.entries(uploadStatus.value).forEach(([key, value]) => {
    const oldIndex = parseInt(key)
    if (oldIndex > index) {
      newStatus[oldIndex - 1] = value
    } else if (oldIndex < index) {
      newStatus[oldIndex] = value
    }
  })
  uploadStatus.value = newStatus
}

// 清空失败文件
const clearFailedFiles = () => {
  files.value = []
  previews.value = []
  uploadStatus.value = {}
}

// 清空已上传文件列表
const clearUploadedFiles = () => {
  uploadedFiles.value = []
}

// 处理添加标签
const handleAddTag = () => {
  if (!formData.value.tagsInput.trim()) return

  const newTags = formData.value.tagsInput
    .split(/[,，\s]+/)
    .map((tag) => tag.trim())
    .filter((tag) => tag && !formData.value.tags.includes(tag))

  formData.value.tags.push(...newTags)
  formData.value.tagsInput = ''
}

// 上传文件
const uploadFiles = async () => {
  uploading.value = true
  error.value = null

  handleAddTag()

  for (let i = 0; i < files.value.length; i++) {
    try {
      const file = files.value[i]
      if (!file) continue

      await uploadAndCreateArtwork(
        file,
        {
          title: formData.value.title,
          artist: formData.value.artist,
          description: formData.value.description,
          category: formData.value.category,
          avatar_url: formData.value.avatar_url,
          tags: formData.value.tags,
        },
        (progress) => {
          uploadStatus.value[i] = { type: 'progress', value: progress }
        },
      )

      // 上传成功：添加到已上传列表，从待上传列表中移除
      uploadedFiles.value.push(file.name)
      files.value.splice(i, 1)
      previews.value.splice(i, 1)

      // 调整索引
      const newStatus: Record<number, any> = {}
      Object.entries(uploadStatus.value).forEach(([key, value]) => {
        const oldIndex = parseInt(key)
        if (oldIndex > i) {
          newStatus[oldIndex - 1] = value
        } else if (oldIndex < i) {
          newStatus[oldIndex] = value
        }
      })
      uploadStatus.value = newStatus
      i-- // 因为移除了当前项，索引需要回退
    } catch (err) {
      error.value = `上传失败: ${(err as Error).message}`
      // 记录该文件的错误状态，但继续处理下一个文件
      uploadStatus.value[i] = {
        type: 'error',
        message: (err as Error).message,
      }
    }
  }

  uploading.value = false
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

onMounted(() => {
  window.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  window.removeEventListener('paste', handlePaste)
})
</script>

<style scoped>
.progress-bar {
  background-color: hsl(var(--p) / 1);
}
</style>
