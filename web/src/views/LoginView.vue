<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
import { NavBar, Field, Button, CellGroup, Tag } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { getErrorMessage } from '@/api/client'
import { savePortalRole, type PortalRole } from '@/utils/portalRole'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const portalRole = computed(() => {
  const q = route.query.role as PortalRole | undefined
  if (q === 'parent' || q === 'mentor') return q
  return 'parent'
})

const roleLabel = computed(() => (portalRole.value === 'mentor' ? '学霸端' : '家长 / 学生端'))
const roleHint = computed(() =>
  portalRole.value === 'mentor'
    ? '登录已认证学霸进入工作台；新用户可提交入驻申请'
    : '登录后预约课程、进入课室观看',
)

const phone = ref('')
const code = ref('')
const sending = ref(false)
const logging = ref(false)
const countdown = ref(0)
const debugCode = ref('')

let timer: ReturnType<typeof setInterval> | null = null

function startCountdown() {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value -= 1
    if (countdown.value <= 0 && timer) {
      clearInterval(timer)
      timer = null
    }
  }, 1000)
}

async function onSendCode() {
  if (!/^1\d{10}$/.test(phone.value)) {
    showFailToast('请输入正确手机号')
    return
  }
  sending.value = true
  try {
    const res = await auth.requestCode(phone.value)
    showSuccessToast('验证码已发送')
    if (res.debug_code) {
      debugCode.value = res.debug_code
      code.value = res.debug_code
    }
    startCountdown()
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    sending.value = false
  }
}

async function onLogin() {
  if (!phone.value || !code.value) {
    showFailToast('请填写手机号和验证码')
    return
  }
  logging.value = true
  savePortalRole(portalRole.value)
  try {
    const result = await auth.login(phone.value, code.value)
    showSuccessToast('登录成功')
    const redirect = route.query.redirect as string | undefined

    if (result.mentor) {
      router.replace(redirect?.startsWith('/mentor') ? redirect : '/mentor')
      return
    }

    if (portalRole.value === 'mentor') {
      router.replace(redirect?.startsWith('/mentor') ? redirect : '/mentor/apply/status')
      return
    }

    router.replace(redirect || '/mentors')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    logging.value = false
  }
}

function switchRole() {
  router.replace({ name: 'welcome' })
}
</script>

<template>
  <div class="page login-page">
    <NavBar title="登录" left-arrow @click-left="switchRole" />

    <div class="hero">
      <Tag type="primary" plain>{{ roleLabel }}</Tag>
      <h1>手机号登录</h1>
      <p>{{ roleHint }}</p>
      <button type="button" class="switch-role" @click="switchRole">切换身份</button>
    </div>

    <CellGroup inset>
      <Field v-model="phone" type="tel" maxlength="11" label="手机号" placeholder="请输入手机号" />
      <Field v-model="code" maxlength="6" label="验证码" placeholder="6 位验证码">
        <template #button>
          <Button
            size="small"
            type="primary"
            plain
            :disabled="countdown > 0 || sending"
            :loading="sending"
            @click="onSendCode"
          >
            {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
          </Button>
        </template>
      </Field>
    </CellGroup>

    <p v-if="debugCode" class="hint">开发模式验证码：{{ debugCode }}</p>

    <div class="actions">
      <Button round block type="primary" :loading="logging" @click="onLogin">登录</Button>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  background: linear-gradient(180deg, #eef2ff 0%, var(--tm-bg) 40%);
}

.hero {
  padding: 24px 24px 20px;
}

.hero h1 {
  margin: 12px 0 8px;
  font-size: 24px;
  color: var(--tm-primary);
}

.hero p {
  margin: 0;
  color: var(--tm-muted);
  font-size: 14px;
  line-height: 1.5;
}

.switch-role {
  margin-top: 12px;
  padding: 0;
  border: none;
  background: none;
  color: var(--tm-primary);
  font-size: 14px;
  cursor: pointer;
}

.hint {
  margin: 12px 16px 0;
  font-size: 12px;
  color: var(--tm-muted);
}

.actions {
  padding: 24px 16px;
}
</style>
