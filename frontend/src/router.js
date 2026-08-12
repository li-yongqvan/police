import { createRouter, createWebHistory } from 'vue-router'
import { useSessionStore } from './stores/session'

const DemoLogin = () => import('./views/DemoLogin.vue')
const OAuthQQ = () => import('./views/OAuthQQ.vue')
const AdminLayout = () => import('./views/AdminLayout.vue')
const AdminOverview = () => import('./views/AdminOverview.vue')
const AdminUsers = () => import('./views/AdminUsers.vue')
const AdminInvites = () => import('./views/AdminInvites.vue')

const DISCOURSE_URL = 'http://122.51.233.225:8080/session/sso'

function redirectToDiscourse() {
  window.location.href = DISCOURSE_URL
}

const dynamicImportErrorPattern =
  /Failed to fetch dynamically imported module|Importing a module script failed|Unable to preload CSS|error loading dynamically imported module/i

function isDynamicImportError(error) {
  const message = error?.message || String(error || '')
  return dynamicImportErrorPattern.test(message)
}

function recoverFromDynamicImportError(to) {
  const target = to?.fullPath || window.location.pathname + window.location.search + window.location.hash
  const reloadKey = 'ai-forum:dynamic-import-reload:' + target
  if (sessionStorage.getItem(reloadKey)) return false
  sessionStorage.setItem(reloadKey, '1')
  window.location.assign(target)
  return true
}

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, top: 16, behavior: 'smooth' }
    return { top: 0 }
  },
  routes: [
    { path: '/', name: 'login', component: DemoLogin },
    { path: '/oauth/qq', name: 'oauth-qq', component: OAuthQQ },
    { path: '/register', redirect: { name: 'login' } },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        { path: '', name: 'admin-overview', component: AdminOverview },
        { path: 'users', name: 'admin-users', component: AdminUsers },
        { path: 'stats', redirect: { name: 'admin-overview' } },
        { path: 'config', redirect: { name: 'admin-overview' } },
        { path: 'invites', name: 'admin-invites', component: AdminInvites },
        { path: 'roles', redirect: { name: 'admin-overview' } },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: () => { redirectToDiscourse(); return false } },
  ],
})

router.onError((error, to) => {
  if (!isDynamicImportError(error)) return
  recoverFromDynamicImportError(to)
})

router.afterEach((to) => {
  sessionStorage.removeItem('ai-forum:dynamic-import-reload:' + to.fullPath)
})

router.beforeEach(async (to) => {
  const session = useSessionStore()
  if (to.name === 'login' && session.token) {
    const valid = await session.ensureValidSession()
    if (!valid) return true
    const role = session.currentUser?.role
    if (['admin', 'platform_admin'].includes(role)) {
      return { path: '/admin' }
    }
    redirectToDiscourse()
    return false
  }
  if (to.name === 'login' || to.name === 'oauth-qq') {
    return true
  }
  if (!session.token) {
    return { name: 'login' }
  }
  const valid = await session.ensureValidSession()
  if (!valid) return { name: 'login' }
  if (to.path.startsWith('/admin') && !session.canAccessAdmin) {
    redirectToDiscourse()
    return false
  }
  const platformOnly = ['admin-invites']
  if (platformOnly.includes(to.name) && session.currentUser?.role !== 'platform_admin') {
    return { name: 'admin-overview' }
  }
  return true
})

export default router

