<script setup>
defineProps({
  panelTitle: { type: String, default: '' },
  panelSubtitle: { type: String, default: '' },
})
</script>

<template>
  <div class="gx-page gx-compose-page">
    <div class="gx-compose-shell">
      <aside class="gx-compose-rail" aria-label="发帖导航">
        <slot name="rail" />
      </aside>

      <section class="gx-compose-panel">
        <header class="gx-compose-panel__head">
          <slot name="panel-head">
            <div v-if="panelTitle" class="gx-compose-panel__head-main">
              <span class="gx-compose-panel__head-icon" aria-hidden="true">
                <slot name="panel-icon" />
              </span>
              <div>
                <h2 class="gx-compose-panel__title">{{ panelTitle }}</h2>
                <p v-if="panelSubtitle" class="gx-compose-panel__subtitle">{{ panelSubtitle }}</p>
              </div>
            </div>
          </slot>
        </header>

        <div class="gx-compose-panel__body">
          <slot />
        </div>

        <footer v-if="$slots.footer" class="gx-compose-panel__footer">
          <slot name="footer" />
        </footer>
      </section>
    </div>
  </div>
</template>
