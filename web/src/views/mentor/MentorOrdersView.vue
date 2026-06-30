<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Tag, Button, Empty, Loading } from 'vant'
import type { MentorOrderDetail } from '@/api/types'
import { fetchMentorOrders, getErrorMessage } from '@/api/client'

const router = useRouter()
const loading = ref(true)
const orders = ref<MentorOrderDetail[]>([])

function statusLabel(s: string) {
  const map: Record<string, string> = {
    RESERVED: '已预约',
    ACTIVE: '上课中',
    COMPLETED: '已完成',
    CANCELLED: '已取消',
  }
  return map[s] || s
}

function slotLabel(o: MentorOrderDetail) {
  if (!o.slot_date) return ''
  return `${o.slot_date} ${o.start_time?.slice(0, 5)}-${o.end_time?.slice(0, 5)}`
}

function canEnter(o: MentorOrderDetail) {
  return o.status === 'RESERVED' || o.status === 'ACTIVE'
}

function canReport(o: MentorOrderDetail) {
  return o.status === 'COMPLETED' && !o.has_report
}

async function load() {
  loading.value = true
  try {
    orders.value = await fetchMentorOrders()
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="我的订单" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>
    <Empty v-else-if="orders.length === 0" description="暂无订单" />

    <div v-else class="list">
      <article v-for="o in orders" :key="o.id" class="card">
        <div class="head">
          <h3>{{ o.child_name || '学员' }}</h3>
          <Tag :type="o.status === 'ACTIVE' ? 'success' : 'primary'">{{ statusLabel(o.status) }}</Tag>
        </div>
        <p class="time">{{ slotLabel(o) }}</p>
        <p class="meta">{{ o.user_phone }}</p>
        <div class="actions">
          <Button
            v-if="canEnter(o)"
            round
            block
            type="primary"
            size="small"
            @click="router.push({ name: 'classroom', params: { orderId: o.id }, query: { role: 'mentor' } })"
          >
            进入课室
          </Button>
          <Button
            v-else-if="canReport(o)"
            round
            block
            type="primary"
            size="small"
            @click="router.push({ name: 'mentor-report', params: { orderId: o.id } })"
          >
            填写成长报告
          </Button>
          <Tag v-else-if="o.has_report" plain type="success">已提交报告</Tag>
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.center-loading {
  margin-top: 80px;
  text-align: center;
}

.list {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card {
  background: var(--tm-card);
  border-radius: 12px;
  padding: 16px;
}

.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.head h3 {
  margin: 0;
  font-size: 17px;
}

.time {
  margin: 10px 0 4px;
  color: var(--tm-primary);
  font-weight: 600;
}

.meta {
  margin: 0 0 12px;
  font-size: 12px;
  color: var(--tm-muted);
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
