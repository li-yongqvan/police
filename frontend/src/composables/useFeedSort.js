import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export const FEED_SORT_OPTIONS = [
  { id: 'hot', label: '热门' },
  { id: 'new', label: '最新' },
  { id: 'featured', label: '精华' },
  { id: 'today', label: '今日' },
]

const VALID = new Set(FEED_SORT_OPTIONS.map((o) => o.id))

export function useFeedSort() {
  const route = useRoute()
  const router = useRouter()

  const sort = computed({
    get() {
      const q = route.query.sort
      return VALID.has(q) ? q : 'hot'
    },
    set(value) {
      const next = VALID.has(value) ? value : 'hot'
      router.replace({ query: { ...route.query, sort: next === 'hot' ? undefined : next } })
    },
  })

  return { sort, options: FEED_SORT_OPTIONS }
}

export function useFeedSortWatch(onChange) {
  const route = useRoute()
  watch(
    () => route.query.sort,
    () => onChange(true),
  )
}
