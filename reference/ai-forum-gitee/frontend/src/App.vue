<template>
  <div id="app">
    <header class="app-header">
      <nav class="navbar">
        <router-link to="/" class="brand">AI智联论坛</router-link>
        <div class="nav-links" :class="{ 'nav-open': mobileMenuOpen }">
          <template v-if="authStore.isLoggedIn">
            <router-link to="/boards" @click="mobileMenuOpen = false">板块</router-link>
            <router-link to="/post/create" @click="mobileMenuOpen = false">发帖</router-link>
            <div class="user-menu">
              <router-link to="/profile" class="user-info" @click="mobileMenuOpen = false">
                <UserAvatar :user="authStore.user" size="sm" />
                <span class="username">{{ authStore.user?.nickname || authStore.user?.username }}</span>
                <LevelBadge :level="authStore.user?.level" size="sm" />
              </router-link>
              <button class="logout-btn" @click="handleLogout">退出</button>
            </div>
          </template>
          <template v-else>
            <router-link to="/login" @click="mobileMenuOpen = false">登录</router-link>
            <router-link to="/register" @click="mobileMenuOpen = false">注册</router-link>
          </template>
          <button class="hamburger" @click="mobileMenuOpen = !mobileMenuOpen">&#9776;</button>
        </div>
      </nav>
    </header>
    <main class="app-main">
      <router-view />
    </main>
    <footer class="app-footer">
      <p>&copy; 2026 AI智联论坛 - 人工智能技术交流社区</p>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import UserAvatar from './components/UserAvatar.vue'
import LevelBadge from './components/LevelBadge.vue'

const authStore = useAuthStore()
const router = useRouter()
const mobileMenuOpen = ref(false)

function handleLogout() {
  authStore.logout()
  mobileMenuOpen.value = false
  router.push('/login')
}
</script>

<style>
:root {
  --primary-color: #4a90d9;
  --primary-hover: #357abd;
  --text-color: #333;
  --text-muted: #999;
  --bg-color: #f5f7fa;
  --border-color: #e8e8e8;
  --danger-color: #e74c3c;
  --success-color: #2ecc71;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: var(--bg-color);
  color: var(--text-color);
  line-height: 1.6;
}

.app-header {
  background-color: var(--primary-color);
  color: white;
  padding: 0 1rem;
  position: sticky;
  top: 0;
  z-index: 100;
}

.navbar {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.navbar .brand {
  color: white;
  text-decoration: none;
  font-size: 1.25rem;
  font-weight: bold;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-links a {
  color: white;
  text-decoration: none;
  padding: 0.5rem 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.nav-links a:hover,
.nav-links a.router-link-active {
  background-color: rgba(255, 255, 255, 0.15);
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.username {
  font-size: 0.875rem;
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.logout-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.8rem;
}

.logout-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.hamburger {
  display: none;
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
}

.app-main {
  max-width: 1200px;
  margin: 2rem auto;
  padding: 0 1rem;
  min-height: calc(100vh - 140px);
}

.app-footer {
  text-align: center;
  padding: 1rem;
  background-color: #fff;
  border-top: 1px solid var(--border-color);
  color: var(--text-muted);
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .hamburger { display: block; }
  .nav-links {
    display: none;
    flex-direction: column;
    position: absolute;
    top: 60px;
    left: 0;
    right: 0;
    background: var(--primary-color);
    padding: 1rem;
    gap: 0.5rem;
  }
  .nav-links.nav-open { display: flex; }
  .navbar { position: relative; }
  .user-menu { flex-direction: row; }
  .username { max-width: 80px; }
}
</style>
