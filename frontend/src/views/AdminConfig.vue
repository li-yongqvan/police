<script setup>
import { computed, onMounted, ref } from 'vue'
import { LayoutGrid, Save, Settings, ShieldCheck } from 'lucide-vue-next'
import GxBreadcrumb from '../components/gx/GxBreadcrumb.vue'
import GxAdminPageHeader from '../components/gx/GxAdminPageHeader.vue'
import { adminApi, forumApi } from '../api'
import { formatApiError } from '../api/errors'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const config = ref(null)
const boards = ref([])
const loading = ref(true)
const saving = ref(false)
const error = ref('')

const breadcrumbItems = [
  { label: '管理后台', to: '/admin' },
  { label: '系统配置' },
]

const enabledBoardCount = computed(() => {
  if (!config.value) return 0
  return boards.value.filter((board) => config.value.boardSwitches?.[board.id]).length
})

const moderationLabel = computed(() => (config.value?.moderationMode === 'manual' ? '人工审核' : '自动审核'))

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [nextConfig, nextBoards] = await Promise.all([adminApi.getConfig(), forumApi.getBoards(true)])
    config.value = nextConfig
    boards.value = nextBoards
  } catch (err) {
    error.value = formatApiError(err, '系统配置加载失败')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!config.value || saving.value) return
  saving.value = true
  try {
    config.value = await adminApi.updateConfig(config.value)
    window.dispatchEvent(new CustomEvent('forum-config-updated'))
    session.setFlash('系统配置已保存。', 'success')
  } catch (err) {
    session.setFlash(formatApiError(err, '系统配置保存失败'), 'info')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page gx-config-page">
    <GxBreadcrumb :items="breadcrumbItems" />

    <GxAdminPageHeader
      eyebrow="系统配置"
      title="基础控制能力"
      description="统一维护发帖准入、审核策略和板块开放状态。"
    />

    <p v-if="loading" class="gx-muted gx-config-loading">正在加载系统配置...</p>
    <section v-else-if="error" class="gx-config-error">
      <p>{{ error }}</p>
      <button type="button" class="gx-btn gx-btn--secondary" @click="load">重试</button>
    </section>

    <template v-else-if="config">
      <section class="gx-config-summary" aria-label="配置状态概览">
        <article class="gx-config-summary__item">
          <Settings :size="18" />
          <span>发帖状态</span>
          <strong>{{ config.postingEnabled ? '已开启' : '已关闭' }}</strong>
        </article>
        <article class="gx-config-summary__item">
          <ShieldCheck :size="18" />
          <span>审核模式</span>
          <strong>{{ moderationLabel }}</strong>
        </article>
        <article class="gx-config-summary__item">
          <LayoutGrid :size="18" />
          <span>开放板块</span>
          <strong>{{ enabledBoardCount }} / {{ boards.length }}</strong>
        </article>
      </section>

      <section class="gx-config-grid">
        <div class="gx-config-panel">
          <div class="gx-config-panel__head">
            <div>
              <p class="gx-eyebrow">Posting</p>
              <h2>发布权限</h2>
            </div>
            <span :class="['gx-config-status', config.postingEnabled ? 'is-on' : 'is-off']">
              {{ config.postingEnabled ? '允许发帖' : '暂停发帖' }}
            </span>
          </div>

          <label class="gx-config-control">
            <span>
              <strong>开启发帖</strong>
              <small>关闭后普通用户不能继续发布新帖子。</small>
            </span>
            <input v-model="config.postingEnabled" class="gx-switch-input" type="checkbox" />
          </label>
        </div>

        <div class="gx-config-panel">
          <div class="gx-config-panel__head">
            <div>
              <p class="gx-eyebrow">Moderation</p>
              <h2>审核策略</h2>
            </div>
            <span class="gx-config-status is-neutral">{{ moderationLabel }}</span>
          </div>

          <label class="gx-config-field">
            <span>
              <strong>审核模式</strong>
              <small>自动审核会直接处理敏感内容，人工审核会进入待审队列。</small>
            </span>
            <select v-model="config.moderationMode">
              <option value="auto">自动审核</option>
              <option value="manual">人工审核</option>
            </select>
          </label>
        </div>
      </section>

      <section class="gx-config-panel gx-config-panel--wide">
        <div class="gx-config-panel__head">
          <div>
            <p class="gx-eyebrow">Boards</p>
            <h2>板块开关</h2>
          </div>
          <span class="gx-config-status is-neutral">{{ enabledBoardCount }} 个已开放</span>
        </div>

        <div class="gx-board-switch-grid">
          <label v-for="board in boards" :key="board.id" class="gx-board-switch">
            <span>
              <strong>{{ board.name }}</strong>
              <small>{{ board.description || '暂无板块说明' }}</small>
            </span>
            <input v-model="config.boardSwitches[board.id]" class="gx-switch-input" type="checkbox" />
          </label>
        </div>
      </section>

      <div class="gx-config-actions">
        <button type="button" class="gx-btn gx-btn--primary" :disabled="saving" @click="save">
          <Save :size="16" />
          {{ saving ? '保存中...' : '保存配置' }}
        </button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.gx-config-page {
  width: 100%;
  max-width: none;
}

.gx-config-loading,
.gx-config-error {
  margin: 0;
}

.gx-config-error {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px;
  border: 1px solid #f1b8b8;
  border-radius: var(--radius-md);
  background: #fff5f5;
  color: #9f1d1d;
}

.gx-config-error p {
  margin: 0;
}

.gx-config-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.gx-config-summary__item {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 4px 10px;
  align-items: center;
  min-height: 82px;
  padding: 14px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}

.gx-config-summary__item svg {
  grid-row: span 2;
  color: var(--color-accent);
}

.gx-config-summary__item span {
  font-size: 13px;
  color: var(--color-muted);
}

.gx-config-summary__item strong {
  min-width: 0;
  font-size: 18px;
  color: var(--color-primary);
}

.gx-config-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.gx-config-panel {
  display: grid;
  align-content: start;
  gap: 18px;
  padding: 18px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  box-shadow: 0 8px 22px rgba(15, 43, 91, 0.06);
}

.gx-config-panel--wide {
  margin-bottom: 16px;
}

.gx-config-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--color-border);
}

