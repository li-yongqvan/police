<script setup>
defineProps({
  eyebrow: { type: String, default: '' },
  title: { type: String, required: true },
  description: { type: String, default: '' },
})
</script>

<template>
  <header class="gx-page-head">
    <p v-if="eyebrow" class="gx-eyebrow">{{ eyebrow }}</p>
    <h1 class="text-title text-gx-primary">{{ title }}</h1>
    <p v-if="description" class="gx-muted mt-1">{{ description }}</p>
    <div v-if="$slots.actions" class="gx-admin-page-head__actions mt-3 flex flex-wrap gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>