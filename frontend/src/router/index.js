import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/boards',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
  },
  {
    path: '/boards',
    name: 'BoardList',
    component: () => import('../views/BoardList.vue'),
  },
  {
    path: '/boards/:id',
    name: 'BoardDetail',
    component: () => import('../views/BoardDetail.vue'),
  },
  {
    path: '/posts/:id',
    name: 'PostDetail',
    component: () => import('../views/PostDetail.vue'),
  },
  {
    path: '/post/create',
    name: 'PostCreate',
    component: () => import('../views/PostCreate.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/post/:id/edit',
    name: 'PostEdit',
    component: () => import('../views/PostEdit.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/AdminLogin.vue'),
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAdmin: true },
    children: [
      { path: '', name: 'AdminDashboard', component: () => import('../views/admin/AdminDashboard.vue') },
      { path: 'stats', name: 'AdminStats', component: () => import('../views/admin/AdminStats.vue') },
      { path: 'audit', name: 'AdminAudit', component: () => import('../views/admin/AdminAudit.vue') },
      { path: 'posts', name: 'AdminPostManagement', component: () => import('../views/admin/AdminPostManagement.vue') },
      { path: 'users', name: 'AdminUserManagement', component: () => import('../views/admin/AdminUserManagement.vue') },
      { path: 'invites', name: 'AdminInviteManagement', component: () => import('../views/admin/AdminInviteManagement.vue') },
      { path: 'boards', name: 'AdminBoardManagement', component: () => import('../views/admin/AdminBoardManagement.vue') },
      { path: 'config', name: 'AdminConfig', component: () => import('../views/admin/AdminConfig.vue') },
      { path: 'sensitive-words', name: 'AdminSensitiveWords', component: () => import('../views/admin/AdminSensitiveWords.vue') },
      { path: 'roles', name: 'AdminRoleManagement', component: () => import('../views/admin/AdminRoleManagement.vue') },
    ]
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')

  if (to.meta.requiresAdmin && !token) {
    next({ name: 'AdminLogin', query: { redirect: to.fullPath } })
    return
  }

  // Check admin role for admin routes
  if (to.meta.requiresAdmin && token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      if (payload.role !== 'admin' && payload.role !== 'platform_admin') {
        next({ name: 'AdminLogin' })
        return
      }
    } catch (e) {
      localStorage.removeItem('access_token')
      next({ name: 'AdminLogin' })
      return
    }
  }

  if (to.meta.requiresAuth && !token) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  next()
})

export default router
