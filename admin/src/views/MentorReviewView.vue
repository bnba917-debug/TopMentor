<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { PendingMentorApplication } from '@/api/types'
import { fetchPendingMentors, reviewMentor, getErrorMessage, maskUrl } from '@/api/client'

const loading = ref(true)
const list = ref<PendingMentorApplication[]>([])

async function load() {
  loading.value = true
  try {
    list.value = await fetchPendingMentors()
  } catch (e) {
    ElMessage.error(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

async function onApprove(row: PendingMentorApplication) {
  try {
    await ElMessageBox.confirm(`确认通过 ${row.real_name} 的入驻申请？`, '审核通过')
    await reviewMentor(row.mentor_id, { action: 'approve' })
    ElMessage.success('已通过')
    await load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(getErrorMessage(e))
  }
}

async function onReject(row: PendingMentorApplication) {
  try {
    const { value } = await ElMessageBox.prompt('请填写驳回原因', '驳回申请', {
      inputPlaceholder: '材料不完整 / 视频不符合要求…',
      inputValidator: (v) => !!v?.trim() || '请填写原因',
    })
    await reviewMentor(row.mentor_id, { action: 'reject', reject_reason: value })
    ElMessage.success('已驳回')
    await load()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(getErrorMessage(e))
  }
}

onMounted(load)
</script>

<template>
  <div class="page-card">
    <div class="head">
      <h2>学霸资质审核</h2>
      <el-button @click="load">刷新</el-button>
    </div>

    <el-table v-loading="loading" :data="list" stripe>
      <el-table-column prop="real_name" label="姓名" width="100" />
      <el-table-column prop="school_name" label="学校" width="120" />
      <el-table-column prop="major" label="专业" width="120" />
      <el-table-column prop="english_score" label="英语成绩" min-width="140" />
      <el-table-column label="身份证（脱敏）" min-width="180">
        <template #default="{ row }">{{ maskUrl(row.id_card_url) }}</template>
      </el-table-column>
      <el-table-column label="学生证（脱敏）" min-width="180">
        <template #default="{ row }">{{ maskUrl(row.student_card_url) }}</template>
      </el-table-column>
      <el-table-column prop="applied_at" label="申请时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="onApprove(row)">通过</el-button>
          <el-button type="danger" link @click="onReject(row)">驳回</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && list.length === 0" description="暂无待审核申请" />
  </div>
</template>

<style scoped>
.head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

h2 {
  margin: 0;
}
</style>
