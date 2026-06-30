<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Button, Empty, Loading, NoticeBar, Tag } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { fetchMentorApplyStatus, getErrorMessage } from '@/api/client'
import type { MentorApplyStatus } from '@/api/types'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const status = ref<MentorApplyStatus | null>(null)

async function load() {
  loading.value = true
  try {
    status.value = await fetchMentorApplyStatus()
    if (status.value.status === 'approved' && auth.isMentor) {
      router.replace('/mentor')
    }
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

function goApply() {
  router.push({ name: 'mentor-apply' })
}

function relogin() {
  auth.logout()
  router.replace({ name: 'login', query: { redirect: '/mentor' } })
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="入驻审核" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" vertical class="loading">加载中...</Loading>

    <template v-else-if="status">
      <div v-if="status.status === 'none'" class="panel">
        <Empty description="您还未提交学霸入驻申请">
          <Button round type="primary" @click="goApply">立即申请</Button>
        </Empty>
      </div>

      <div v-else-if="status.status === 'pending'" class="panel">
        <div class="icon pending">⏳</div>
        <h2>审核中</h2>
        <p>我们已收到您的入驻资料，通常 1–3 个工作日内完成审核。</p>
        <NoticeBar
          wrapable
          :scrollable="false"
          text="审核期间请保持手机畅通，如有疑问请联系平台客服。"
        />
        <p v-if="status.applied_at" class="meta">提交时间：{{ new Date(status.applied_at).toLocaleString() }}</p>
      </div>

      <div v-else-if="status.status === 'rejected'" class="panel">
        <div class="icon rejected">✕</div>
        <h2>审核未通过</h2>
        <NoticeBar
          v-if="status.reject_reason"
          wrapable
          :scrollable="false"
          :text="`原因：${status.reject_reason}`"
        />
        <p class="hint">您可以根据反馈修改资料后重新提交。</p>
        <Button round block type="primary" @click="goApply">修改并重新申请</Button>
      </div>

      <div v-else-if="status.status === 'approved'" class="panel">
        <div class="icon approved">✓</div>
        <h2>审核已通过</h2>
        <Tag type="success">认证学霸</Tag>
        <p class="hint">请重新登录以进入学霸工作台。</p>
        <Button round block type="primary" @click="relogin">重新登录</Button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.loading {
  margin-top: 80px;
}

.panel {
  padding: 32px 20px;
  text-align: center;
}

.icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.icon.pending {
  background: #fef3c7;
}

.icon.rejected {
  background: #fee2e2;
  color: #dc2626;
}

.icon.approved {
  background: #dcfce7;
  color: #16a34a;
}

.panel h2 {
  margin: 0 0 12px;
  font-size: 22px;
}

.panel p {
  margin: 0 0 16px;
  color: var(--tm-muted, #64748b);
  line-height: 1.6;
}

.hint {
  margin-top: 16px !important;
}

.meta {
  font-size: 12px;
  margin-top: 20px !important;
}
</style>
