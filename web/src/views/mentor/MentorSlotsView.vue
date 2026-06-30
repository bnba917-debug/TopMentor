<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { NavBar, Button, Loading, Tag, Empty } from 'vant'
import type { MentorSlot } from '@/api/types'
import { fetchMentorPortalSlots, setMentorSlots, getErrorMessage } from '@/api/client'

const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const slots = ref<MentorSlot[]>([])

const presetTimes = ['09:00:00', '10:00:00', '14:00:00', '15:00:00', '16:00:00', '19:00:00', '20:00:00']

const tomorrow = computed(() => {
  const d = new Date()
  d.setDate(d.getDate() + 1)
  return d.toISOString().slice(0, 10)
})

const dayAfter = computed(() => {
  const d = new Date()
  d.setDate(d.getDate() + 2)
  return d.toISOString().slice(0, 10)
})

function slotKey(date: string, time: string) {
  return `${date}_${time}`
}

const selected = ref<Set<string>>(new Set())

function isSelected(date: string, time: string) {
  return selected.value.has(slotKey(date, time))
}

function toggle(date: string, time: string) {
  const key = slotKey(date, time)
  const next = new Set(selected.value)
  if (next.has(key)) {
    next.delete(key)
  } else {
    next.add(key)
  }
  selected.value = next
}

function applyExisting(list: MentorSlot[]) {
  const next = new Set<string>()
  for (const s of list) {
    if (s.status === 0) {
      next.add(slotKey(s.slot_date, s.start_time))
    }
  }
  selected.value = next
}

async function load() {
  loading.value = true
  try {
    const list = await fetchMentorPortalSlots({ from: tomorrow.value, to: dayAfter.value })
    slots.value = list
    applyExisting(list)
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    loading.value = false
  }
}

async function onSave() {
  saving.value = true
  try {
    const dates = [tomorrow.value, dayAfter.value]
    const payload = dates.flatMap((date) =>
      presetTimes.map((time) => ({
        slot_date: date,
        start_time: time,
        available: isSelected(date, time),
      })),
    )
    await setMentorSlots(payload)
    showSuccessToast('时段已更新')
    await load()
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="开放时段" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" class="center-loading" vertical>加载中...</Loading>

    <template v-else>
      <p class="hint">点击时段切换开放/关闭（未来 2 天）</p>

      <section v-for="date in [tomorrow, dayAfter]" :key="date" class="day-block">
        <h3>{{ date }}</h3>
        <div class="grid">
          <Tag
            v-for="time in presetTimes"
            :key="slotKey(date, time)"
            size="large"
            :type="isSelected(date, time) ? 'primary' : 'default'"
            @click="toggle(date, time)"
          >
            {{ time.slice(0, 5) }}
          </Tag>
        </div>
      </section>

      <Empty v-if="presetTimes.length === 0" description="暂无预设时段" />

      <div class="actions">
        <Button round block type="primary" :loading="saving" @click="onSave">保存时段</Button>
      </div>
    </template>
  </div>
</template>

<style scoped>
.center-loading {
  margin-top: 80px;
  text-align: center;
}

.hint {
  margin: 12px 16px;
  font-size: 13px;
  color: var(--tm-muted);
}

.day-block {
  padding: 0 16px 16px;
}

.day-block h3 {
  margin: 0 0 10px;
  font-size: 15px;
}

.grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.actions {
  padding: 8px 16px 24px;
}
</style>
