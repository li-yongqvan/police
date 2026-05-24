<script setup>
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
</script>

<template>
  <section class="panel profile-card">
    <p class="eyebrow">个人主页</p>
    <h2>{{ session.currentUser?.name }}</h2>
    <p>{{ session.currentUser?.department }}</p>
    <p class="detail-copy">{{ session.currentUser?.bio }}</p>

    <div class="profile-grid">
      <div>
        <strong>{{ session.currentUser?.role }}</strong>
        <span>当前角色</span>
      </div>
      <div>
        <strong>{{ session.currentUser?.status }}</strong>
        <span>账号状态</span>
      </div>
      <div>
        <strong>论坛 MVP</strong>
        <span>展示模式</span>
      </div>
    </div>
  </section>
</template>
