<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const activeMenu = computed(() => route.path)

function onLogout() {
  auth.logout()
  router.push({ name: 'login' })
}
</script>

<template>
  <el-container class="layout">
    <el-aside width="220px" class="aside">
      <div class="brand">TopMentor</div>
      <el-menu :default-active="activeMenu" router>
        <el-menu-item index="/dashboard">
          <el-icon><DataAnalysis /></el-icon>
          <span>财务大盘</span>
        </el-menu-item>
        <el-menu-item index="/mentors">
          <el-icon><User /></el-icon>
          <span>学霸审核</span>
        </el-menu-item>
        <el-menu-item index="/courseware">
          <el-icon><Reading /></el-icon>
          <span>课件管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <span>运营后台</span>
        <div class="header-right">
          <span class="username">{{ auth.username }}</span>
          <el-button link type="primary" @click="onLogout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  min-height: 100vh;
}

.aside {
  background: #001529;
  color: #fff;
}

.brand {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.aside :deep(.el-menu) {
  border-right: none;
  background: #001529;
}

.aside :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.75);
}

.aside :deep(.el-menu-item.is-active) {
  background: #1890ff !important;
  color: #fff;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  color: #606266;
  font-size: 14px;
}
</style>
