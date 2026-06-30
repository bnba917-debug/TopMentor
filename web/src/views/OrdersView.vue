<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Tag, Button, Empty, Loading } from 'vant'
import type { CourseOrderDetail } from '@/api/types'
import { fetchUserOrders, getErrorMessage } from '@/api/client'

const router = useRouter()
const loading = ref(true)
const orders = ref<CourseOrderDetail[]>([])

function statusLabel(s: string) {
  const map: Record<string, string> = {
    RESERVED: '已预约',
    ACTIVE: '上课中',
    COMPLETED: '已完成',
    CANCELLED: '已取消',
  }
  return map[s] || s
}

function canEnter(o: CourseOrderDetail) {
  return o.status === 'RESERVED' || o.status === 'ACTIVE'
}

function slotLabel(o: CourseOrderDetail) {
  if (!o.slot_date) return ''
  return `${o.slot_date} ${o.start_time?.slice(0, 5)}-${o.end_time?.slice(0, 5)}`
}

async function load() {
  loading.value = true
  try {
    orders.value = await fetchUserOrders()
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
    <NavBar title="我的预约" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>
    <Empty v-else-if="orders.length === 0" description="暂无预约" />

    <div v-else class="list">
      <article v-for="o in orders" :key="o.id" class="card">
        <div class="head">
          <h3>{{ o.mentor_name }}</h3>
          <Tag :type="o.status === 'ACTIVE' ? 'success' : 'primary'">{{ statusLabel(o.status) }}</Tag>
        </div>
        <p class="time">{{ slotLabel(o) }}</p>
        <p class="meta">订单 {{ o.id.slice(0, 12) }}…</p>
        <Button
          v-if="canEnter(o)"
          round
          block
          type="primary"
          size="small"
          @click="router.push({ name: 'classroom', params: { orderId: o.id } })"
        >
          进入课室
        </Button>
        <Button
          v-else-if="o.status === 'COMPLETED'"
          round
          block
          plain
          type="primary"
          size="small"
          @click="router.push({ name: 'report', params: { orderId: o.id } })"
        >
          查看成长报告
        </Button>
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
</style>
