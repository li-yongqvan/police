<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'

const rows = ref([])

onMounted(async () => {
  rows.value = await adminApi.getDailyStats(7)
})
</script>

<template>
  <div class="gx-page gx-admin-page">
    <GxBreadcrumb :items="breadcrumbItems" />
    <GxAdminPageHeader eyebrow="趋势统计" title="近 7 日运营数据" />

    <section class="gx-card gx-table-wrap">
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
            <td>{{ row.new_users ?? 0 }}</td>
            <td>{{ row.new_posts ?? 0 }}</td>
            <td>{{ row.new_comments ?? 0 }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="!rows.length" class="gx-muted">暂无历史统计数据，服务运行一段时间后会自动沉淀。</p>
    </section>
  </div>
</template>
