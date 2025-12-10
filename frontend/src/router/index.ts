import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/pages/HomeView.vue'
import FavoriteView from '@/pages/FavoriteView.vue'
import UploadView from '@/pages/UploadView.vue'
import ArtworkDetail from '@/pages/ArtworkDetail.vue'
import RandomView from '@/pages/RandomView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: HomeView },
    { path: '/favorites', component: FavoriteView },
    { path: '/upload', component: UploadView },
    { path: '/random', component: RandomView },
    {
      path: '/artwork/:id',
      name: 'artwork-detail',
      component: ArtworkDetail, // 或懒加载
      props: true, // 可以直接通过 props 接收 id
    },
  ],
})

export default router
