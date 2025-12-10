<template>
  <div v-if="totalPages > 1" class="w-max">
    <div
      class="flex items-center border rounded-lg justify-center p-1 h-full bg-card w-full">
      <!-- 上一页 -->
      <button
        @click="handlePrevious"
        :disabled="isFirstPage"
        class="btn btn-sm btn-icon btn-ghost"
        :class="{ 'hover:bg-transparent': isFirstPage }"
        aria-label="上一页">
        <span class="icon-[lucide--chevron-left] size-3.5 sm:size-4"></span>
      </button>

      <!-- 页码按钮 -->
      <div
        class="flex items-center space-x-0.5 sm:space-x-1 min-w-0 overflow-x-auto scrollbar-hide">
        <template v-for="page in visiblePages" :key="page">
          <button
            v-if="typeof page === 'number'"
            @click="handlePageChange(page)"
            class="btn btn-icon-xs btn-ghost"
            :class="
              currentPage === page
                ? 'bg-primary  text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            ">
            {{ page }}
          </button>
          <span v-else class="btn btn-icon-xs btn-ghost">
            {{ page }}
          </span>
        </template>
      </div>

      <!-- 下一页 -->
      <button
        @click="handleNext"
        :disabled="isLastPage"
        class="btn btn-icon-xs btn-ghost"
        :class="{ 'hover:bg-transparent': isLastPage }"
        aria-label="下一页">
        <span class="icon-[lucide--chevron-right] size-3.5 sm:size-4"></span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

export interface PageNavigatorProps {
  total: number;
  pageSize?: number;
  maxVisiblePagesMobile?: number;
  maxVisiblePagesDesktop?: number;
}

const props = withDefaults(defineProps<PageNavigatorProps>(), {
  pageSize: 12,
  maxVisiblePagesMobile: 5,
  maxVisiblePagesDesktop: 8,
});

const currentPage = defineModel<number>("currentPage", { default: 1 });

const totalPages = computed(() => Math.ceil(props.total / props.pageSize));

const isFirstPage = computed(() => currentPage.value === 1);
const isLastPage = computed(() => currentPage.value === totalPages.value);

const isMounted = ref(false);

onMounted(() => {
  isMounted.value = true;
});

const isMobile = computed(() => {
  if (!isMounted.value) return false; // SSR & hydration 期间固定 false
  return window.innerWidth < 640;
});

const visiblePages = computed(() => {
  const pages: (number | string)[] = [];
  const totalPageCount = totalPages.value;
  const current = currentPage.value;

  // 移动端显示更少的页码
  const maxVisiblePages = isMobile.value
    ? props.maxVisiblePagesMobile
    : props.maxVisiblePagesDesktop;

  // 总页数较少，直接显示所有页
  if (totalPageCount <= maxVisiblePages) {
    for (let i = 1; i <= totalPageCount; i++) pages.push(i);
    return pages;
  }

  // 移动端逻辑：只显示 首页 + 当前页附近 + 尾页
  if (isMobile.value) {
    pages.push(1); // 首页

    if (current > 3) {
      pages.push("...");
    }

    // 当前页附近只显示 1 个
    const start = Math.max(2, current - 1);
    const end = Math.min(totalPageCount - 1, current + 1);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (current < totalPageCount - 2) {
      pages.push("...");
    }

    pages.push(totalPageCount); // 尾页
  } else {
    // 桌面端逻辑：显示更多页码
    pages.push(1);

    if (current > 3) {
      pages.push("...");
    }

    const start = Math.max(2, current - 2);
    const end = Math.min(totalPageCount - 1, current + 2);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (current < totalPageCount - 3) {
      pages.push("...");
    }

    pages.push(totalPageCount);
  }

  return pages;
});

// 事件处理函数
const handlePageChange = (page: number) => {
  if (page === currentPage.value) return;
  currentPage.value = page;
};

const handlePrevious = () => {
  if (!isFirstPage.value) {
    currentPage.value--;
  }
};

const handleNext = () => {
  if (!isLastPage.value) {
    currentPage.value++;
  }
};

const resetPage = () => {
  currentPage.value = 1;
};

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

defineExpose({
  resetPage,
  goToPage,
  totalPages,
});
</script>

<style scoped>
/* 隐藏滚动条但保持滚动功能 */
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
</style>
