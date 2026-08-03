<script setup>
import { Shield } from 'lucide-vue-next'
import Alert from '../ui/Alert.vue'

defineProps({
  mode: { type: String, default: 'login' },
  eyebrow: { type: String, default: '广西警察学院 · 校内论坛' },
  headline: { type: String, default: '万千帖子，齐聚 AI智联平台。' },
  formTitle: { type: String, default: '' },
  formHint: { type: String, default: '' },
  submitLabel: { type: String, default: '登录' },
  loading: { type: Boolean, default: false },
  error: { type: String, default: '' },
  statusHint: { type: String, default: '' },
})

const emit = defineEmits(['submit'])
</script>

<template>
  <main class="gx-auth-spotify">
    <div class="gx-auth-spotify__bg" aria-hidden="true">
      <div class="gx-auth-spotify__grid" />
      <div class="gx-auth-spotify__orb gx-auth-spotify__orb--top" />
      <div class="gx-auth-spotify__orb gx-auth-spotify__orb--side" />
      <div class="gx-auth-spotify__scanline" />
    </div>

    <section class="gx-auth-spotify__body">
      <div class="gx-auth-spotify__intro">
        <div class="gx-auth-spotify__brand" aria-hidden="true">
          <Shield class="gx-auth-spotify__brand-icon" :size="32" stroke-width="1.75" />
        </div>
        <p class="gx-auth-spotify__eyebrow">{{ eyebrow }}</p>
        <h1 class="gx-auth-spotify__headline">{{ headline }}</h1>
        <p v-if="formHint && !formTitle" class="gx-auth-spotify__lede">{{ formHint }}</p>
      </div>

      <div class="gx-auth-spotify__panel">
        <header v-if="formTitle" class="gx-auth-spotify__panel-head">
          <h2 class="gx-auth-spotify__form-title">{{ formTitle }}</h2>
          <p v-if="formHint" class="gx-auth-spotify__form-hint">{{ formHint }}</p>
        </header>

        <form class="gx-auth-spotify__form" @submit.prevent="emit('submit')">
          <div class="gx-auth-spotify__fields">
            <slot />
          </div>

          <div v-if="statusHint || error" class="gx-auth-spotify__feedback">
            <p v-if="statusHint" class="gx-auth-spotify__status">{{ statusHint }}</p>
            <Alert v-if="error" variant="destructive" class="gx-auth-spotify__alert">{{ error }}</Alert>
          </div>

          <button type="submit" class="gx-auth-spotify__btn gx-auth-spotify__btn--primary" :disabled="loading">
            {{ loading ? '请稍候…' : submitLabel }}
          </button>
        </form>
      </div>

      <div v-if="$slots.footer" class="gx-auth-spotify__secondary">
        <slot name="footer" />
      </div>
    </section>
  </main>
</template>

<style scoped>
.gx-auth-spotify {
  --gx-auth-blue: var(--color-brand, #3d7eef);
  --gx-auth-blue-hover: #5b9cff;
  --gx-auth-bg: #030712;
  --gx-auth-panel: rgba(18, 24, 38, 0.92);
  --gx-auth-field: #1a2332;
  --gx-auth-border: rgba(100, 149, 237, 0.18);
  --gx-auth-gap-section: 40px;
  --gx-auth-gap-field: 22px;
  --gx-auth-panel-pad: clamp(24px, 5vw, 32px);
  position: relative;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: var(--gx-auth-bg);
  color: #fff;
}

.gx-auth-spotify__bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse 120% 80% at 50% -20%, rgba(29, 78, 216, 0.45), transparent 55%),
    radial-gradient(ellipse 70% 50% at 100% 50%, rgba(37, 99, 235, 0.15), transparent 50%),
    radial-gradient(ellipse 60% 40% at 0% 80%, rgba(61, 126, 239, 0.12), transparent 45%),
    linear-gradient(180deg, #050a14 0%, #030712 45%, #020617 100%);
}

.gx-auth-spotify__grid {
  position: absolute;
  inset: 0;
  opacity: 0.55;
  background-image:
    linear-gradient(rgba(61, 126, 239, 0.07) 1px, transparent 1px),
    linear-gradient(90deg, rgba(61, 126, 239, 0.07) 1px, transparent 1px);
  background-size: 56px 56px;
  mask-image: radial-gradient(ellipse 90% 70% at 50% 35%, #000 20%, transparent 75%);
}

.gx-auth-spotify__orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
}

