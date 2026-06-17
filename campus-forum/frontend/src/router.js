import { createRouter, createWebHistory } from 'vue-router'
import { useSessionStore } from './stores/session'

const DemoLogin = () => import('./views/DemoLogin.vue')
const CommunityLayout = () => import('./views/CommunityLayout.vue')
const CommunityHome = () => import('./views/CommunityHome.vue')
const BoardView = () => import('./views/BoardView.vue')
const PostDetail = () => import('./views/PostDetail.vue')
const NewPost = () => import('./views/NewPost.vue')
const ProfileView = () => import('./views/ProfileView.vue')
const AdminLayout = () => import('./views/AdminLayout.vue')
const AdminOverview = () => import('./views/AdminOverview.vue')
const AdminAudit = () => import('./views/AdminAudit.vue')
const AdminConfig = () => import('./views/AdminConfig.vue')
const AdminUsers = () => import('./views/AdminUsers.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'login', component: DemoLogin },
    {
      path: '/community',
      component: CommunityLayout,
      children: [
        { path: '', name: 'community-home', component: CommunityHome },
        { path: 'boards/:slug', name: 'board', component: BoardView },
        { path: 'posts/new', name: 'new-post', component: NewPost },
        { path: 'posts/:id', name: 'post-detail', component: PostDetail },
        { path: 'profile', name: 'profile', component: ProfileView },
      ],
    },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        { path: '', name: 'admin-overview', component: AdminOverview },
        { path: 'audit', name: 'admin-audit', component: AdminAudit },
        { path: 'config', name: 'admin-config', component: AdminConfig },
        { path: 'users', name: 'admin-users', component: AdminUsers },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const session = useSessionStore()
  if (to.path === '/') {
    return true
  }

  if (!session.token) {
    return { name: 'login' }
  }

  if (to.path.startsWith('/admin') && !session.canAccessAdmin) {
    return { name: 'community-home' }
  }

  return true
})

export default router
