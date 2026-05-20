<script setup>
import { useRouter } from 'vue-router'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()

const roles = [
  {
    role: 'student',
    title: '学生用户',
    intro: '查看论坛、发帖求助、评论互动、体验社区主流程。',
  },
  {
    role: 'admin',
    title: '协会管理员',
    intro: '查看内容审核、活动运营和用户封禁的最小闭环。',
  },
  {
    role: 'platform_admin',
    title: '中台管理员',
    intro: '查看系统配置、监管规则和数据概览。',
  },
]

async function enter(role) {
  await session.loginAs(role)
  session.setFlash('演示账号已就绪，欢迎进入 AI 智联论坛。', 'success')
  if (role === 'student') {
    router.push('/community')
  } else {
    router.push('/admin')
  }
}
</script>

<template>
  <main class="login-shell">
    <section class="login-hero">
      <p class="eyebrow">AI 智联论坛 MVP</p>
      <h1>一个可以直接拿去展示的学院级 AI 社区原型。</h1>
      <p class="lead">
        这版聚焦社区氛围、问答交流、公告运营和最小中台闭环。你可以从不同角色登录，顺着前台到后台把整个故事讲完整。
      </p>
      <div class="metric-row">
        <div>
          <strong>3</strong>
          <span>核心板块</span>
        </div>
        <div>
          <strong>3</strong>
          <span>Go 服务</span>
        </div>
        <div>
          <strong>1</strong>
          <span>最小中台</span>
        </div>
      </div>
    </section>

    <section class="login-panel panel">
      <div class="section-copy">
        <p class="eyebrow">演示登录</p>
        <h2>选择这次展示要切入的身份</h2>
      </div>
      <div class="role-grid">
        <article v-for="item in roles" :key="item.role" class="role-card">
          <h3>{{ item.title }}</h3>
          <p>{{ item.intro }}</p>
          <button class="primary-button" @click="enter(item.role)">进入演示</button>
        </article>
      </div>
    </section>
  </main>
</template>
