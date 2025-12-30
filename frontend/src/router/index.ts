import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

// 布局组件

// 页面组件 - 使用懒加载
const Home = () => import('@/pages/HomeView.vue')
const Browse = () => import('@/pages/BrowseView.vue')
const Random = () => import('@/pages/RandomView.vue')
const Bookmarks = () => import('@/pages/BookmarksView.vue')
const ArtworkDetail = () => import('@/pages/ArtworkDetail.vue')
const Upload = () => import('@/pages/UploadView.vue')
const NotFound = () => import('@/pages/NotFound.vue')
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: {
      title: '主页',
      requiresAuth: false,
    },
  },
  {
    path: '/browse',
    name: 'browse',
    component: Browse,
    meta: {
      title: '浏览作品',
      requiresAuth: false,
    },
  },
  {
    path: '/random',
    name: 'random',
    component: Random,
    meta: {
      title: '随机发现',
      requiresAuth: false,
    },
  },
  {
    path: '/bookmarks',
    name: 'bookmarks',
    component: Bookmarks,
    meta: {
      title: '我的收藏',
      requiresAuth: false,
    },
  },
  {
    path: '/artwork/:id',
    name: 'artwork-detail',
    component: ArtworkDetail,
    meta: {
      title: '作品详情',
      requiresAuth: false,
    },
    props: (route) => ({
      id: parseInt(route.params.id as string),
    }),
  },
  {
    path: '/upload',
    name: 'upload',
    component: Upload,
    meta: {
      title: '上传作品',
      requiresAuth: true,
    },
  },

  // 默认重定向到首页
  {
    path: '/',
    redirect: '/home',
  },

  // 404 路由 - 必须放在最后
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: {
      title: '页面未找到',
    },
  },
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior: (to, from, savedPosition) => {
    // 如果有保存的位置，返回到保存位置
    if (savedPosition) {
      return savedPosition
    }
    // 否则滚动到顶部
    return { top: 0 }
  },
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  const title = to.meta.title as string
  if (title) {
    document.title = `${title} - 普拉娜作品库`
  }

  // 认证检查
  const requiresAuth = to.meta.requiresAuth as boolean
  const isAuthenticated = !!localStorage.getItem('api_key')

  if (requiresAuth && !isAuthenticated) {
    // 重定向到首页
    next({ name: 'home' })
  } else {
    next()
  }
})

// 全局后置钩子
router.afterEach((to) => {
  // 分析、日志等
  // console.log(`导航到: ${to.path}`)
})

export default router
