<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { NavBar, Field, Button, CellGroup, Picker, Popup, Cell } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { getErrorMessage } from '@/api/client'

const router = useRouter()
const auth = useAuthStore()

const childName = ref(auth.user?.child_name || '')
const childAge = ref(String(auth.user?.child_age || 6))
const englishLevel = ref(auth.user?.english_level || 'beginner')
const saving = ref(false)

const showLevelPicker = ref(false)
const levelOptions = [
  { text: '启蒙', value: 'beginner' },
  { text: '进阶', value: 'intermediate' },
  { text: '高级', value: 'advanced' },
]

const levelLabel = ref(
  levelOptions.find((o) => o.value === englishLevel.value)?.text || '启蒙',
)

watch(
  () => auth.user,
  (u) => {
    if (!u) return
    childName.value = u.child_name
    childAge.value = String(u.child_age)
    englishLevel.value = u.english_level
    levelLabel.value =
      levelOptions.find((o) => o.value === u.english_level)?.text || '启蒙'
  },
)

function onLevelConfirm({ selectedOptions }: { selectedOptions: { text: string; value: string }[] }) {
  englishLevel.value = selectedOptions[0]?.value || 'beginner'
  levelLabel.value = selectedOptions[0]?.text || '启蒙'
  showLevelPicker.value = false
}

async function onSave() {
  saving.value = true
  try {
    await auth.saveProfile({
      child_name: childName.value,
      child_age: Number(childAge.value),
      english_level: englishLevel.value as 'beginner' | 'intermediate' | 'advanced',
    })
    showSuccessToast('保存成功')
    router.push('/mentors')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    saving.value = false
  }
}

function onLogout() {
  auth.logout()
  router.replace('/login')
}
</script>

<template>
  <div class="page">
    <NavBar title="孩子档案" left-arrow @click-left="router.back()" />

    <CellGroup inset title="基本信息">
      <Field v-model="childName" label="孩子姓名" placeholder="请输入" />
      <Field v-model="childAge" type="digit" label="年龄" placeholder="6-14" />
      <Field
        v-model="levelLabel"
        is-link
        readonly
        label="英语基础"
        @click="showLevelPicker = true"
      />
    </CellGroup>

    <CellGroup inset title="课时账户">
      <Field :model-value="String(auth.user?.available_lessons ?? 0)" readonly label="可用课时" />
      <Cell title="购买课时包" is-link @click="router.push('/packages')" />
      <Cell title="我的预约" is-link @click="router.push('/orders')" />
    </CellGroup>

    <div class="actions">
      <Button round block type="primary" :loading="saving" @click="onSave">保存</Button>
      <Button round block plain type="danger" class="logout" @click="onLogout">退出登录</Button>
    </div>

    <Popup v-model:show="showLevelPicker" position="bottom" round>
      <Picker
        :columns="levelOptions"
        @confirm="onLevelConfirm"
        @cancel="showLevelPicker = false"
      />
    </Popup>
  </div>
</template>

<style scoped>
.actions {
  padding: 24px 16px;
}

.logout {
  margin-top: 12px;
}
</style>
