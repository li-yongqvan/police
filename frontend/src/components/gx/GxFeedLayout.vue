<script setup>
defineProps({
  showAside: { type: Boolean, default: true },
})
</script>

<template>
  <div class="gx-feed-layout" :class="{ 'gx-feed-layout--no-aside': !showAside }">
    <main class="gx-feed-layout__main">
      <slot name="header" />
      <slot name="sort" />
      <div class="gx-feed-layout__stream">
        <slot />
      </div>
    </main>
    <aside v-if="showAside" class="gx-feed-layout__aside">
      <slot name="aside" />
    </aside>
  </div>
</template>
