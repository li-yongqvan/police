<template>
  <footer class="gx-footer">
    <div class="gx-footer__pattern" aria-hidden="true" />
    <div class="gx-footer__inner">
      <p class="gx-footer__org">AI智联平台 · 广西警察学院校内交流论坛</p>
      <p class="gx-footer__links">
        <a href="mailto:feedback@gx-police-campus.local">意见反馈</a>
        <span>·</span>
        <RouterLink to="/community/about">关于本站</RouterLink>
      </p>
      <p class="gx-footer__legal">
        本论坛为广西警察学院校内交流平台，视觉风格借鉴警校文化，<strong>不使用</strong>国家警徽等法定标识。所有图标均为原创简化设计，不代表公安机关官方形象。
      </p>
    </div>
  </footer>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>
