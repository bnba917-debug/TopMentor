<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NavBar, CellGroup, Cell, Tag, Button, Image as VanImage, Loading } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { fetchMentorProfile, mediaUrl } from '@/api/client'
import type { MentorPortalProfile } from '@/api/types'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(true)
const profile = ref<MentorPortalProfile | null>(null)

async function loadProfile() {
  loading.value = true
  try {
    profile.value = await fetchMentorProfile()
    auth.syncMentorSummary(profile.value)
  } catch {
    profile.value = null
  } finally {
    loading.value = false
  }
}

function onLogout() {
  auth.logout()
  router.replace('/login')
}

onMounted(loadProfile)
</script>

<template>
  <div class="page mentor-center">
    <NavBar title="个人中心" />

    <Loading v-if="loading" vertical class="loading">加载中...</Loading>

    <template v-else>
      <div class="hero">
        <VanImage
          round
          width="72"
          height="72"
          fit="cover"
          :src="mediaUrl(profile?.avatar_url || auth.mentor?.avatar_url) || undefined"
        >
          <template #error>
            <div class="avatar-fallback">{{ (profile?.real_name || auth.mentor?.real_name || '?').slice(0, 1) }}</div>
          </template>
        </VanImage>
        <div class="hero-text">
          <h2>{{ profile?.real_name || auth.mentor?.real_name }}</h2>
          <p>{{ profile?.school_name || auth.mentor?.school_name }}</p>
          <div class="tags">
            <Tag v-if="profile?.is_verified" type="success">已认证</Tag>
            <Tag v-else type="warning">待认证</Tag>
            <Tag v-if="profile?.major" plain type="primary">{{ profile.major }}</Tag>
          </div>
        </div>
      </div>

      <CellGroup inset title="工作台">
        <Cell title="我的订单" is-link @click="router.push({ name: 'mentor-orders' })" />
        <Cell title="开放时段" is-link @click="router.push({ name: 'mentor-slots' })" />
        <Cell title="我的钱包" is-link @click="router.push({ name: 'mentor-wallet' })" />
      </CellGroup>

      <CellGroup inset title="资料与展示">
        <Cell title="编辑个人资料" label="姓名、学校、标签、简介" is-link @click="router.push({ name: 'mentor-profile' })" />
        <Cell
          title="宣传视频"
          :label="profile?.intro_video_url ? '已上传，点击管理' : '上传 30 秒自荐视频，吸引家长预约'"
          is-link
          @click="router.push({ name: 'mentor-profile', hash: '#video' })"
        />
        <Cell
          title="预览家长看到的页面"
          is-link
          @click="router.push({ name: 'mentor-detail', params: { id: auth.mentor?.id } })"
        />
      </CellGroup>

      <div class="footer">
        <Button round block plain type="danger" @click="onLogout">退出登录</Button>
      </div>
    </template>

    <p v-if="!loading && !profile" class="error">加载资料失败，请稍后重试</p>
  </div>
</template>

<style scoped>
.mentor-center {
  min-height: 100vh;
  background: var(--tm-bg, #f7f8fa);
}

.loading {
  margin-top: 80px;
}

.hero {
  display: flex;
  gap: 16px;
  align-items: center;
  padding: 24px 20px 8px;
}

.hero-text h2 {
  margin: 0 0 4px;
  font-size: 22px;
  color: var(--tm-primary, #2563eb);
}

.hero-text p {
  margin: 0 0 8px;
  color: var(--tm-muted, #64748b);
  font-size: 14px;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.avatar-fallback {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: #dbeafe;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
}

.footer {
  padding: 24px 16px 32px;
}

.error {
  text-align: center;
  color: #ef4444;
  padding: 16px;
}
</style>
