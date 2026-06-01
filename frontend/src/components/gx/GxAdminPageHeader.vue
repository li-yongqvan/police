<script setup>
defineProps({
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>

<template>
  <header class="gx-admin-page-head mb-6">
    <h1 class="text-title text-gx-primary">{{ title }}</h1>
    <p v-if="description" class="mt-1 text-caption">{{ description }}</p>
    <div v-if="$slots.actions" class="gx-admin-page-head__actions mt-3 flex flex-wrap gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>
