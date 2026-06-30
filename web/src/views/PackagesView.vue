<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { NavBar, Button, Tag, Loading, Cell, CellGroup } from 'vant'
import type { LessonPackage, PaymentChannels } from '@/api/types'
import {
  createRecharge,
  fetchPackages,
  fetchPaymentChannels,
  formatYuan,
  getErrorMessage,
} from '@/api/client'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const paying = ref(false)
const packages = ref<LessonPackage[]>([])
const channels = ref<PaymentChannels | null>(null)

async function load() {
  loading.value = true
  try {
    const [pkgList, ch] = await Promise.all([fetchPackages(), fetchPaymentChannels()])
    packages.value = pkgList
    channels.value = ch
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

function channelLabel(ch: string) {
  const map: Record<string, string> = {
    mock: '模拟支付（开发）',
    wechat_jsapi: '微信支付（微信内）',
    wechat_h5: '微信支付（浏览器）',
    alipay_h5: '支付宝',
  }
  return map[ch] || ch
}

function pickChannel(): string {
  if (!channels.value) return 'mock'
  if (channels.value.mode === 'mock') return 'mock'
  const ua = navigator.userAgent.toLowerCase()
  if (ua.includes('micromessenger') && channels.value.channels.includes('wechat_jsapi')) {
    return 'wechat_jsapi'
  }
  if (channels.value.channels.includes('wechat_h5')) return 'wechat_h5'
  if (channels.value.channels.includes('alipay_h5')) return 'alipay_h5'
  return 'mock'
}

async function onBuy(pkg: LessonPackage) {
  if (!auth.isLoggedIn) {
    router.push({ name: 'login', query: { redirect: '/packages' } })
    return
  }
  paying.value = true
  try {
    const channel = pickChannel()
    const result = await createRecharge(pkg.id, channel)

    if (result.mock_paid) {
      showSuccessToast(`充值成功，当前 ${result.available_lessons} 课时`)
      if (auth.user && result.available_lessons != null) {
        auth.setAvailableLessons(result.available_lessons)
      }
      return
    }

    if (result.jsapi_params) {
      showFailToast('微信 JSAPI 支付待商户号接入后启用')
      return
    }
    if (result.pay_url) {
      window.location.href = result.pay_url
      return
    }
    showSuccessToast('订单已创建')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    paying.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="购买课时" left-arrow @click-left="router.back()" />

    <div v-if="channels" class="mode-banner">
      支付模式：<Tag type="primary">{{ channels.mode === 'mock' ? '模拟（立即到账）' : '正式' }}</Tag>
    </div>

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>

    <div v-else class="list">
      <div v-for="pkg in packages" :key="pkg.id" class="card">
        <div class="card-head">
          <h3>{{ pkg.name }}</h3>
          <Tag v-if="pkg.is_trial" type="danger">体验</Tag>
        </div>
        <p class="lessons">{{ pkg.lesson_count }} 课时 · ¥{{ formatYuan(pkg.price_cents) }}</p>
        <Button
          round
          block
          type="primary"
          size="small"
          :loading="paying"
          @click="onBuy(pkg)"
        >
          {{ channels?.mode === 'mock' ? '模拟购买' : '立即购买' }}
        </Button>
      </div>
    </div>

    <CellGroup v-if="channels && !loading" inset title="可用支付方式" class="channels">
      <Cell v-for="ch in channels.channels" :key="ch" :title="channelLabel(ch)" />
    </CellGroup>
  </div>
</template>

<style scoped>
.mode-banner {
  padding: 12px 16px;
  font-size: 13px;
  color: var(--tm-muted);
}

.center-loading {
  margin-top: 80px;
  text-align: center;
}

.list {
  padding: 0 16px 16px;
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
  align-items: center;
  gap: 8px;
}

.card-head h3 {
  margin: 0;
  font-size: 18px;
}

.lessons {
  margin: 12px 0 16px;
  font-size: 15px;
  color: var(--tm-primary);
  font-weight: 600;
}

.channels {
  margin-top: 8px;
}
</style>
