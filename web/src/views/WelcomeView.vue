<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Button } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { loginPathForRole, savePortalRole } from '@/utils/portalRole'

const router = useRouter()
const auth = useAuthStore()

function enterParent() {
  savePortalRole('parent')
  if (auth.isLoggedIn) {
    router.replace(auth.isMentor ? { name: 'mentor-home' } : { name: 'mentors' })
    return
  }
  router.push(loginPathForRole('parent'))
}

function enterMentor() {
  savePortalRole('mentor')
  if (auth.isLoggedIn) {
    router.replace(auth.isMentor ? { name: 'mentor-home' } : { name: 'mentor-apply-status' })
    return
  }
  router.push(loginPathForRole('mentor', '/mentor/apply/status'))
}

function browseWithoutLogin() {
  savePortalRole('parent')
  router.push({ name: 'mentors' })
}
</script>

<template>
  <div class="welcome page">
    <div class="brand">
      <p class="logo">TopMentor</p>
      <h1>名校学霸，1 对 1 陪伴</h1>
      <p class="subtitle">请选择您的身份，我们将带您进入对应功能</p>
    </div>

    <div class="cards">
      <button type="button" class="role-card parent" @click="enterParent">
        <span class="icon">👨‍👩‍👧</span>
        <strong>我是家长 / 学生</strong>
        <span class="desc">浏览学霸、预约课程、进入课室观看</span>
        <span class="action">登录 / 注册 →</span>
      </button>

      <button type="button" class="role-card mentor" @click="enterMentor">
        <span class="icon">🎓</span>
        <strong>我是学霸</strong>
        <span class="desc">接单授课、管理时段、填写成长报告</span>
        <span class="action">登录 / 申请入驻 →</span>
      </button>
    </div>

    <div class="footer">
      <Button round block plain @click="browseWithoutLogin">先逛逛学霸广场（无需登录）</Button>
    </div>
  </div>
</template>

<style scoped>
.welcome {
  min-height: 100vh;
  padding: 32px 20px 24px;
  background: linear-gradient(180deg, #eef2ff 0%, #f8fafc 45%, var(--tm-bg, #f7f8fa) 100%);
  box-sizing: border-box;
}

.brand {
  text-align: center;
  margin-bottom: 28px;
}

.logo {
  margin: 0 0 8px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.08em;
  color: var(--tm-primary, #2563eb);
}

.brand h1 {
  margin: 0 0 10px;
  font-size: 26px;
  color: #0f172a;
}

.subtitle {
  margin: 0;
  font-size: 14px;
  color: var(--tm-muted, #64748b);
  line-height: 1.5;
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.role-card {
  width: 100%;
  text-align: left;
  border: none;
  border-radius: 16px;
  padding: 20px;
  cursor: pointer;
  background: #fff;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.role-card:active {
  transform: scale(0.98);
}

.role-card.parent {
  border: 2px solid #dbeafe;
}

.role-card.mentor {
  border: 2px solid #ede9fe;
}

.icon {
  font-size: 32px;
  display: block;
  margin-bottom: 10px;
}

.role-card strong {
  display: block;
  font-size: 18px;
  color: #0f172a;
  margin-bottom: 6px;
}

.desc {
  display: block;
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
  margin-bottom: 12px;
}

.action {
  font-size: 14px;
  font-weight: 600;
  color: var(--tm-primary, #2563eb);
}

.role-card.mentor .action {
  color: #7c3aed;
}

.footer {
  margin-top: 24px;
}
</style>
