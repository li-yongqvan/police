<template>
  <div class="admin-config" v-loading="loading">
    <h2>系统配置</h2>
    <el-card v-for="section in configSections" :key="section.name" style="margin-bottom: 16px">
      <template #header><h4>{{ section.name }}</h4></template>
      <div v-for="key in section.keys" :key="key" class="config-item">
        <label>{{ configLabels[key] || key }}</label>
        <el-input v-model="configValues[key]" :placeholder="configValues[key]" />
        <span class="config-desc">{{ configDescriptions[key] || '' }}</span>
      </div>
      <el-button type="primary" size="small" @click="saveSection(section.keys)">保存本组</el-button>
    </el-card>
    <el-button type="success" @click="saveAll">全部保存</el-button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getConfig, updateConfig } from '@/api/admin'

const configValues = ref({})
const loading = ref(false)

const configSections = [
  { name: '发帖配置', keys: ['post_max_title_length', 'post_max_content_length', 'daily_post_limit_per_user'] },
  { name: '等级要求', keys: ['post_requires_level', 'comment_requires_level', 'upload_requires_level'] },
  { name: '板块开关', keys: ['board_ai-learning_enabled', 'board_announcements_enabled', 'board_tech-help_enabled'] },
  { name: '审核策略', keys: ['sensitive_word_action', 'max_attachment_size_mb'] },
]

const configLabels = {
  'post_max_title_length': '帖子标题最大长度',
  'post_max_content_length': '帖子内容最大长度',
  'post_requires_level': '发帖所需最低等级',
  'comment_requires_level': '评论所需最低等级',
  'upload_requires_level': '上传附件所需最低等级',
  'max_attachment_size_mb': '附件最大大小（MB）',
  'board_ai-learning_enabled': 'AI学习交流区开关',
  'board_announcements_enabled': '协会公告区开关',
  'board_tech-help_enabled': '技术问答区开关',
  'sensitive_word_action': '敏感词处理方式',
  'daily_post_limit_per_user': '每用户每日发帖上限',
}

const configDescriptions = {
  'sensitive_word_action': 'pending_review = 待审核, reject = 直接拒绝',
}

async function fetchConfig() {
  loading.value = true
  try {
    const res = await getConfig()
    configValues.value = res.data.configs || {}
  } catch (e) {
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

async function saveSection(keys) {
  for (const key of keys) {
    try {
      await updateConfig(key, configValues.value[key])
    } catch (e) {
      ElMessage.error(`保存 ${key} 失败`)
      return
    }
  }
  ElMessage.success('配置已保存')
}

async function saveAll() {
  const allKeys = configSections.flatMap(s => s.keys)
  await saveSection(allKeys)
}

onMounted(fetchConfig)
</script>

<style scoped>
.admin-config h2 { margin-bottom: 16px; }
.config-item { margin-bottom: 12px; display: flex; align-items: center; gap: 12px; }
.config-item label { width: 180px; font-weight: 500; }
.config-item .el-input { width: 200px; }
.config-desc { color: #999; font-size: 12px; margin-left: 8px; }
</style>
