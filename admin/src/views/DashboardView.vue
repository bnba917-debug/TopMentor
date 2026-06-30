<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { FinanceSummary } from '@/api/types'
import { fetchFinanceSummary, getErrorMessage } from '@/api/client'

const loading = ref(true)
const summary = ref<FinanceSummary | null>(null)

async function load() {
  loading.value = true
  try {
    summary.value = await fetchFinanceSummary()
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div v-loading="loading" class="page-card">
    <h2>财务大盘</h2>
    <el-row v-if="summary" :gutter="16" class="stats">
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">¥{{ summary.total_recharge_yuan.toFixed(2) }}</div>
          <div class="stat-label">累计充值</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">{{ summary.unspent_lessons }}</div>
          <div class="stat-label">未消耗课时（负债）</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">{{ summary.completed_orders }}</div>
          <div class="stat-label">已完成订单</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">¥{{ summary.total_mentor_earn_yuan.toFixed(2) }}</div>
          <div class="stat-label">学霸结算总额</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">¥{{ summary.total_withdraw_yuan.toFixed(2) }}</div>
          <div class="stat-label">提现总额</div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="stat-box">
          <div class="stat-value">{{ summary.pending_mentors }}</div>
          <div class="stat-label">待审核学霸</div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
h2 {
  margin: 0 0 20px;
}

.stats {
  margin-top: 8px;
}

.stat-box {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 16px;
}
</style>