.gx-auth-spotify__orb--top {
  top: -8%;
  left: 50%;
  width: min(520px, 90vw);
  height: min(320px, 45vh);
  transform: translateX(-50%);
  background: rgba(59, 130, 246, 0.35);
}

.gx-auth-spotify__orb--side {
  right: -12%;
  bottom: 18%;
  width: 280px;
  height: 280px;
  background: rgba(37, 99, 235, 0.22);
}

.gx-auth-spotify__scanline {
  position: absolute;
  inset: 0;
  opacity: 0.04;
  background: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(255, 255, 255, 0.15) 2px,
    rgba(255, 255, 255, 0.15) 3px
  );
}

.gx-auth-spotify__body {
  position: relative;
  z-index: 1;
  flex: 1;
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 440px;
  margin: 0 auto;
  padding: clamp(32px, 8vh, 56px) 20px 48px;
  overflow-y: auto;
}

.gx-auth-spotify__intro {
  margin-bottom: var(--gx-auth-gap-section);
  text-align: center;
}

.gx-auth-spotify__brand {
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: rgba(61, 126, 239, 0.12);
  border: 1px solid rgba(61, 126, 239, 0.35);
  box-shadow: 0 0 24px rgba(61, 126, 239, 0.25);
}

.gx-auth-spotify__brand-icon {
  color: #fff;
}

.gx-auth-spotify__eyebrow {
  margin: 0 0 12px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.45);
}

.gx-auth-spotify__headline {
  margin: 0;
  font-size: clamp(1.5rem, 4.5vw, 1.85rem);
  font-weight: 800;
  line-height: 1.25;
  letter-spacing: -0.02em;
}

.gx-auth-spotify__lede {
  margin: 14px auto 0;
  max-width: 22em;
  font-size: 15px;
  line-height: 1.65;
  color: rgba(255, 255, 255, 0.5);
}

.gx-auth-spotify__panel {
  padding: var(--gx-auth-panel-pad);
  border-radius: var(--radius-md, 12px);
  background: var(--gx-auth-panel);
  border: 1px solid var(--gx-auth-border);
  box-shadow:
    0 16px 48px rgba(0, 0, 0, 0.45),
    0 0 0 1px rgba(61, 126, 239, 0.06) inset;
  backdrop-filter: blur(12px);
}

.gx-auth-spotify__panel-head {
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--gx-auth-border);
}

.gx-auth-spotify__form-title {
  margin: 0;
  font-size: 1.35rem;
  font-weight: 700;
  line-height: 1.3;
  color: #fff;
}

.gx-auth-spotify__form-hint {
  margin: 10px 0 0;
  font-size: 14px;
  line-height: 1.65;
  color: rgba(255, 255, 255, 0.5);
}

.gx-auth-spotify__form {
  display: flex;
  flex-direction: column;
}

.gx-auth-spotify__fields {
  display: flex;
  flex-direction: column;
  gap: var(--gx-auth-gap-field);
}

.gx-auth-spotify__fields :deep(.gx-auth-field) {
  display: grid;
  gap: 10px;
}

.gx-auth-spotify__fields :deep(.gx-auth-field label) {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.02em;
  color: rgba(255, 255, 255, 0.7);
}

.gx-auth-spotify__fields :deep(.gx-auth-field input) {
  width: 100%;
  height: 52px;
  padding: 0 16px;
  border: 1px solid var(--gx-auth-border);
  border-radius: var(--radius-sm, 8px);
  background: var(--gx-auth-field);
  color: #fff;
  font-size: 16px;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.gx-auth-spotify__fields :deep(.gx-auth-field input::placeholder) {
  color: rgba(255, 255, 255, 0.32);
}

.gx-auth-spotify__fields :deep(.gx-auth-field input:focus) {
  outline: none;
  border-color: var(--gx-auth-blue);
  box-shadow: 0 0 0 3px rgba(61, 126, 239, 0.22);
}

.gx-auth-spotify__feedback {
  display: grid;
  gap: 12px;
  margin-top: 20px;
}

.gx-auth-spotify__status {
  margin: 0;
  font-size: 14px;
  line-height: 1.5;
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
}

.gx-auth-spotify__alert {
  margin: 0;
  border-radius: 10px;
  border-color: rgba(255, 120, 120, 0.35) !important;
  background: rgba(200, 22, 35, 0.2) !important;
  color: #ffc9c9 !important;
}

.gx-auth-spotify__btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 44px;
  margin-top: 32px;
  padding: 0 28px;
  border: none;
  border-radius: var(--radius-md, 12px);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition: transform 0.12s, background 0.15s, opacity 0.15s;
}

