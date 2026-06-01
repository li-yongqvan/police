<script setup>
import { RouterLink } from 'vue-router'

defineProps({
  items: { type: Array, required: true },
})
</script>

<template>
  <nav class="gx-breadcrumb" aria-label="面包屑导航">
    <ol class="gx-breadcrumb__list">
      <li v-for="(item, index) in items" :key="index" class="gx-breadcrumb__item">
        <RouterLink v-if="item.to && index < items.length - 1" :to="item.to" class="gx-breadcrumb__link">
          {{ item.label }}
        </RouterLink>
        <span v-else class="gx-breadcrumb__current" aria-current="page">{{ item.label }}</span>
        <span v-if="index < items.length - 1" class="gx-breadcrumb__sep" aria-hidden="true">/</span>
      </li>
    </ol>
  </nav>
</template>
