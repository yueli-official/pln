<template>
  <div class="min-h-[calc(100vh-100px)] flex items-center justify-center px-4">
    <div class="text-center">
      <!-- 404 图标 -->
      <div class="mb-8">
        <span
          class="text-9xl font-bold bg-linear-to-r from-primary to-primary/50 bg-clip-text text-transparent"
        >
          404
        </span>
      </div>

      <!-- 标题 -->
      <h1 class="text-4xl font-bold mb-4">页面未找到</h1>

      <!-- 描述 -->
      <p class="text-muted-foreground text-lg mb-8 max-w-md mx-auto">
        您访问的页面不存在，可能已被移除或链接错误。
      </p>

      <!-- 返回按钮 -->
      <div class="flex items-center justify-center gap-4">
        <RouterLink
          to="/home"
          class="px-6 py-3 rounded-lg bg-primary text-white font-medium hover:shadow-lg hover:shadow-primary/30 transition-all"
        >
          <span class="icon-[lucide--home] mr-2"></span>
          返回首页
        </RouterLink>
        <button
          @click="goBack"
          class="px-6 py-3 rounded-lg border border-border/50 bg-card/50 font-medium hover:bg-card transition-all"
        >
          <span class="icon-[lucide--arrow-left] mr-2"></span>
          返回上一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = (): void => {
  router.back()
}
</script>
