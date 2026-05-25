<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const rows = ref([])

onMounted(async () => {
  rows.value = await adminApi.getDailyStats(7)
})
</script>

<template>
  <section class="panel content-panel">
    <div class="section-title">
      <div>
        <p class="eyebrow">趋势统计</p>
        <h3>近 7 日运营数据</h3>
      </div>
    </div>
    <div class="stats-table-wrap">
      <table class="stats-table">
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
      <p v-if="!rows.length" class="login-hint">暂无历史统计数据，服务运行一段时间后会自动沉淀。</p>
    </div>
  </section>
</template>
