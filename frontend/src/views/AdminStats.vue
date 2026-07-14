<script setup>
import { computed, onMounted, ref } from 'vue'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi } from '../api'

const rows = ref([])
const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '趋势统计' },
]

const newUserTotal = computed(() => rows.value.reduce((sum, row) => sum + (row.new_users ?? row.newUsers ?? 0), 0))
const newPostTotal = computed(() => rows.value.reduce((sum, row) => sum + (row.new_posts ?? row.newPosts ?? 0), 0))
const newCommentTotal = computed(() => rows.value.reduce((sum, row) => sum + (row.new_comments ?? row.newComments ?? 0), 0))

onMounted(async () => {
  rows.value = await adminApi.getDailyStats(7)
})
</script>

<template>
  <div class="gx-page gx-admin-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="趋势统计" title="近 7 日运营数据" />

    <section class="gx-admin-summary">
      <article class="gx-admin-summary__item">
        <span>新增用户</span>
        <strong>{{ newUserTotal }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>新增帖子</span>
        <strong>{{ newPostTotal }}</strong>
      </article>
      <article class="gx-admin-summary__item">
        <span>新增评论</span>
        <strong>{{ newCommentTotal }}</strong>
      </article>
    </section>

    <section class="gx-card">
      <div class="gx-section-head">
        <strong>每日趋势</strong>
        <span class="gx-muted">{{ rows.length }} 天记录</span>
      </div>
      <div class="gx-table-wrap">
        <table class="gx-table">
          <thead>
            <tr>
              <th>日期</th>
              <th>新增用户</th>
              <th>新增帖子</th>
              <th>新增评论</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in rows" :key="row.stat_date || row.date">
              <td>{{ row.stat_date || row.date }}</td>
              <td>{{ row.new_users ?? row.newUsers ?? 0 }}</td>
              <td>{{ row.new_posts ?? row.newPosts ?? 0 }}</td>
              <td>{{ row.new_comments ?? row.newComments ?? 0 }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p v-if="!rows.length" class="gx-empty">暂无历史统计数据，服务运行一段时间后会自动沉淀。</p>
    </section>
  </div>
</template>
