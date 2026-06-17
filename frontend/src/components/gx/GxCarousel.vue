<script setup>
import { onMounted, onUnmounted, ref } from 'vue'

const slides = [
  { title: '忠诚为民', sub: '勤学砺警 交流共进' },
  { title: '勤学砺警', sub: '交流共进' },
]
const index = ref(0)
let timer

onMounted(() => {
  timer = window.setInterval(() => {
    index.value = (index.value + 1) % slides.length
  }, 6000)
})
onUnmounted(() => window.clearInterval(timer))

function pause() {
  window.clearInterval(timer)
}
</script>

<template>
  <section
    class="gx-carousel relative overflow-hidden rounded-gx-lg border border-gx-primary/15 bg-gradient-to-br from-gx-primary to-gx-primary-dark p-6 text-white shadow-[0_8px_32px_rgba(15,43,91,0.2)]"
    aria-label="论坛标语轮播"
    @mouseenter="pause"
  >
    <div class="pointer-events-none absolute inset-0 opacity-40" aria-hidden="true">
      <div class="absolute -right-8 -top-8 h-32 w-32 rounded-full bg-gx-gold/30 blur-2xl" />
      <div class="absolute -bottom-6 left-1/3 h-24 w-24 rounded-full bg-gx-accent/20 blur-xl" />
    </div>
    <div
      v-for="(slide, i) in slides"
      :key="i"
      class="gx-carousel__slide relative transition-opacity duration-500"
      :class="{ 'is-active': i === index }"
    >
      <p class="gx-carousel__title text-2xl font-bold tracking-wide md:text-3xl">{{ slide.title }}</p>
      <p class="gx-carousel__sub gx-carousel__sub--desktop mt-2 text-gx-gold/90">{{ slide.sub }}</p>
      <p class="gx-carousel__sub gx-carousel__sub--mobile mt-2 text-gx-gold/90">勤学砺警 · 交流共进</p>
    </div>
    <div class="relative mt-4 flex gap-2" role="tablist">
      <button
        v-for="(_, i) in slides"
        :key="i"
        type="button"
        class="h-2 w-2 rounded-full bg-white/30 transition-all"
        :class="{ 'w-6 bg-gx-gold': i === index }"
        :aria-label="`第 ${i + 1} 张`"
        @click="index = i"
      />
    </div>
  </section>
</template>

<style scoped>
.gx-carousel__slide {
  position: absolute;
  inset: 0;
  padding: 1.5rem;
  opacity: 0;
  pointer-events: none;
}
.gx-carousel__slide.is-active {
  position: relative;
  opacity: 1;
  pointer-events: auto;
}
.gx-carousel {
  min-height: 100px;
  max-height: 160px;
  padding: 1rem 1.25rem !important;
}
.gx-carousel__title {
  font-size: 1.25rem !important;
}
@media (max-width: 767.98px) {
  .gx-carousel__sub--desktop {
    display: none;
  }
}
@media (min-width: 768px) {
  .gx-carousel__sub--mobile {
    display: none;
  }
}
</style>
