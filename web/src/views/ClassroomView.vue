<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NavBar, Button, Tag, Loading, NoticeBar } from 'vant'
import { useClassroom } from '@/composables/useClassroom'

const route = useRoute()
const router = useRouter()
const orderId = route.params.orderId as string
const role = (route.query.role as 'user' | 'mentor') || 'user'

const {
  joining,
  inRoom,
  mockMode,
  receiveOnly,
  mediaHint,
  lessonStarted,
  peerOnline,
  peerLabel,
  statusHint,
  timerLabel,
  timerDisplay,
  mentorName,
  errorMsg,
  localVideoRef,
  remoteVideoRef,
  enter,
  leave,
} = useClassroom(orderId, role)

async function onLeave() {
  await leave()
  router.replace(role === 'mentor' ? '/mentor/orders' : '/orders')
}

onMounted(() => {
  enter()
})
</script>

<template>
  <div class="classroom page">
    <NavBar :title="mentorName ? `课室 · ${mentorName}` : '互动课室'" />

    <div class="timer-bar">
      <span>{{ timerLabel }}</span>
      <strong :class="{ waiting: !lessonStarted }">{{ timerDisplay }}</strong>
      <Tag v-if="peerOnline" type="success">{{ peerLabel }}在线</Tag>
      <Tag v-else-if="inRoom && lessonStarted" type="warning">{{ peerLabel }}未进房</Tag>
      <Tag v-if="mockMode" type="warning">模拟音视频</Tag>
      <Tag v-else-if="receiveOnly" type="primary">仅观看</Tag>
    </div>

    <NoticeBar
      v-if="statusHint"
      class="status-hint"
      wrapable
      :scrollable="false"
      :text="statusHint"
    />

    <NoticeBar v-if="mediaHint" wrapable :scrollable="false" :text="mediaHint" />

    <Loading v-if="joining" vertical class="loading">进房中...</Loading>

    <p v-if="errorMsg" class="error">{{ errorMsg }}</p>

    <div v-if="inRoom" class="video-grid">
      <div class="video-box">
        <div ref="localVideoRef" class="video-inner" />
        <p v-if="mockMode || receiveOnly" class="mock-label">
          我（{{ role === 'user' ? '家长' : '学霸' }}）{{ receiveOnly ? '· 本机无摄像头' : '' }}
        </p>
      </div>
      <div class="video-box remote">
        <div ref="remoteVideoRef" class="video-inner" />
        <div v-if="!peerOnline && lessonStarted" class="peer-waiting">
          <p>{{ peerLabel }}尚未进房</p>
          <p class="peer-waiting-sub">音视频连接后将自动显示画面</p>
        </div>
        <p v-if="mockMode" class="mock-label">{{ role === 'user' ? '学霸' : '家长' }}（模拟）</p>
        <p v-else-if="receiveOnly && peerOnline" class="mock-label">{{ role === 'user' ? '学霸' : '家长' }}</p>
      </div>
    </div>

    <div v-if="inRoom" class="footer">
      <Button
        round
        block
        :type="role === 'mentor' ? 'danger' : 'default'"
        @click="onLeave"
      >
        {{ role === 'mentor' ? '结束本节课' : '退出观看' }}
      </Button>
    </div>
  </div>
</template>

<style scoped>
.classroom {
  background: #0f172a;
  color: #fff;
  min-height: 100vh;
}

.timer-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.08);
}

.timer-bar strong {
  font-size: 22px;
  color: #fbbf24;
}

.timer-bar strong.waiting {
  font-size: 16px;
  color: #94a3b8;
}

.status-hint {
  margin-bottom: 0;
}

.loading {
  margin-top: 80px;
  color: #fff;
}

.error {
  color: #f87171;
  padding: 16px;
  text-align: center;
}

.video-grid {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.video-box {
  background: #1e293b;
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 4/3;
  position: relative;
}

.video-inner {
  width: 100%;
  height: 100%;
}

.video-box.remote .video-inner {
  background: linear-gradient(135deg, #312e81, #1e293b);
}

.peer-waiting {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  text-align: center;
  background: rgba(15, 23, 42, 0.72);
  pointer-events: none;
}

.peer-waiting p {
  margin: 0;
  font-size: 15px;
  color: #e2e8f0;
}

.peer-waiting-sub {
  font-size: 12px !important;
  color: #94a3b8 !important;
}

.mock-label {
  position: absolute;
  bottom: 8px;
  left: 12px;
  margin: 0;
  font-size: 12px;
  opacity: 0.8;
}

.footer {
  padding: 16px;
}
</style>