.gx-auth-spotify__btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.gx-auth-spotify__btn--primary {
  background: var(--gx-auth-blue);
  color: #fff;
}

.gx-auth-spotify__btn--primary:hover:not(:disabled) {
  background: var(--gx-auth-blue-hover);
  transform: scale(1.01);
}

.gx-auth-spotify__secondary {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin-top: var(--gx-auth-gap-section);
  padding-top: 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.gx-auth-spotify__secondary :deep(a) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 54px;
  padding: 0 28px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.55);
  background: transparent;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  text-decoration: none;
  transition: background 0.15s, border-color 0.15s;
}

.gx-auth-spotify__secondary :deep(a:hover) {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.85);
}

.gx-auth-spotify__secondary :deep(p) {
  margin: 0;
  max-width: 20em;
  font-size: 14px;
  line-height: 1.65;
  text-align: center;
  color: rgba(255, 255, 255, 0.42);
}

@media (min-width: 768px) {
  .gx-auth-spotify__body {
    padding-left: 28px;
    padding-right: 28px;
  }
}

@media (max-width: 479px) {
  .gx-auth-spotify {
    justify-content: flex-start;
    padding-top: env(safe-area-inset-top);
    padding-bottom: env(safe-area-inset-bottom);
  }
  .gx-auth-spotify__body {
    max-width: none;
    padding: max(20px, env(safe-area-inset-top)) 16px calc(24px + env(safe-area-inset-bottom));
  }
  .gx-auth-spotify__intro {
    margin-bottom: 24px;
  }
  .gx-auth-spotify__headline {
    font-size: 1.35rem;
  }
  .gx-auth-spotify__panel {
    padding: 20px 16px;
    border-radius: var(--radius-md, 12px);
  }
  .gx-auth-spotify__btn {
    min-height: 44px;
    margin-top: 24px;
  }
  .gx-auth-spotify__secondary {
    margin-top: 24px;
    padding-top: 20px;
  }
}

@media (max-height: 480px) and (orientation: landscape) {
  .gx-auth-spotify {
    justify-content: flex-start;
  }
  .gx-auth-spotify__body {
    max-width: 440px;
    padding: max(12px, env(safe-area-inset-top)) 20px calc(18px + env(safe-area-inset-bottom));
  }
  .gx-auth-spotify__intro {
    margin-bottom: 12px;
  }
  .gx-auth-spotify__brand {
    width: 40px;
    height: 40px;
    margin-bottom: 8px;
  }
  .gx-auth-spotify__brand-icon {
    width: 24px;
    height: 24px;
  }
  .gx-auth-spotify__eyebrow {
    margin-bottom: 4px;
    font-size: 11px;
  }
  .gx-auth-spotify__headline {
    font-size: 1.15rem;
    line-height: 1.18;
  }
  .gx-auth-spotify__panel {
    padding: 14px 16px;
  }
  .gx-auth-spotify__panel-head {
    margin-bottom: 12px;
    padding-bottom: 10px;
  }
  .gx-auth-spotify__form-title {
    font-size: 1.05rem;
  }
  .gx-auth-spotify__form-hint {
    margin-top: 4px;
    line-height: 1.35;
  }
  .gx-auth-spotify__fields {
    gap: 10px;
  }
  .gx-auth-spotify__fields :deep(.gx-auth-field) {
    gap: 6px;
  }
  .gx-auth-spotify__fields :deep(.gx-auth-field input) {
    height: 44px;
  }
  .gx-auth-spotify__feedback {
    margin-top: 12px;
  }
  .gx-auth-spotify__btn {
    margin-top: 14px;
  }
  .gx-auth-spotify__secondary {
    margin-top: 12px;
    padding-top: 12px;
  }
}
</style>
