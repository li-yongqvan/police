<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { formatDisplayDate } from '../../utils/displayName'

const props = defineProps({
  variant: { type: String, default: 'home' },
  stats: { type: Object, default: null },
  boards: { type: Array, default: () => [] },
  board: { type: Object, default: null },
})

const boardRules = computed(() => {
  const name = props.board?.name || '本板块'
  return [
    `发帖须与「${name}」主题相关，禁止无关广告与灌水。`,
    '尊重师长与同学，禁止人身攻击、泄露隐私或涉密内容。',
    '转载资料请注明出处，不得传播违规或未经核实信息。',
    '置顶与精华帖由管理员认定，请勿重复刷屏顶帖。',
    '违规内容将被删帖、禁言或移交学院处理。',
  ]
})

const todayMetrics = computed(() => {
  const s = props.stats
  if (!s) return []
  return [
    { label: '新增帖子', value: s.posts_today ?? 0 },
    { label: '新增回复', value: s.comments_today ?? s.replies_today ?? '—' },
    { label: '活跃用户', value: s.online_users ?? '—' },
    { label: '在线用户', value: s.online_users ?? '—' },
  ]
})

const statRows = computed(() => {
  const s = props.stats
  if (!s) return []
  return [
    { label: '注册用户', value: s.total_users ?? '—' },
    { label: '在线同学', value: s.online_users ?? '—' },
    { label: '帖子总数', value: s.total_posts ?? '—' },
    { label: '今日发帖', value: s.posts_today ?? '—' },
  ]
})

const activeBoards = computed(() =>
  [...props.boards].sort((a, b) => (b.postCount ?? 0) - (a.postCount ?? 0)).slice(0, 6),
)

const updatedAt = computed(() => formatDisplayDate(new Date().toISOString()))
</script>

<template>
  <div class="gx-feed-aside space-y-4">
    <template v-if="variant === 'board' && board">
      <div class="gx-panel">
        <h3 class="gx-panel__title">板块规则</h3>
        <ol class="gx-feed-aside__rules">
          <li v-for="(rule, i) in boardRules" :key="i">{{ rule }}</li>
        </ol>
        <RouterLink :to="`/community/boards/${board.slug}`" class="gx-feed-aside__more">
          查看更多规则 →
        </RouterLink>
      </div>

      <div v-if="todayMetrics.length" class="gx-panel">
        <h3 class="gx-panel__title">今日数据</h3>
        <dl class="gx-feed-aside__metrics">
          <div v-for="row in todayMetrics" :key="row.label">
            <dt>{{ row.label }}</dt>
            <dd>{{ row.value }}</dd>
          </div>
        </dl>
        <p class="gx-feed-aside__updated text-caption text-gx-muted">更新时间：{{ updatedAt }}</p>
      </div>
    </template>

    <template v-else>
      <div v-if="statRows.length" class="gx-panel">
        <h3 class="gx-panel__title">校园概况</h3>
        <dl class="gx-home-side-stats">
          <div v-for="row in statRows" :key="row.label">
            <dt>{{ row.label }}</dt>
            <dd>{{ row.value }}</dd>
          </div>
        </dl>
      </div>

      <div v-if="activeBoards.length" class="gx-panel">
        <h3 class="gx-panel__title">板块热度</h3>
        <ul class="gx-home-board-rank">
          <li v-for="b in activeBoards" :key="b.id">
            <RouterLink :to="`/community/boards/${b.slug}`">
              <span class="gx-home-board-rank__name">{{ b.name }}</span>
              <span class="gx-home-board-rank__count">{{ b.postCount ?? 0 }} 帖</span>
            </RouterLink>
          </li>
        </ul>
      </div>

      <div class="gx-panel">
        <h3 class="gx-panel__title">发帖须知</h3>
        <ul class="gx-feed-aside__tips">
          <li>标题简明，正文文明理性</li>
          <li>涉密、违规内容请勿发布</li>
          <li>部分帖子需审核后展示</li>
        </ul>
      </div>
    </template>
  </div>
</template>

<style scoped>
.gx-feed-aside__tips,
.gx-feed-aside__rules {
  margin: 0;
  padding-left: 1.1rem;
  font-size: 13px;
  color: var(--color-muted);
  line-height: 1.65;
}

.gx-feed-aside__rules {
  list-style: decimal;
}

.gx-feed-aside__more {
  display: inline-block;
  margin-top: 10px;
  font-size: 13px;
  font-weight: 600;
  color: var(--color-primary);
}

.gx-feed-aside__more:hover {
  color: var(--color-brand);
}

.gx-feed-aside__metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin: 0;
}

.gx-feed-aside__metrics dt {
  font-size: 12px;
  color: var(--color-muted);
}

.gx-feed-aside__metrics dd {
  margin: 4px 0 0;
  font-size: 22px;
  font-weight: 700;
  color: var(--color-primary);
}

.gx-feed-aside__updated {
  margin: 12px 0 0;
}
</style>
