<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Search, Tag, Empty, Loading, Icon, NoticeBar } from 'vant'
import type { Mentor } from '@/api/types'
import { fetchMentors, getErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const mentors = ref<Mentor[]>([])
const keyword = ref('')

async function load() {
  loading.value = true
  try {
    const res = await fetchMentors({ school: keyword.value || undefined })
    mentors.value = res.list
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

function genderLabel(g: string) {
  if (g === 'male') return '男'
  if (g === 'female') return '女'
  return ''
}

function goDetail(id: number) {
  router.push({ name: 'mentor-detail', params: { id } })
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="学霸广场">
      <template #right>
        <span
          class="nav-link"
          @click="
            router.push(
              auth.isLoggedIn
                ? auth.isMentor
                  ? '/mentor'
                  : '/profile'
                : { name: 'welcome' },
            )
          "
        >
          {{ auth.isLoggedIn ? '我的' : '登录' }}
        </span>
      </template>
    </NavBar>

    <NoticeBar
      v-if="auth.isLoggedIn && !auth.isMentor"
      text="想成为陪伴学霸？提交入驻申请，审核通过即可接单。"
      color="#2563eb"
      background="#eef2ff"
      @click="router.push({ name: 'mentor-apply-status' })"
    />

    <NoticeBar
      v-else-if="!auth.isLoggedIn"
      text="首次使用？请先选择家长端或学霸端身份"
      color="#2563eb"
      background="#eef2ff"
      @click="router.push({ name: 'welcome' })"
    />

    <Search
      v-model="keyword"
      placeholder="搜索学校，如：清华"
      show-action
      @search="load"
      @clear="load"
    >
      <template #action>
        <div @click="load">搜索</div>
      </template>
    </Search>

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>

    <Empty v-else-if="mentors.length === 0" description="暂无学霸" />

    <div v-else class="list">
      <article v-for="m in mentors" :key="m.id" class="card" @click="goDetail(m.id)">
        <div class="card-head">
          <h3>{{ m.real_name }}</h3>
          <span class="school">{{ m.school_name }}</span>
        </div>
        <p class="major">{{ m.major }} · {{ genderLabel(m.gender) }}</p>
        <p class="score">{{ m.english_score }}</p>
        <div class="tags">
          <Tag v-for="tag in m.tags" :key="tag" plain type="primary" size="medium">{{ tag }}</Tag>
        </div>
        <div v-if="m.intro_video_url" class="video-hint">
          <Icon name="play-circle-o" /> 30 秒自荐视频
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.nav-link {
  font-size: 14px;
  color: var(--tm-primary);
}

.center-loading {
  margin-top: 80px;
  text-align: center;
}

.list {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
  background: var(--tm-card);
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-head h3 {
  margin: 0;
  font-size: 18px;
}

.school {
  font-size: 13px;
  color: var(--tm-primary);
}

.major,
.score {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--tm-muted);
}

.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.video-hint {
  margin-top: 10px;
  font-size: 12px;
  color: var(--tm-primary);
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
