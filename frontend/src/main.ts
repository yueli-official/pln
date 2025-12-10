import { createApp } from 'vue'
import { createPinia } from 'pinia'

import '@yuelioi/toast/dist/toast.css'

import App from './App.vue'
import router from './router'

import '@/styles/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
