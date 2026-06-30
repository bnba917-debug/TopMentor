<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { NavBar, Field, Button, CellGroup, Rate } from 'vant'
import { submitMentorReport, getErrorMessage } from '@/api/client'

const route = useRoute()
const router = useRouter()

const speaking = ref(4)
const confidence = ref(4)
const vocabulary = ref(4)
const comment = ref('')
const submitting = ref(false)

async function onSubmit() {
  if (comment.value.trim().length < 10) {
    showFailToast('评语至少 10 个字')
    return
  }
  submitting.value = true
  try {
    await submitMentorReport({
      order_id: route.params.orderId as string,
      speaking_score: speaking.value,
      confidence_score: confidence.value,
      vocabulary_score: vocabulary.value,
      comment: comment.value.trim(),
    })
    showSuccessToast('报告已提交')
    router.replace({ name: 'mentor-orders' })
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="page">
    <NavBar title="成长报告" left-arrow @click-left="router.back()" />

    <CellGroup inset title="评分（1-5 星）">
      <Field label="口语表达">
        <template #input>
          <Rate v-model="speaking" :count="5" />
        </template>
      </Field>
      <Field label="自信心">
        <template #input>
          <Rate v-model="confidence" :count="5" />
        </template>
      </Field>
      <Field label="词汇运用">
        <template #input>
          <Rate v-model="vocabulary" :count="5" />
        </template>
      </Field>
    </CellGroup>

    <CellGroup inset title="评语">
      <Field
        v-model="comment"
        rows="4"
        autosize
        type="textarea"
        maxlength="500"
        show-word-limit
        placeholder="请描述本节课表现与建议（至少 10 字）"
      />
    </CellGroup>

    <div class="actions">
      <Button round block type="primary" :loading="submitting" @click="onSubmit">提交报告</Button>
    </div>
  </div>
</template>

<style scoped>
.actions {
  padding: 24px 16px;
}
</style>
