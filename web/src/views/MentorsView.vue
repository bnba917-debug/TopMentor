<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Search, Tag, Empty, Loading, Icon, NoticeBar } from 'vant'
import type { Mentor } from '@/api/types'
import { fetchMentors, getErrorMessage, mediaUrl } from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const mentors = ref<Mentor[]>([])
const total = ref(0)
const keyword = ref('')

const hotSchools = ['清华大学', '北京大学', '复旦大学']

const activeSchool = computed(() => keyword.value)

async function load() {
  loading.value = true
  try {
    const res = await fetchMentors({ school: keyword.value || undefined })
    mentors.value = res.list
    total.value = res.total
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

function avatarText(name: string) {
  return name?.trim().slice(0, 1) || '学'
}

function avatarStyle(id: number) {
  const hues = [234, 262, 199, 24, 330]
  const hue = hues[id % hues.length]
  return {
    background: `linear-gradient(135deg, hsl(${hue} 72% 58%), hsl(${(hue + 36) % 360} 68% 46%))`,
  }
}

function schoolShort(school: string) {
  return school.replace(/大学|学院|学校/g, '').slice(0, 4) || school
}

function pickSchool(school: string) {
  keyword.value = school === keyword.value ? '' : school
  load()
}

function goDetail(id: number) {
  router.push({ name: 'mentor-detail', params: { id } })
}

onMounted(load)
</script>

<template>
  <div class="mentors-page page">
    <div class="hero">
      <NavBar title="学霸广场" :border="false" class="hero-nav">
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

      <div class="hero-copy">
        <p class="hero-kicker">TopMentor · 名校 1 对 1</p>
        <h1>找到适合孩子的陪伴学霸</h1>
        <p class="hero-desc">清华北大等名校学长学姐，视频自荐 + 在线预约</p>
      </div>

      <div class="search-wrap">
        <Search
          v-model="keyword"
          shape="round"
          background="transparent"
          placeholder="搜索学校，如：清华"
          show-action
          @search="load"
          @clear="load"
        >
          <template #action>
            <button type="button" class="search-btn" @click="load">搜索</button>
          </template>
        </Search>
      </div>

      <div class="school-chips">
        <button
          v-for="school in hotSchools"
          :key="school"
          type="button"
          class="chip"
          :class="{ active: activeSchool === school }"
          @click="pickSchool(school)"
        >
          {{ schoolShort(school) }}
        </button>
      </div>
    </div>

    <div class="content">
      <NoticeBar
        v-if="auth.isLoggedIn && !auth.isMentor"
        class="notice"
        text="想成为陪伴学霸？提交入驻申请，审核通过即可接单。"
        color="#4338ca"
        background="#eef2ff"
        @click="router.push({ name: 'mentor-apply-status' })"
      />

      <NoticeBar
        v-else-if="!auth.isLoggedIn"
        class="notice"
        text="首次使用？请先选择家长端或学霸端身份"
        color="#4338ca"
        background="#eef2ff"
        @click="router.push({ name: 'welcome' })"
      />

      <div v-if="!loading && mentors.length > 0" class="list-meta">
        <span>共 {{ total }} 位认证学霸</span>
        <span v-if="keyword" class="filter-tag">筛选：{{ keyword }}</span>
      </div>

      <Loading v-if="loading" class="center-loading" vertical color="#4f46e5">加载中...</Loading>

      <Empty v-else-if="mentors.length === 0" image="search" description="暂无匹配的学霸，换个学校试试">
        <button type="button" class="reset-btn" @click="keyword = ''; load()">查看全部</button>
      </Empty>

      <div v-else class="list">
        <article v-for="m in mentors" :key="m.id" class="card" @click="goDetail(m.id)">
          <div class="card-top">
            <div class="avatar-wrap">
              <img v-if="m.avatar_url" :src="mediaUrl(m.avatar_url)" :alt="m.real_name" class="avatar-img" />
              <div v-else class="avatar-fallback" :style="avatarStyle(m.id)">
                {{ avatarText(m.real_name) }}
              </div>
              <span v-if="m.intro_video_url" class="video-dot" title="有自荐视频">
                <Icon name="play-circle-o" />
              </span>
            </div>

            <div class="card-main">
              <div class="name-row">
                <h3>{{ m.real_name }}</h3>
                <span class="gender">{{ genderLabel(m.gender) }}</span>
              </div>
              <p class="school-line">
                <span class="school-badge">{{ schoolShort(m.school_name) }}</span>
                {{ m.major }}
              </p>
              <p v-if="m.bio" class="bio">{{ m.bio }}</p>
            </div>
          </div>

          <div class="score-row">
            <Icon name="award-o" />
            <span>{{ m.english_score }}</span>
          </div>

          <div v-if="m.tags?.length" class="tags">
            <Tag v-for="tag in m.tags.slice(0, 4)" :key="tag" round plain type="primary" size="medium">
              {{ tag }}
            </Tag>
          </div>

          <div class="card-foot">
            <span v-if="m.intro_video_url" class="video-hint">
              <Icon name="video-o" /> 30 秒自荐视频
            </span>
            <span class="book-hint">查看详情 · 预约 <Icon name="arrow" /></span>
          </div>
        </article>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mentors-page {
  background: var(--tm-bg);
  padding-bottom: 32px;
}

.hero {
  background: linear-gradient(160deg, #eef2ff 0%, #f8fafc 45%, var(--tm-bg) 100%);
  padding-bottom: 8px;
}

.hero-nav {
  background: transparent;
}

.hero-nav :deep(.van-nav-bar__title) {
  font-weight: 700;
  color: #1e1b4b;
}

.nav-link {
  font-size: 14px;
  font-weight: 600;
  color: var(--tm-primary);
}

.hero-copy {
  padding: 4px 20px 12px;
}

.hero-kicker {
  margin: 0 0 6px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.06em;
  color: var(--tm-primary);
}

.hero-copy h1 {
  margin: 0 0 8px;
  font-size: 24px;
  line-height: 1.3;
  color: #0f172a;
}

.hero-desc {
  margin: 0;
  font-size: 14px;
  line-height: 1.5;
  color: var(--tm-muted);
}

.search-wrap {
  padding: 0 12px;
}

.search-wrap :deep(.van-search__content) {
  background: #fff;
  box-shadow: 0 4px 16px rgba(79, 70, 229, 0.08);
}

.search-btn {
  border: none;
  background: none;
  color: var(--tm-primary);
  font-size: 14px;
  font-weight: 600;
  padding: 0 4px;
  cursor: pointer;
}

.school-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px 16px 4px;
}

.chip {
  border: 1px solid #e0e7ff;
  background: #fff;
  color: #4338ca;
  font-size: 13px;
  padding: 6px 14px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.chip.active {
  background: var(--tm-primary);
  border-color: var(--tm-primary);
  color: #fff;
}

.content {
  padding: 0 16px;
}

.notice {
  margin-top: 12px;
  border-radius: 10px;
  overflow: hidden;
}

.list-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 16px 0 10px;
  font-size: 13px;
  color: var(--tm-muted);
}

.filter-tag {
  color: var(--tm-primary);
  background: #eef2ff;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
}

.center-loading {
  margin-top: 64px;
  text-align: center;
}

.reset-btn {
  margin-top: 12px;
  border: none;
  background: var(--tm-primary);
  color: #fff;
  padding: 10px 20px;
  border-radius: 999px;
  font-size: 14px;
  cursor: pointer;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.card {
  background: var(--tm-card);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.06);
  border: 1px solid rgba(79, 70, 229, 0.06);
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.card:active {
  transform: scale(0.99);
}

.card-top {
  display: flex;
  gap: 14px;
}

.avatar-wrap {
  position: relative;
  flex-shrink: 0;
}

.avatar-img,
.avatar-fallback {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  object-fit: cover;
}

.avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24px;
  font-weight: 700;
}

.video-dot {
  position: absolute;
  right: -4px;
  bottom: -4px;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #fff;
  color: var(--tm-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
}

.card-main {
  min-width: 0;
  flex: 1;
}

.name-row {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.name-row h3 {
  margin: 0;
  font-size: 18px;
  color: #0f172a;
}

.gender {
  font-size: 12px;
  color: var(--tm-muted);
}

.school-line {
  margin: 6px 0 0;
  font-size: 13px;
  color: #475569;
  line-height: 1.4;
}

.school-badge {
  display: inline-block;
  margin-right: 6px;
  padding: 2px 8px;
  border-radius: 6px;
  background: #eef2ff;
  color: var(--tm-primary);
  font-size: 12px;
  font-weight: 600;
}

.bio {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--tm-muted);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.score-row {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  padding: 6px 12px;
  border-radius: 10px;
  background: #fffbeb;
  color: #b45309;
  font-size: 13px;
  font-weight: 500;
}

.tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.card-foot {
  margin-top: 14px;
  padding-top: 12px;
  border-top: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  font-size: 12px;
}

.video-hint {
  color: var(--tm-primary);
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.book-hint {
  color: var(--tm-muted);
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

@media (min-width: 768px) {
  .hero-copy h1 {
    font-size: 28px;
  }

  .content {
    max-width: 960px;
    margin: 0 auto;
  }

  .list {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(79, 70, 229, 0.12);
  }
}
</style>
