<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showFailToast, showSuccessToast } from 'vant'
import { NavBar, Tag, Button, Loading, Empty, NoticeBar } from 'vant'
import type { Mentor, MentorSlot } from '@/api/types'
import {
  createBooking,
  fetchMentor,
  fetchMentorSlots,
  fetchLessonBalance,
  getErrorMessage,
  mediaUrl,
} from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const booking = ref(false)
const mentor = ref<Mentor | null>(null)
const slots = ref<MentorSlot[]>([])

const availableSlots = computed(() => slots.value.filter((s) => s.status === 0))

function formatSlotLabel(s: MentorSlot) {
  const start = s.start_time.slice(0, 5)
  const end = s.end_time.slice(0, 5)
  return `${s.slot_date} ${start}-${end}`
}

async function load() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const [m, list] = await Promise.all([fetchMentor(id), fetchMentorSlots(id)])
    mentor.value = m
    slots.value = list
    if (auth.isLoggedIn) {
      const balance = await fetchLessonBalance()
      auth.setAvailableLessons(balance.available_lessons)
    }
  } catch (e) {
    showFailToast(getErrorMessage(e))
    router.back()
  } finally {
    loading.value = false
  }
}

async function onBook(slot: MentorSlot) {
  if (!auth.isLoggedIn) {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
    return
  }

  const lessons = auth.user?.available_lessons ?? 0
  if (lessons < 1) {
    try {
      await showConfirmDialog({
        title: '课时不足',
        message: '当前可用课时为 0，是否前往购买课时包？',
        confirmButtonText: '去购买',
      })
      router.push({ name: 'packages' })
    } catch {
      /* cancelled */
    }
    return
  }

  try {
    await showConfirmDialog({
      title: '确认预约',
      message: `将扣除 1 课时预约：${formatSlotLabel(slot)}`,
    })
  } catch {
    return
  }

  booking.value = true
  try {
    const result = await createBooking(slot.id)
    auth.setAvailableLessons(result.available_lessons)
    showSuccessToast('预约成功')
    router.push({ name: 'orders' })
  } catch (e) {
    const msg = getErrorMessage(e)
    showFailToast(msg)
    if (msg.includes('课时')) {
      setTimeout(() => router.push({ name: 'packages' }), 1500)
    }
  } finally {
    booking.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="学霸详情" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>

    <template v-else-if="mentor">
      <div class="header">
        <div v-if="mentor.avatar_url" class="avatar-wrap">
          <img :src="mediaUrl(mentor.avatar_url)" alt="" class="avatar" />
        </div>
        <h1>{{ mentor.real_name }}</h1>
        <p>{{ mentor.school_name }} · {{ mentor.major }}</p>
        <p class="score">{{ mentor.english_score }}</p>
        <p v-if="mentor.bio" class="bio">{{ mentor.bio }}</p>
        <div class="tags">
          <Tag v-for="tag in mentor.tags" :key="tag" type="primary" plain>{{ tag }}</Tag>
        </div>
      </div>

      <div v-if="mentor.intro_video_url" class="video-box">
        <video controls playsinline :src="mediaUrl(mentor.intro_video_url)" />
        <p>30 秒自荐视频</p>
      </div>

      <div class="slots-section">
        <h2>可预约时间</h2>
        <NoticeBar
          v-if="auth.isLoggedIn"
          :text="`当前可用课时：${auth.user?.available_lessons ?? 0} 节`"
          left-icon="info-o"
          color="#4f46e5"
          background="#eef2ff"
          @click="router.push('/packages')"
        />
        <p class="hint">每节 45 分钟，预约即扣 1 课时 · <span class="buy-link" @click="router.push('/packages')">购买课时</span></p>

        <Empty v-if="availableSlots.length === 0" description="暂无可约时间" />

        <div v-else class="slot-list">
          <Button
            v-for="s in availableSlots"
            :key="s.id"
            block
            round
            plain
            type="primary"
            class="slot-btn"
            :loading="booking"
            @click="onBook(s)"
          >
            {{ formatSlotLabel(s) }}
          </Button>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.center-loading {
  margin-top: 80px;
  text-align: center;
}

.header {
  padding: 24px 16px;
  background: linear-gradient(180deg, #eef2ff, var(--tm-bg));
}

.header h1 {
  margin: 0 0 8px;
  font-size: 26px;
}

.header p {
  margin: 4px 0;
  color: var(--tm-muted);
}

.avatar-wrap {
  margin-bottom: 12px;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
}

.bio {
  margin-top: 10px !important;
  line-height: 1.6;
  color: #334155 !important;
  font-size: 14px;
}

.score {
  font-weight: 600;
  color: var(--tm-primary) !important;
}

.tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.video-box {
  margin: 16px;
  background: var(--tm-card);
  border-radius: 12px;
  overflow: hidden;
}

.video-box video {
  width: 100%;
  display: block;
  background: #000;
}

.video-box p {
  margin: 0;
  padding: 8px 12px;
  font-size: 12px;
  color: var(--tm-muted);
}

.slots-section {
  padding: 16px;
}

.slots-section h2 {
  margin: 0 0 4px;
  font-size: 18px;
}

.hint {
  margin: 8px 0 16px;
  font-size: 13px;
  color: var(--tm-muted);
}

.buy-link {
  color: var(--tm-primary);
  font-weight: 600;
}

.slot-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.slot-btn {
  margin: 0;
}
</style>