.gx-config-panel__head h2 {
  margin: 2px 0 0;
  color: var(--color-primary);
  font-size: 18px;
  line-height: 1.3;
}

.gx-config-panel__head .gx-eyebrow {
  margin: 0;
  font-size: 12px;
  letter-spacing: 0;
}

.gx-config-status {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 28px;
  padding: 0 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border);
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}

.gx-config-status.is-on {
  border-color: #8ac29b;
  background: #edf8f0;
  color: #1f7a3d;
}

.gx-config-status.is-off {
  border-color: #e7b0b0;
  background: #fff2f2;
  color: #a12a2a;
}

.gx-config-status.is-neutral {
  border-color: #d7deea;
  background: #f3f6fb;
  color: var(--color-primary);
}

.gx-config-control,
.gx-config-field,
.gx-board-switch {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px;
  align-items: center;
}

.gx-config-control strong,
.gx-config-field strong,
.gx-board-switch strong {
  display: block;
  color: var(--color-text);
  font-size: 15px;
  line-height: 1.35;
}

.gx-config-control small,
.gx-config-field small,
.gx-board-switch small {
  display: block;
  margin-top: 4px;
  color: var(--color-muted);
  font-size: 13px;
  line-height: 1.45;
}

.gx-config-field select {
  width: min(220px, 42vw);
  min-height: 42px;
  padding: 0 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-surface);
  color: var(--color-primary);
  font-weight: 600;
}

.gx-switch-input {
  position: relative;
  width: 52px;
  height: 30px;
  flex: 0 0 auto;
  margin: 0;
  border: 1px solid #c8d1de;
  border-radius: 999px;
  appearance: none;
  background: #d8dee8;
  cursor: pointer;
  transition:
    background-color 0.16s ease,
    border-color 0.16s ease;
}

.gx-switch-input::after {
  content: '';
  position: absolute;
  top: 3px;
  left: 3px;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  background: #fff;
  box-shadow: 0 2px 6px rgba(15, 43, 91, 0.22);
  transition: transform 0.16s ease;
}

.gx-switch-input:checked {
  border-color: var(--color-accent);
  background: var(--color-accent);
}

.gx-switch-input:checked::after {
  transform: translateX(22px);
}

.gx-switch-input:focus-visible,
.gx-config-field select:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px rgba(15, 43, 91, 0.16);
}

.gx-board-switch-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.gx-board-switch {
  min-height: 88px;
  padding: 14px;
  border: 1px solid #e4e9f2;
  border-radius: var(--radius-md);
  background: #f8fafc;
}

.gx-config-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}

@media (max-width: 760px) {
  .gx-config-summary,
  .gx-config-grid,
  .gx-board-switch-grid {
    grid-template-columns: 1fr;
  }

  .gx-config-panel {
    padding: 16px;
  }

  .gx-config-panel__head,
  .gx-config-control,
  .gx-config-field,
  .gx-board-switch {
    gap: 12px;
  }

  .gx-config-field {
    grid-template-columns: 1fr;
  }

  .gx-config-field select {
    width: 100%;
  }

  .gx-config-actions {
    justify-content: stretch;
  }

  .gx-config-actions .gx-btn {
    width: 100%;
  }
}
</style>
