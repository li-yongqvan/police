<script setup>
import { onMounted, ref } from 'vue'
import { adminApi } from '../api'
import { useSessionStore } from '../stores/session'

const session = useSessionStore()
const words = ref([])
const newWord = ref('')
const category = ref('general')

async function load() {
  words.value = await adminApi.listSensitiveWords()
}

async function add() {
  if (!newWord.value.trim()) return
  await adminApi.addSensitiveWord(newWord.value.trim(), category.value)
  newWord.value = ''
  session.setFlash('敏感词已添加', 'success')
  await load()
}

async function remove(id) {
  await adminApi.deleteSensitiveWord(id)
  session.setFlash('已删除', 'success')
  await load()
}

onMounted(load)
</script>

<template>
  <div class="gx-page gx-admin-page">
    <header class="gx-page-head">
      <p class="gx-eyebrow">敏感词库</p>
      <h1>内容合规关键词</h1>
    </header>

    <section class="gx-card gx-form">
      <label>
        <span>词条</span>
        <input v-model="newWord" type="text" placeholder="输入敏感词" />
      </label>
      <label>
        <span>分类</span>
        <input v-model="category" type="text" />
      </label>
      <button type="button" class="gx-btn gx-btn--primary" @click="add">添加</button>
    </section>

    <section class="gx-card">
      <div class="gx-admin-list">
        <article v-for="item in words" :key="item.id" class="gx-admin-row">
          <div>
            <strong>{{ item.word }}</strong>
            <p class="gx-muted">{{ item.category || 'general' }}</p>
          </div>
          <button type="button" class="gx-btn gx-btn--danger" @click="remove(item.id)">删除</button>
        </article>
      </div>
    </section>
  </div>
</template>
