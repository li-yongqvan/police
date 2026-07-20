import { createApp } from 'vue'
import { createPinia } from 'pinia'
import * as Sentry from '@sentry/vue'
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

const sentryDsn = import.meta.env.VITE_SENTRY_DSN

if (sentryDsn) {
  Sentry.init({
    app,
    dsn: sentryDsn,
    environment: import.meta.env.MODE,
    integrations: [Sentry.browserTracingIntegration({ router })],
    tracesSampleRate: Number(import.meta.env.VITE_SENTRY_TRACES_SAMPLE_RATE ?? 0),
  })
}

setUnauthorizedHandler(() => {
  const session = useSessionStore()
  session.logout()
  router.push({ name: 'login' })
})

app.mount('#app')
