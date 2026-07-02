<script setup>
import { RouterLink } from 'vue-router'
import GxIcon from '../components/gx/GxIcon.vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import { GX_NAV_ITEMS } from '../composables/useGxNav'

const breadcrumbItems = [
  { label: '首页', to: '/community' },
  { label: '关于' },
]

const boardIntro = {
  study: '课程学习、备考与学术讨论',
  training: '技能训练、实习心得与实务交流',
  notice: '学院与协会通知、活动信息（请以官方通知为准）',
  club: '社团活动、文化建设与风采展示',
}

const registerSteps = [
  { title: '领取邀请码', desc: '向辅导员或协会管理员领取，一人一码，勿在公开群转发。' },
  { title: '打开注册页', desc: '在登录页选择「注册」，填写学号/用户名与密码。' },
  { title: '填写邀请码', desc: '提交邀请码完成账号创建并登录。' },
  { title: '完善资料', desc: '首次登录按提示填写院系、区队、年级等信息。' },
]

const postTips = [
  { icon: 'edit', label: '文明发帖', desc: '遵守校纪校规，理性交流，避免人身攻击与广告引流。' },
  { icon: 'shield', label: '审核说明', desc: '部分内容将进入人工审核，通过后在对应板块展示。' },
  { icon: 'book', label: '附件支持', desc: '可在发帖或详情页上传/下载附件（等级不足时部分功能受限）。' },
]

const reportTips = [
  '在帖子详情页点击「举报」，简要说明理由即可。',
  '管理员将在工作日内处理，严重违规可能删帖并限制账号。',
  '遇紧急或人身安全相关问题，请同时联系辅导员或学院值班老师。',
]
</script>

<template>
  <div class="gx-page gx-about-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <div class="gx-about-layout">
      <section class="gx-about-hero" aria-labelledby="about-hero-title">
        <div class="gx-about-hero__badge" aria-hidden="true">
          <GxIcon name="shield" :size="28" />
        </div>
        <div class="gx-about-hero__content">
          <p class="gx-about-hero__eyebrow">关于本站</p>
          <h1 id="about-hero-title" class="gx-about-hero__title">AI智联平台 · 校内交流论坛</h1>
          <p class="gx-about-hero__lede">
            面向广西警察学院师生的校内交流平台，聚焦学业研讨、警务实训、校园公告与社团风采。
          </p>
          <p class="gx-about-hero__legal">
            视觉风格借鉴警校文化；不使用国家警徽等法定标识。图标均为原创简化设计，不代表公安机关官方形象。
          </p>
          <div class="gx-about-hero__actions">
            <RouterLink to="/register" class="gx-about-hero__cta gx-about-hero__cta--primary">
              注册账号
            </RouterLink>
            <RouterLink to="/community" class="gx-about-hero__cta gx-about-hero__cta--ghost">
              进入社区
            </RouterLink>
          </div>
        </div>
      </section>

      <div class="gx-about-grid">
        <article class="gx-about-topic">
          <header class="gx-about-topic__head">
            <span class="gx-about-topic__icon gx-about-topic__icon--user" aria-hidden="true">
              <GxIcon name="user" :size="22" />
            </span>
            <div>
              <h2 class="gx-about-topic__title">如何注册</h2>
              <p class="gx-about-topic__subtitle">四步完成校内账号开通</p>
            </div>
          </header>
          <ol class="gx-about-steps">
            <li v-for="(step, index) in registerSteps" :key="step.title" class="gx-about-step">
              <span class="gx-about-step__num">{{ index + 1 }}</span>
              <div class="gx-about-step__body">
                <strong>{{ step.title }}</strong>
                <p>{{ step.desc }}</p>
              </div>
            </li>
          </ol>
          <p class="gx-about-topic__note">
            邀请码无效或已使用时，请向发放老师重新申请。
          </p>
        </article>

        <article class="gx-about-topic">
          <header class="gx-about-topic__head">
            <span class="gx-about-topic__icon gx-about-topic__icon--book" aria-hidden="true">
              <GxIcon name="book" :size="22" />
            </span>
            <div>
              <h2 class="gx-about-topic__title">板块说明</h2>
              <p class="gx-about-topic__subtitle">按主题浏览与发帖</p>
            </div>
          </header>
          <ul class="gx-about-boards">
            <li v-for="board in GX_NAV_ITEMS" :key="board.key" class="gx-about-board">
              <span class="gx-about-board__icon" aria-hidden="true">
                <GxIcon :name="board.icon" :size="18" />
              </span>
              <div class="gx-about-board__text">
                <RouterLink :to="`/community/boards/${board.key}`" class="gx-about-board__name">
                  {{ board.label }}
                </RouterLink>
                <p>{{ boardIntro[board.key] }}</p>
              </div>
            </li>
          </ul>
        </article>

        <article class="gx-about-topic">
          <header class="gx-about-topic__head">
            <span class="gx-about-topic__icon gx-about-topic__icon--edit" aria-hidden="true">
              <GxIcon name="edit" :size="22" />
            </span>
            <div>
              <h2 class="gx-about-topic__title">发帖与审核</h2>
              <p class="gx-about-topic__subtitle">发布内容前请知悉</p>
            </div>
          </header>
          <ul class="gx-about-features">
            <li v-for="tip in postTips" :key="tip.label" class="gx-about-feature">
              <span class="gx-about-feature__icon" aria-hidden="true">
                <GxIcon :name="tip.icon" :size="18" />
              </span>
              <div>
                <strong>{{ tip.label }}</strong>
                <p>{{ tip.desc }}</p>
              </div>
            </li>
          </ul>
          <RouterLink to="/community/posts/new" class="gx-about-topic__link">
            前往发帖 →
          </RouterLink>
        </article>

        <article class="gx-about-topic gx-about-topic--alert">
          <header class="gx-about-topic__head">
            <span class="gx-about-topic__icon gx-about-topic__icon--flag" aria-hidden="true">
              <GxIcon name="flag" :size="22" />
            </span>
            <div>
              <h2 class="gx-about-topic__title">违规与举报</h2>
              <p class="gx-about-topic__subtitle">共同维护文明交流环境</p>
            </div>
          </header>
          <ul class="gx-about-bullets">
            <li v-for="(line, i) in reportTips" :key="i">{{ line }}</li>
          </ul>
        </article>
      </div>
    </div>
  </div>
</template>
