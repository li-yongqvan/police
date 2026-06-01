<script setup>
import { computed, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import GxIcon from './GxIcon.vue'
import { GX_MOBILE_TABS, getLastBoardPath, isNavActive, rememberBoardSlug } from '../../composables/useGxNav'

const route = useRoute()

const tabs = computed(() =>
  GX_MOBILE_TABS.map((tab) =>
    tab.matchBoard ? { ...tab, to: getLastBoardPath() } : tab,
  ),
)

watch(
  () => route.params.slug,
  (slug) => {
    if (route.name === 'board' && slug) rememberBoardSlug(slug)
  },
  { immediate: true },
)
</script>

<template>
  <nav class="gx-tabbar" aria-label="底部导航">
    <RouterLink
      v-for="tab in tabs"
      :key="tab.name"
      :to="tab.to"
      class="gx-tabbar__item flex flex-col items-center gap-0.5 py-2 text-[11px]"
      :class="{ 'is-active': isNavActive(route, tab) }"
    >
      <GxIcon :name="tab.icon" :size="22" />
      <span>{{ tab.label }}</span>
    </RouterLink>
  </nav>
</template>

<style scoped>
.gx-tabbar__item {
  color: var(--color-muted);
  transition: color 0.15s ease;
}
.gx-tabbar__item.is-active {
  color: var(--color-primary);
  font-weight: 600;
}
</style>
