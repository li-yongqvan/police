import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { setUnauthorizedHandler } from './api/http'
import { useSessionStore } from './stores/session'
import './style.css'
import '../../mobile-web/styles/index.css'

window.addEventListener('vite:preloadError', (event) => {
  event.preventDefault()
  window.location.reload()
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

setUnauthorizedHandler(() => {
  const session = useSessionStore()
  session.logout()
  router.push({ name: 'login' })
})

app.mount('#app')
