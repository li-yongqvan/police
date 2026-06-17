<script setup>
defineProps({
  class: { type: null, default: undefined },
})
</script>

<template>
  <div class="inspira-grid-bg pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
    <div class="inspira-grid-bg__gradient" />
    <div class="inspira-grid-bg__grid" />
  </div>
</template>

<style scoped>
.inspira-grid-bg__gradient {
  position: absolute;
  inset: -20%;
  background:
    radial-gradient(circle at 20% 20%, rgba(212, 184, 106, 0.18), transparent 45%),
    radial-gradient(circle at 80% 10%, rgba(200, 22, 35, 0.08), transparent 40%),
    radial-gradient(circle at 50% 80%, rgba(15, 43, 91, 0.25), transparent 50%);
  animation: inspira-drift 18s ease-in-out infinite alternate;
}

.inspira-grid-bg__grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(15, 43, 91, 0.06) 1px, transparent 1px),
    linear-gradient(90deg, rgba(15, 43, 91, 0.06) 1px, transparent 1px);
  background-size: 32px 32px;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 75%);
}

@keyframes inspira-drift {
  from {
    transform: translate3d(0, 0, 0) scale(1);
  }
  to {
    transform: translate3d(-2%, 2%, 0) scale(1.05);
  }
}
</style>
