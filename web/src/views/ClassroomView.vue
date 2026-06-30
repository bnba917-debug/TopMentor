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

    <div class="classroom-inner">
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
        <div class="video-box local">
          <div ref="localVideoRef" class="video-inner" />
          <span class="video-badge">我 · {{ role === 'user' ? '家长' : '学霸' }}</span>
          <p v-if="receiveOnly" class="video-note">本机无摄像头</p>
        </div>

        <div class="video-box remote">
          <div ref="remoteVideoRef" class="video-inner" />
          <span class="video-badge primary">{{ role === 'user' ? '学霸' : '家长' }}</span>
          <div v-if="!peerOnline && lessonStarted && !mockMode" class="peer-waiting">
            <p>{{ peerLabel }}尚未进房</p>
            <p class="peer-waiting-sub">音视频连接后将自动显示画面</p>
          </div>
          <div v-else-if="mockMode" class="peer-waiting">
            <p>模拟模式：无真实音视频</p>
            <p class="peer-waiting-sub">请用两个浏览器分别登录家长端与学霸端，并在服务器配置 Agora 密钥后测试</p>
          </div>
        </div>
      </div>

      <div v-if="inRoom" class="footer">
        <Button
          round
          block
          :type="role === 'mentor' ? 'danger' : 'default'"
          class="leave-btn"
          @click="onLeave"
        >
          {{ role === 'mentor' ? '结束本节课' : '退出观看' }}
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.classroom {
  background: #0f172a;
  color: #fff;
  min-height: 100vh;
}

.classroom-inner {
  width: 100%;
  max-width: 960px;
  margin: 0 auto;
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
  padding: 12px 16px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.video-box {
  background: #1e293b;
  border-radius: 12px;
  overflow: hidden;
  aspect-ratio: 16 / 9;
  max-height: 42vh;
  position: relative;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.video-box.remote {
  aspect-ratio: 16 / 9;
  max-height: 48vh;
}

.video-inner {
  width: 100%;
  height: 100%;
}

.video-inner :deep(video) {
  object-fit: cover;
}

.video-box.remote .video-inner {
  background: linear-gradient(135deg, #312e81, #1e293b);
}

.video-badge {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 2;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 12px;
  background: rgba(15, 23, 42, 0.72);
  backdrop-filter: blur(4px);
}

.video-badge.primary {
  background: rgba(79, 70, 229, 0.85);
}

.video-note {
  position: absolute;
  bottom: 8px;
  left: 10px;
  margin: 0;
  font-size: 11px;
  color: #cbd5e1;
  z-index: 2;
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
  z-index: 1;
}

.peer-waiting p {
  margin: 0;
  font-size: 15px;
  color: #e2e8f0;
}

.peer-waiting-sub {
  font-size: 12px !important;
  color: #94a3b8 !important;
  max-width: 320px;
  line-height: 1.5;
}

.footer {
  padding: 0 16px 20px;
}

.leave-btn {
  max-width: 360px;
  margin: 0 auto;
}

/* 平板 / 电脑：对方大画面 + 自己小窗 */
@media (min-width: 768px) {
  .classroom-inner {
    padding: 0 20px 24px;
  }

  .timer-bar {
    border-radius: 12px;
    margin-top: 12px;
  }

  .video-grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 240px;
    grid-template-rows: auto;
    gap: 14px;
    align-items: stretch;
    padding-top: 16px;
  }

  .video-box.local {
    grid-column: 2;
    grid-row: 1;
    max-height: 135px;
    aspect-ratio: 4 / 3;
    align-self: start;
  }

  .video-box.remote {
    grid-column: 1;
    grid-row: 1;
    max-height: min(52vh, 420px);
    aspect-ratio: auto;
    min-height: 280px;
  }

  .peer-waiting-sub {
    max-width: 420px;
  }

  .footer {
    padding-top: 4px;
  }
}

@media (min-width: 1024px) {
  .video-grid {
    grid-template-columns: minmax(0, 1fr) 280px;
    gap: 16px;
  }

  .video-box.local {
    max-height: 158px;
  }

  .video-box.remote {
    max-height: min(50vh, 460px);
  }
}
</style>
