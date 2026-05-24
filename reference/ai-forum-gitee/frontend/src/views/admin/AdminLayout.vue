<template>
  <el-container class="admin-layout">
    <el-aside width="220px" class="admin-sidebar">
      <div class="sidebar-header">
        <h2>管理后台</h2>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/admin">
          <el-icon><DataBoard /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/admin/stats">
          <el-icon><TrendCharts /></el-icon>
          <span>数据统计</span>
        </el-menu-item>
        <el-menu-item index="/admin/audit">
          <el-icon><DocumentChecked /></el-icon>
          <span>内容审核</span>
        </el-menu-item>
        <el-menu-item index="/admin/posts">
          <el-icon><Document /></el-icon>
          <span>帖子管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/invites">
          <el-icon><Ticket /></el-icon>
          <span>邀请码管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/boards">
          <el-icon><Grid /></el-icon>
          <span>板块管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/config">
          <el-icon><Setting /></el-icon>
          <span>系统配置</span>
        </el-menu-item>
        <el-menu-item index="/admin/sensitive-words">
          <el-icon><Warning /></el-icon>
          <span>敏感词管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/roles">
          <el-icon><Key /></el-icon>
          <span>角色权限</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="admin-header">
        <div class="header-right">
          <span class="username">{{ adminStore.adminUser?.username || '管理员' }}</span>
          <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useAdminStore } from '@/stores/admin'
import { DataBoard, TrendCharts, DocumentChecked, Document, User, Ticket, Grid, Setting, Warning, Key } from '@element-plus/icons-vue'

const adminStore = useAdminStore()

function handleLogout() {
  adminStore.clearAdminUser()
  localStorage.removeItem('access_token')
  window.location.href = '/admin/login'
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  min-width: 1200px;
}
.admin-sidebar {
  background-color: #304156;
}
.sidebar-header {
  padding: 20px;
  text-align: center;
}
.sidebar-header h2 {
  color: #fff;
  font-size: 18px;
  margin: 0;
}
.admin-header {
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.username {
  color: #333;
}
.admin-main {
  background: #f0f2f5;
  padding: 20px;
}
</style>
