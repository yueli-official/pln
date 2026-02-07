<template>
  <div
    class="bg-linear-to-br from-background via-background to-background/95 text-foreground py-16 px-4 relative overflow-hidden"
  >
    <!-- 装饰背景 -->
    <div class="absolute top-0 left-0 w-96 h-96 bg-primary/5 rounded-full blur-3xl -z-10"></div>
    <div class="absolute bottom-0 right-0 w-96 h-96 bg-primary/5 rounded-full blur-3xl -z-10"></div>

    <div class="max-w-7xl mx-auto">
      <!-- 页面标题区域 -->
      <div class="text-center mb-16">
        <h1
          class="text-3xl md:text-5xl font-bold mb-4 bg-linear-to-r from-primary via-primary to-primary/70 bg-clip-text text-transparent"
        >
          普拉娜图像档案
        </h1>
        <p class="text-lg text-muted-foreground max-w-2xl mx-auto">
          收录与普拉娜相关的精选插画与同人创作
        </p>
      </div>

      <!-- 导航卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
        <!-- 浏览插画 -->
        <RouterLink
          to="/browse"
          class="group relative overflow-hidden rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm p-8 hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:shadow-primary/20"
        >
          <div
            class="absolute inset-0 bg-linear-to-br from-primary/0 to-primary/5 opacity-0 group-hover:opacity-100 transition-opacity"
          ></div>
          <div class="relative">
            <div class="text-5xl mb-4">🎨</div>
            <h3 class="text-2xl font-semibold mb-2">浏览插画</h3>
            <p class="text-sm text-muted-foreground">查看所有普拉娜相关插画内容</p>
          </div>
          <div
            class="absolute top-6 right-6 text-primary opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <span class="icon-[lucide--arrow-right] text-2xl"></span>
          </div>
        </RouterLink>

        <!-- 随机邂逅 -->
        <RouterLink
          to="/random"
          class="group relative overflow-hidden rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm p-8 hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:shadow-primary/20"
        >
          <div
            class="absolute inset-0 bg-linear-to-br from-primary/0 to-primary/5 opacity-0 group-hover:opacity-100 transition-opacity"
          ></div>
          <div class="relative">
            <div class="text-5xl mb-4">✨</div>
            <h3 class="text-2xl font-semibold mb-2">随机邂逅</h3>
            <p class="text-sm text-muted-foreground">随机展示一张普拉娜插画</p>
          </div>
          <div
            class="absolute top-6 right-6 text-primary opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <span class="icon-[lucide--arrow-right] text-2xl"></span>
          </div>
        </RouterLink>

        <!-- 我的收藏 -->
        <RouterLink
          to="/bookmarks"
          class="group relative overflow-hidden rounded-xl border border-border/50 bg-card/50 backdrop-blur-sm p-8 hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:shadow-primary/20"
        >
          <div
            class="absolute inset-0 bg-linear-to-br from-primary/0 to-primary/5 opacity-0 group-hover:opacity-100 transition-opacity"
          ></div>
          <div class="relative">
            <div class="text-5xl mb-4">❤️</div>
            <h3 class="text-2xl font-semibold mb-2">我的收藏</h3>
            <p class="text-sm text-muted-foreground">保存你最喜欢的普拉娜插画</p>
          </div>
          <div
            class="absolute top-6 right-6 text-primary opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <span class="icon-[lucide--arrow-right] text-2xl"></span>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>
