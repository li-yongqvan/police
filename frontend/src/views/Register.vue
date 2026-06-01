<script setup>
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import GxAuthShell from '../components/gx/GxAuthShell.vue'
import { userApi } from '../api'
import { useSessionStore } from '../stores/session'

const router = useRouter()
const session = useSessionStore()
const form = ref({
  username: '',
  password: '',
  invitationCode: '',
  department: '',
  squad: '',
  grade: '',
})
const loading = ref(false)
const error = ref('')

async function submit() {
  if (!form.value.username.trim() || !form.value.password || !form.value.invitationCode.trim()) {
    error.value = '请填写学号/用户名、密码和邀请码'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const result = await userApi.register({
      username: form.value.username.trim(),
      password: form.value.password,
      invitationCode: form.value.invitationCode.trim(),
      department: form.value.department,
      squad: form.value.squad,
      grade: form.value.grade,
    })
    session.persistSession(result)
    session.setFlash('注册成功，欢迎加入 AI智联平台', 'success')
    router.push(session.routeAfterLogin(result.user))
  } catch (err) {
    error.value = err.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <GxAuthShell
    headline="加入校园论坛，开启交流之旅。"
    form-title="创建账号"
    form-hint="请填写学号、密码与邀请码；院系信息可在注册后补充"
    submit-label="注册并登录"
    :loading="loading"
    :error="error"
    @submit="submit"
  >
    <div class="gx-auth-field">
      <label for="reg-username">学号 / 用户名</label>
      <input id="reg-username" v-model="form.username" type="text" autocomplete="username" />
    </div>
    <div class="gx-auth-field">
      <label for="reg-password">密码</label>
      <input id="reg-password" v-model="form.password" type="password" autocomplete="new-password" />
    </div>
    <div class="gx-auth-field">
      <label for="reg-invite">邀请码</label>
      <input id="reg-invite" v-model="form.invitationCode" type="text" placeholder="学校发放的邀请码" />
    </div>
    <div class="gx-auth-field">
      <label for="reg-dept">院系</label>
      <input id="reg-dept" v-model="form.department" type="text" placeholder="例如：刑事科学技术学院" />
    </div>
    <div class="gx-auth-field">
      <label for="reg-squad">区队</label>
      <input id="reg-squad" v-model="form.squad" type="text" placeholder="例如：一区队" />
    </div>
    <div class="gx-auth-field">
      <label for="reg-grade">年级</label>
      <input id="reg-grade" v-model="form.grade" type="text" placeholder="例如：2024" />
    </div>

    <template #footer>
      <p>已有账号？</p>
      <RouterLink to="/">返回登录</RouterLink>
    </template>
  </GxAuthShell>
</template>
