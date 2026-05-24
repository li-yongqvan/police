<template>
  <div class="admin-stats" v-loading="loading">
    <h2>数据统计</h2>

    <!-- Overview Cards -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <h3>注册用户数</h3>
            <p class="stat-value">{{ overview.total_users || 0 }}</p>
            <p class="stat-sub">今日新增 {{ overview.users_today || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <h3>总帖子数</h3>
            <p class="stat-value">{{ overview.total_posts || 0 }}</p>
            <p class="stat-sub">今日新增 {{ overview.posts_today || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <h3>总评论数</h3>
            <p class="stat-value">{{ overview.total_comments || 0 }}</p>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-card">
            <h3>封禁用户</h3>
            <p class="stat-value">{{ overview.banned_users || 0 }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Trend Chart -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="chart-header">
              <h3>趋势图</h3>
              <el-radio-group v-model="trendDays" size="small" @change="fetchTrendData">
                <el-radio-button :label="7">7天</el-radio-button>
                <el-radio-button :label="14">14天</el-radio-button>
                <el-radio-button :label="30">30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header><h3>用户等级分布</h3></template>
          <div ref="levelChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Board Activity Chart -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card>
          <template #header><h3>板块活跃度排行</h3></template>
          <div ref="boardChartRef" class="chart-container" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { getStatsOverview, getDailyStats } from '@/api/admin'
import * as echarts from 'echarts'

const loading = ref(false)
const overview = ref({})
const trendDays = ref(7)

const trendChartRef = ref(null)
const levelChartRef = ref(null)
const boardChartRef = ref(null)

let trendChart = null
let levelChart = null
let boardChart = null

async function fetchOverview() {
  loading.value = true
  try {
    const { data } = await getStatsOverview()
    overview.value = data || {}
  } catch (e) {
    console.error('Failed to fetch stats overview:', e)
  } finally {
    loading.value = false
  }
}

async function fetchTrendData() {
  try {
    const { data } = await getDailyStats(trendDays.value)
    const rows = data || []

    const dates = rows.map(r => r.stat_date)
    const users = rows.map(r => r.new_users || 0)
    const posts = rows.map(r => r.new_posts || 0)
    const comments = rows.map(r => r.new_comments || 0)

    if (!trendChart) {
      trendChart = echarts.init(trendChartRef.value)
    }
    trendChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['注册用户', '发帖量', '评论量'] },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dates, boundaryGap: false },
      yAxis: { type: 'value' },
      series: [
        { name: '注册用户', type: 'line', data: users, smooth: true, itemStyle: { color: '#409eff' } },
        { name: '发帖量', type: 'line', data: posts, smooth: true, itemStyle: { color: '#67c23a' } },
        { name: '评论量', type: 'line', data: comments, smooth: true, itemStyle: { color: '#e6a23c' } }
      ]
    })
  } catch (e) {
    console.error('Failed to fetch trend data:', e)
  }
}

function fetchLevelDistribution() {
  if (!levelChart) {
    levelChart = echarts.init(levelChartRef.value)
  }
  levelChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      label: { show: true, formatter: 'Lv{b}: {c}人' },
      data: [
        { value: 0, name: '0' },
        { value: 0, name: '1' },
        { value: 0, name: '2' },
        { value: 0, name: '3' },
        { value: 0, name: '4' },
        { value: 0, name: '5' }
      ]
    }]
  })
}

function fetchBoardActivity() {
  if (!boardChart) {
    boardChart = echarts.init(boardChartRef.value)
  }
  boardChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: ['示例板块'], inverse: true },
    series: [
      { name: '发帖量', type: 'bar', data: [0], itemStyle: { color: '#409eff' } },
      { name: '评论量', type: 'bar', data: [0], itemStyle: { color: '#67c23a' } }
    ]
  })
}

function handleResize() {
  trendChart?.resize()
  levelChart?.resize()
  boardChart?.resize()
}

onMounted(() => {
  fetchOverview()
  fetchTrendData()
  fetchLevelDistribution()
  fetchBoardActivity()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  levelChart?.dispose()
  boardChart?.dispose()
})
</script>

<style scoped>
.admin-stats h2 { margin-bottom: 20px; }
.stats-row { margin-bottom: 20px; }
.stat-card { text-align: center; }
.stat-card h3 { color: #909399; font-size: 14px; margin: 0 0 10px 0; }
.stat-value { font-size: 28px; font-weight: bold; color: #303133; margin: 0 0 8px 0; }
.stat-sub { font-size: 12px; color: #909399; margin: 0; }
.chart-row { margin-bottom: 20px; }
.chart-header { display: flex; justify-content: space-between; align-items: center; }
.chart-header h3 { margin: 0; }
.chart-container { width: 100%; height: 350px; }
</style>
