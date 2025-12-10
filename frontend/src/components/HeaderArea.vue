<template>
  <!-- Header -->
  <header class="sticky top-0 z-50 bg-card shadow-sm border-b border-border">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <nav class="flex justify-between items-center h-16">
        <!-- Logo -->
        <RouterLink
          to="/"
          class="text-xl md:text-2xl font-bold text-primary hover:opacity-80 transition-opacity"
        >
          图片库
        </RouterLink>

        <!-- 桌面端导航 -->
        <nav class="hidden md:flex items-center">
          <div class="flex items-center space-x-1 bg-muted rounded-full p-1">
            <RouterLink
              to="/"
              class="px-4 py-2 text-sm font-medium transition-all duration-200 rounded-full cursor-pointer"
              :class="[
                $route.path === '/'
                  ? 'text-primary-foreground bg-primary shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-background/60',
              ]"
            >
              首页
            </RouterLink>

            <RouterLink
              to="/random"
              class="px-4 py-2 text-sm font-medium transition-all duration-200 rounded-full cursor-pointer"
              :class="[
                $route.path.startsWith('/random')
                  ? 'text-primary-foreground bg-primary shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-background/60',
              ]"
            >
              随机
            </RouterLink>
            <RouterLink
              to="/favorites"
              class="px-4 py-2 text-sm font-medium transition-all duration-200 rounded-full cursor-pointer"
              :class="[
                $route.path.startsWith('/favorites')
                  ? 'text-primary-foreground bg-primary shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-background/60',
              ]"
            >
              收藏
            </RouterLink>
          </div>
        </nav>

        <!-- 右侧操作 -->
        <div class="flex items-center space-x-2 md:space-x-4">
          <!-- 上传按钮 -->
          <RouterLink to="/upload" class="btn btn-sm btn-primary gap-2 hidden sm:flex">
            <span class="icon-[lucide--upload] size-4"></span>
            上传
          </RouterLink>

          <!-- 移动端上传按钮 -->
          <RouterLink to="/upload" class="btn btn-sm btn-primary sm:hidden">
            <span class="icon-[lucide--upload] size-4"></span>
          </RouterLink>

          <!-- 主题切换 -->
          <ThemeToggle></ThemeToggle>

          <!-- 设置按钮 -->
          <button
            onclick="settingModal.showModal()"
            class="btn btn-sm btn-ghost"
            aria-label="Settings"
          >
            <span class="icon-[lucide--settings] size-5"></span>
          </button>

          <!-- 移动端菜单按钮 -->
          <button
            @click="toggleMobileMenu"
            class="md:hidden p-2 rounded-md text-foreground hover:bg-accent transition-colors"
            aria-label="Toggle menu"
          >
            <span v-if="!isMobileMenuOpen" class="icon-[lucide--menu] size-6"></span>
            <span v-else class="icon-[lucide--x] size-6"></span>
          </button>
        </div>
      </nav>

      <!-- 移动端菜单 -->
      <div v-if="isMobileMenuOpen" class="md:hidden py-4 space-y-2 border-t border-border">
        <RouterLink
          to="/"
          class="block px-3 py-2 rounded-md text-base font-medium text-foreground hover:bg-accent transition-colors"
          @click="closeMobileMenu"
        >
          首页
        </RouterLink>
        <RouterLink
          to="/favorites"
          class="block px-3 py-2 rounded-md text-base font-medium text-foreground hover:bg-accent transition-colors"
          @click="closeMobileMenu"
        >
          收藏
        </RouterLink>
      </div>
    </div>
  </header>

  <!-- 设置弹窗 -->
  <dialog id="settingModal" class="modal" ref="settingModal">
    <div class="modal-content modal-content-md">
      <div class="modal-header mb-0 pb-1">
        <h3 class="modal-title">管理员设置</h3>
        <span onclick="settingModal.close()" class="modal-close">关闭</span>
      </div>

      <div class="modal-body flex gap-4 flex-col">
        <div class="label-float label-float-required">
          <input type="text" v-model="apiKeyInput" placeholder=" " id="apikey" /><label
            for="apikey"
          >
            API Key</label
          >
        </div>

        <!-- 当前保存的 API Key 状态 -->
        <div v-if="currentApiKey" class="alert alert-info">
          <div class="flex items-center gap-2">
            <span class="icon-[lucide--check-circle] size-5"></span>
            <span>已保存 API Key</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-3 justify-between">
          <button @click="saveApiKey" class="btn btn-primary flex-1 gap-2">
            <span class="icon-[lucide--save] size-4"></span>
            保存
          </button>
          <button
            v-if="currentApiKey"
            @click="clearApiKey"
            class="btn btn-destructive flex-1 gap-2"
          >
            <span class="icon-[lucide--trash-2] size-4"></span>
            清除
          </button>
        </div>
      </div>
    </div>

    <!-- 背景遮罩 -->
    <form method="dialog" class="modal-backdrop">
      <button @click="closeSettings">close</button>
    </form>
  </dialog>
</template>

<script setup lang="ts">
import { ThemeToggle } from '@yuelioi/ui'
import { ref, onMounted } from 'vue'

import { toast } from '@yuelioi/toast'

const settingModal = ref<HTMLDialogElement>()

const isMobileMenuOpen = ref(false)
const apiKeyInput = ref('')
const showApiKey = ref(false)
const currentApiKey = ref('')

// 关闭设置
function closeSettings() {
  settingModal.value?.close()
  apiKeyInput.value = ''
  showApiKey.value = false
}

// 保存 API Key
function saveApiKey() {
  if (apiKeyInput.value.trim()) {
    localStorage.setItem('api_key', apiKeyInput.value.trim())
    currentApiKey.value = apiKeyInput.value.trim()
    toast.success('API Key 已保存')
    closeSettings()
  } else {
    toast.error('请输入有效的 API Key')
  }
}

// 清除 API Key
function clearApiKey() {
  if (confirm('确定要清除保存的 API Key 吗？')) {
    localStorage.removeItem('api_key')
    currentApiKey.value = ''
    apiKeyInput.value = ''
    toast.info('API Key 已清除')
  }
}

// 切换菜单
function toggleMobileMenu() {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

function closeMobileMenu() {
  isMobileMenuOpen.value = false
}

// 初始化
onMounted(() => {
  currentApiKey.value = localStorage.getItem('api_key') || ''
})
</script>
