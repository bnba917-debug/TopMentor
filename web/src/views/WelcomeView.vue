<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Icon } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { loginPathForRole, savePortalRole } from '@/utils/portalRole'

const router = useRouter()
const auth = useAuthStore()

const features = [
  { icon: 'medal-o', text: '清华北大等名校学霸' },
  { icon: 'friends-o', text: '1 对 1 实时陪伴' },
  { icon: 'video-o', text: '在线视频互动课' },
]

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
    <div class="bg-glow bg-glow-a" aria-hidden="true" />
    <div class="bg-glow bg-glow-b" aria-hidden="true" />

    <div class="welcome-inner">
      <header class="brand">
        <div class="logo-badge">
          <Icon name="gem-o" size="22" color="#fff" />
        </div>
        <p class="logo-text">TopMentor</p>
        <h1>名校学霸，1 对 1 陪伴</h1>
        <p class="subtitle">为 6–14 岁孩子连接清华、北大等名校学长学姐，轻松预约在线课程</p>
      </header>

      <div class="features">
        <span v-for="item in features" :key="item.text" class="feature-chip">
          <Icon :name="item.icon" size="14" />
          {{ item.text }}
        </span>
      </div>

      <section class="cards" aria-label="选择身份">
        <button type="button" class="role-card parent" @click="enterParent">
          <div class="card-icon parent-icon">
            <Icon name="home-o" size="26" color="#4338ca" />
          </div>
          <div class="card-body">
            <div class="card-title-row">
              <strong>我是家长 / 学生</strong>
              <Icon name="arrow" class="card-arrow" />
            </div>
            <span class="desc">浏览学霸、预约课程、进入课室观看</span>
            <span class="action">登录 / 注册</span>
          </div>
        </button>

        <button type="button" class="role-card mentor" @click="enterMentor">
          <div class="card-icon mentor-icon">
            <Icon name="certificate" size="26" color="#6d28d9" />
          </div>
          <div class="card-body">
            <div class="card-title-row">
              <strong>我是学霸</strong>
              <Icon name="arrow" class="card-arrow" />
            </div>
            <span class="desc">接单授课、管理时段、填写成长报告</span>
            <span class="action mentor-action">登录 / 申请入驻</span>
          </div>
        </button>
      </section>

      <footer class="footer">
        <button type="button" class="browse-link" @click="browseWithoutLogin">
          <Icon name="search" size="16" />
          先逛逛学霸广场
          <span class="browse-note">无需登录</span>
        </button>
        <p class="copyright">TopPeerTalk · 让优质陪伴触手可及</p>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.welcome {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
  background: linear-gradient(165deg, #eef2ff 0%, #f8fafc 38%, #f5f6fa 100%);
  box-sizing: border-box;
}

.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.bg-glow-a {
  width: 280px;
  height: 280px;
  top: -80px;
  right: -60px;
  background: rgba(99, 102, 241, 0.22);
}

.bg-glow-b {
  width: 220px;
  height: 220px;
  bottom: 10%;
  left: -70px;
  background: rgba(167, 139, 250, 0.18);
}

.welcome-inner {
  position: relative;
  z-index: 1;
  max-width: 480px;
  margin: 0 auto;
  padding: 48px 20px 32px;
}

.brand {
  text-align: center;
  margin-bottom: 24px;
}

.logo-badge {
  width: 52px;
  height: 52px;
  margin: 0 auto 14px;
  border-radius: 16px;
  background: linear-gradient(135deg, #6366f1, #4f46e5);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 12px 28px rgba(79, 70, 229, 0.35);
}

.logo-text {
  margin: 0 0 10px;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--tm-primary);
}

.brand h1 {
  margin: 0 0 12px;
  font-size: 28px;
  line-height: 1.25;
  color: #0f172a;
  font-weight: 800;
}

.subtitle {
  margin: 0 auto;
  max-width: 320px;
  font-size: 14px;
  line-height: 1.65;
  color: var(--tm-muted);
}

.features {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
  margin-bottom: 28px;
}

.feature-chip {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(79, 70, 229, 0.12);
  font-size: 12px;
  color: #475569;
  backdrop-filter: blur(4px);
}

.cards {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.role-card {
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  text-align: left;
  border: 1px solid transparent;
  border-radius: 18px;
  padding: 18px;
  cursor: pointer;
  background: #fff;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.07);
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

.role-card:active {
  transform: scale(0.985);
}

.role-card.parent {
  border-color: #e0e7ff;
}

.role-card.parent:hover {
  border-color: #c7d2fe;
  box-shadow: 0 14px 36px rgba(79, 70, 229, 0.14);
}

.role-card.mentor {
  border-color: #ede9fe;
}

.role-card.mentor:hover {
  border-color: #ddd6fe;
  box-shadow: 0 14px 36px rgba(109, 40, 217, 0.12);
}

.card-icon {
  flex-shrink: 0;
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.parent-icon {
  background: linear-gradient(135deg, #eef2ff, #e0e7ff);
}

.mentor-icon {
  background: linear-gradient(135deg, #f5f3ff, #ede9fe);
}

.card-body {
  flex: 1;
  min-width: 0;
}

.card-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 6px;
}

.card-title-row strong {
  font-size: 17px;
  color: #0f172a;
}

.card-arrow {
  color: #cbd5e1;
  font-size: 14px;
}

.role-card:hover .card-arrow {
  color: var(--tm-primary);
}

.desc {
  display: block;
  font-size: 13px;
  color: #64748b;
  line-height: 1.55;
  margin-bottom: 10px;
}

.action {
  font-size: 13px;
  font-weight: 700;
  color: var(--tm-primary);
}

.mentor-action {
  color: #7c3aed;
}

.footer {
  margin-top: 28px;
  text-align: center;
}

.browse-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: none;
  background: none;
  color: #475569;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  padding: 10px 16px;
  border-radius: 999px;
  transition: background 0.15s ease, color 0.15s ease;
}

.browse-link:hover {
  background: rgba(255, 255, 255, 0.8);
  color: var(--tm-primary);
}

.browse-note {
  font-size: 11px;
  font-weight: 500;
  color: #94a3b8;
  padding: 2px 8px;
  border-radius: 999px;
  background: #f1f5f9;
}

.copyright {
  margin: 16px 0 0;
  font-size: 11px;
  color: #94a3b8;
  letter-spacing: 0.02em;
}

@media (min-width: 768px) {
  .welcome-inner {
    max-width: 720px;
    padding-top: 64px;
  }

  .brand h1 {
    font-size: 34px;
  }

  .subtitle {
    max-width: 420px;
    font-size: 15px;
  }

  .cards {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
  }

  .role-card {
    flex-direction: column;
    align-items: stretch;
    min-height: 200px;
  }

  .card-icon {
    width: 56px;
    height: 56px;
  }
}
</style>
