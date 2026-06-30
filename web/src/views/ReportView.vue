<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showFailToast } from 'vant'
import { NavBar, Loading, Rate, Empty } from 'vant'
import type { GrowthReport } from '@/api/types'
import { fetchGrowthReport, getErrorMessage } from '@/api/client'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const report = ref<GrowthReport | null>(null)

async function load() {
  loading.value = true
  try {
    report.value = await fetchGrowthReport(route.params.orderId as string)
  } catch (e) {
    showFailToast(getErrorMessage(e))
    router.back()
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="成长报告" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>
    <Empty v-else-if="!report" description="报告不存在" />

    <div v-else class="content">
      <header>
        <h2>{{ report.child_name || '学员' }} 的成长报告</h2>
        <p>授课学霸：{{ report.mentor_name }}</p>
        <p class="date">{{ report.created_at?.slice(0, 10) }}</p>
      </header>

      <section class="scores">
        <div class="score-row">
          <span>口语表达</span>
          <Rate :model-value="report.speaking_score" readonly />
        </div>
        <div class="score-row">
          <span>自信心</span>
          <Rate :model-value="report.confidence_score" readonly />
        </div>
        <div class="score-row">
          <span>词汇运用</span>
          <Rate :model-value="report.vocabulary_score" readonly />
        </div>
      </section>

      <section class="comment">
        <h3>老师评语</h3>
        <p>{{ report.comment }}</p>
      </section>
    </div>
  </div>
</template>

<style scoped>
.center-loading {
  margin-top: 80px;
  text-align: center;
}

.content {
  padding: 16px;
}

header h2 {
  margin: 0 0 8px;
  font-size: 20px;
  color: var(--tm-primary);
}

header p {
  margin: 0 0 4px;
  color: var(--tm-muted);
  font-size: 14px;
}

.date {
  margin-top: 8px !important;
}

.scores {
  margin: 20px 0;
  padding: 16px;
  background: var(--tm-card);
  border-radius: 12px;
}

.score-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.comment h3 {
  margin: 0 0 8px;
  font-size: 16px;
}

.comment p {
  margin: 0;
  line-height: 1.6;
  color: #374151;
}
</style>
