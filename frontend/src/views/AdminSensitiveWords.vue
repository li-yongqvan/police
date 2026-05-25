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
  <section class="panel form-panel">
    <p class="eyebrow">敏感词库</p>
    <h2>内容合规关键词</h2>
    <div class="form-grid">
      <label>
        <span>词条</span>
        <input v-model="newWord" type="text" placeholder="输入敏感词" />
      </label>
      <label>
        <span>分类</span>
        <input v-model="category" type="text" />
      </label>
    </div>
    <button class="primary-button" @click="add">添加</button>

    <div class="user-list" style="margin-top: 22px">
      <article v-for="item in words" :key="item.id" class="user-row">
        <div>
          <strong>{{ item.word }}</strong>
          <p>{{ item.category || 'general' }}</p>
        </div>
        <button class="danger-button" @click="remove(item.id)">删除</button>
      </article>
    </div>
  </section>
</template>
