<template>
  <div class="bg-background text-foreground min-h-[calc(100vh-4rem)] px-4 py-16 relative overflow-hidden">
    <!-- 背景光晕 -->
    <div class="pointer-events-none absolute -top-32 -left-32 w-[36rem] h-[36rem] bg-primary/10 rounded-full blur-[120px]"></div>
    <div class="pointer-events-none absolute -bottom-32 -right-32 w-[36rem] h-[36rem] bg-primary/8 rounded-full blur-[120px]"></div>

    <div class="max-w-5xl mx-auto relative">
      <!-- 英雄区 -->
      <div class="text-center mb-20">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-medium mb-6">
          <span class="icon-[lucide--sparkles] size-3"></span>
          碧蓝档案 · 普拉娜图像档案
        </div>
        <h1 class="text-4xl md:text-6xl font-bold mb-5 leading-tight">
          <span class="bg-linear-to-r from-primary via-primary/90 to-primary/60 bg-clip-text text-transparent">
            普拉娜
          </span>
          <br />
          <span class="text-foreground/90 text-3xl md:text-5xl font-semibold">图像档案库</span>
        </h1>
        <p class="text-muted-foreground text-base md:text-lg max-w-xl mx-auto leading-relaxed">
          收录与普拉娜相关的精选插画、官方立绘与同人创作
        </p>
      </div>

      <!-- 导航卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <RouterLink
          to="/browse"
          class="group relative overflow-hidden rounded-2xl border border-border/60 bg-card/60 backdrop-blur-sm p-7 hover:border-primary/40 hover:bg-card/80 transition-all duration-300 hover:shadow-xl hover:shadow-primary/10 hover:-translate-y-0.5"
        >
          <div class="absolute inset-0 bg-linear-to-br from-primary/0 via-primary/0 to-primary/8 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl"></div>
          <div class="relative">
            <div class="w-12 h-12 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center mb-5 group-hover:bg-primary/15 transition-colors">
              <span class="icon-[lucide--images] size-6 text-primary"></span>
            </div>
            <h3 class="text-lg font-semibold mb-1.5">浏览插画</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">查看全部普拉娜相关插画内容</p>
          </div>
          <span class="absolute bottom-7 right-7 icon-[lucide--arrow-right] size-4 text-primary opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all duration-200"></span>
        </RouterLink>

        <RouterLink
          to="/random"
          class="group relative overflow-hidden rounded-2xl border border-border/60 bg-card/60 backdrop-blur-sm p-7 hover:border-primary/40 hover:bg-card/80 transition-all duration-300 hover:shadow-xl hover:shadow-primary/10 hover:-translate-y-0.5"
        >
          <div class="absolute inset-0 bg-linear-to-br from-primary/0 via-primary/0 to-primary/8 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl"></div>
          <div class="relative">
            <div class="w-12 h-12 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center mb-5 group-hover:bg-primary/15 transition-colors">
              <span class="icon-[lucide--shuffle] size-6 text-primary"></span>
            </div>
            <h3 class="text-lg font-semibold mb-1.5">随机邂逅</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">随机展示一批普拉娜插画</p>
          </div>
          <span class="absolute bottom-7 right-7 icon-[lucide--arrow-right] size-4 text-primary opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all duration-200"></span>
        </RouterLink>

        <RouterLink
          to="/bookmarks"
          class="group relative overflow-hidden rounded-2xl border border-border/60 bg-card/60 backdrop-blur-sm p-7 hover:border-primary/40 hover:bg-card/80 transition-all duration-300 hover:shadow-xl hover:shadow-primary/10 hover:-translate-y-0.5"
        >
          <div class="absolute inset-0 bg-linear-to-br from-primary/0 via-primary/0 to-primary/8 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-2xl"></div>
          <div class="relative">
            <div class="w-12 h-12 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center mb-5 group-hover:bg-primary/15 transition-colors">
              <span class="icon-[lucide--bookmark] size-6 text-primary"></span>
            </div>
            <h3 class="text-lg font-semibold mb-1.5">我的收藏</h3>
            <p class="text-sm text-muted-foreground leading-relaxed">保存你最喜欢的普拉娜插画</p>
          </div>
          <span class="absolute bottom-7 right-7 icon-[lucide--arrow-right] size-4 text-primary opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all duration-200"></span>
        </RouterLink>
      </div>
    </div>
  </div>
</template>
