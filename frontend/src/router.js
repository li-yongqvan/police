import { createRouter, createWebHistory } from 'vue-router'
import { useSessionStore } from './stores/session'

const DemoLogin = () => import('./views/DemoLogin.vue')
const Register = () => import('./views/Register.vue')
const CommunityLayout = () => import('./views/CommunityLayout.vue')
const CommunityHome = () => import('./views/CommunityHome.vue')
const BoardView = () => import('./views/BoardView.vue')
const PostDetail = () => import('./views/PostDetail.vue')
const NewPost = () => import('./views/NewPost.vue')
const EditPost = () => import('./views/EditPost.vue')
const ProfileView = () => import('./views/ProfileView.vue')
const UserPublic = () => import('./views/UserPublic.vue')
const MessagesView = () => import('./views/MessagesView.vue')
const AboutView = () => import('./views/AboutView.vue')
const MyLibraryView = () => import('./views/MyLibraryView.vue')
const AdminLayout = () => import('./views/AdminLayout.vue')
const AdminOverview = () => import('./views/AdminOverview.vue')
const AdminAudit = () => import('./views/AdminAudit.vue')
const AdminConfig = () => import('./views/AdminConfig.vue')
const AdminUsers = () => import('./views/AdminUsers.vue')
const AdminPosts = () => import('./views/AdminPosts.vue')
const AdminBoards = () => import('./views/AdminBoards.vue')
const AdminInvites = () => import('./views/AdminInvites.vue')
const AdminSensitiveWords = () => import('./views/AdminSensitiveWords.vue')
const AdminRoles = () => import('./views/AdminRoles.vue')
const AdminStats = () => import('./views/AdminStats.vue')
const AdminReports = () => import('./views/AdminReports.vue')

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/boards', redirect: { name: 'login' } },
    { path: '/boards/:id', redirect: { name: 'community-home' } },
    { path: '/login', redirect: { name: 'login' } },
    { path: '/post/create', redirect: { name: 'new-post' } },
    { path: '/posts/:id', redirect: (to) => ({ name: 'post-detail', params: { id: to.params.id } }) },
    { path: '/profile', redirect: { name: 'profile' } },
    { path: '/admin/login', redirect: { name: 'login' } },
    { path: '/', name: 'login', component: DemoLogin },
    { path: '/register', name: 'register', component: Register },
    {
      path: '/community',
      component: CommunityLayout,
      children: [
        { path: '', name: 'community-home', component: CommunityHome },
        { path: 'boards/:slug', name: 'board', component: BoardView },
        { path: 'posts/new', name: 'new-post', component: NewPost },
        { path: 'posts/:id/edit', name: 'edit-post', component: EditPost },
        { path: 'posts/:id', name: 'post-detail', component: PostDetail },
        { path: 'profile', name: 'profile', component: ProfileView },
        {
          path: 'my/posts',
          name: 'my-posts',
          component: MyLibraryView,
          meta: { libraryMode: 'posts' },
        },
        {
          path: 'my/favorites',
          name: 'my-favorites',
          component: MyLibraryView,
          meta: { libraryMode: 'favorites' },
        },
        {
          path: 'my/history',
          name: 'my-history',
          component: MyLibraryView,
          meta: { libraryMode: 'history' },
        },
        { path: 'users/:id', name: 'user-public', component: UserPublic },
        { path: 'messages', name: 'messages', component: MessagesView },
        { path: 'about', name: 'about', component: AboutView },
      ],
    },
    {
      path: '/admin',
      component: AdminLayout,
      children: [
        { path: '', name: 'admin-overview', component: AdminOverview },
        { path: 'audit', name: 'admin-audit', component: AdminAudit },
        { path: 'reports', name: 'admin-reports', component: AdminReports },
        { path: 'posts', name: 'admin-posts', component: AdminPosts },
        { path: 'users', name: 'admin-users', component: AdminUsers },
        { path: 'boards', name: 'admin-boards', component: AdminBoards },
        { path: 'config', name: 'admin-config', component: AdminConfig },
        { path: 'invites', name: 'admin-invites', component: AdminInvites },
        { path: 'sensitive', name: 'admin-sensitive', component: AdminSensitiveWords },
        { path: 'roles', name: 'admin-roles', component: AdminRoles },
        { path: 'stats', name: 'admin-stats', component: AdminStats },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const session = useSessionStore()
  if ((to.name === 'login' || to.name === 'register') && session.token) {
    const role = session.currentUser?.role
    if (['admin', 'platform_admin'].includes(role)) {
      return { path: '/admin' }
    }
    return { path: '/community' }
  }
  if (to.name === 'login' || to.name === 'register') {
    return true
  }

  if (!session.token) {
    return { name: 'login' }
  }

  if (to.path.startsWith('/admin') && !session.canAccessAdmin) {
    return { name: 'community-home' }
  }

  const platformOnly = ['admin-invites', 'admin-sensitive', 'admin-roles']
  if (platformOnly.includes(to.name) && session.currentUser?.role !== 'platform_admin') {
    return { name: 'admin-overview' }
  }

  return true
})

export default router
