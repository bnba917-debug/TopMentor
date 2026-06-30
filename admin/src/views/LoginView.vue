<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { getErrorMessage } from '@/api/client'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const username = ref('admin')
const password = ref('admin123')
const loading = ref(false)

async function onSubmit() {
  if (!username.value || !password.value) {
    ElMessage.warning('请输入账号和密码')
    return
  }
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.replace(redirect)
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <el-card class="login-card">
      <h1>TopMentor 运营后台</h1>
      <p class="subtitle">资质审核 · 课件管理 · 财务大盘</p>
      <el-form @submit.prevent="onSubmit">
        <el-form-item label="账号">
          <el-input v-model="username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" show-password autocomplete="current-password" />
        </el-form-item>
        <el-button type="primary" native-type="submit" :loading="loading" style="width: 100%">
          登录
        </el-button>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #eef2ff, #f5f7fa);
}

.login-card {
  width: 400px;
}

h1 {
  margin: 0 0 8px;
  font-size: 22px;
  text-align: center;
}

.subtitle {
  margin: 0 0 24px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
