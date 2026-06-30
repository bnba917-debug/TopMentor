<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog, showFailToast, showSuccessToast } from 'vant'
import {
  NavBar,
  Field,
  Button,
  CellGroup,
  Picker,
  Popup,
  Tag,
  Image as VanImage,
  Loading,
  Progress,
} from 'vant'
import {
  fetchMentorProfile,
  updateMentorProfile,
  uploadMentorFile,
  mediaUrl,
  getErrorMessage,
} from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { MentorPortalProfile } from '@/api/types'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(true)
const saving = ref(false)
const uploadingAvatar = ref(false)
const uploadingVideo = ref(false)
const videoProgress = ref(0)

const realName = ref('')
const schoolName = ref('')
const major = ref('')
const englishScore = ref('')
const bio = ref('')
const tagsText = ref('')
const gender = ref('unknown')
const avatarUrl = ref('')
const introVideoUrl = ref('')

const showGenderPicker = ref(false)
const genderOptions = [
  { text: '男', value: 'male' },
  { text: '女', value: 'female' },
  { text: '不展示', value: 'unknown' },
]

const genderLabel = computed(
  () => genderOptions.find((o) => o.value === gender.value)?.text || '不展示',
)

const tagList = computed(() =>
  tagsText.value
    .split(/[,，]/)
    .map((t) => t.trim())
    .filter(Boolean)
    .slice(0, 6),
)

function fillForm(p: MentorPortalProfile) {
  realName.value = p.real_name
  schoolName.value = p.school_name
  major.value = p.major
  gender.value = p.gender || 'unknown'
  englishScore.value = p.english_score
  bio.value = p.bio
  tagsText.value = (p.tags || []).join('，')
  avatarUrl.value = p.avatar_url
  introVideoUrl.value = p.intro_video_url
}

function isAuthError(e: unknown): boolean {
  const msg = getErrorMessage(e)
  return msg.includes('学霸账号') || msg.includes('未登录') || msg.includes('Token')
}

async function promptRelogin(message: string) {
  try {
    await showConfirmDialog({
      title: '需要重新登录',
      message,
      confirmButtonText: '去登录',
    })
    auth.logout()
    router.replace({ name: 'login', query: { redirect: '/mentor/profile' } })
  } catch {
    /* cancelled */
  }
}

async function load() {
  loading.value = true
  try {
    const p = await fetchMentorProfile()
    fillForm(p)
    auth.syncMentorSummary(p)
  } catch (e) {
    if (isAuthError(e)) {
      await promptRelogin('当前登录状态无法编辑学霸资料，请重新登录后再试。')
    } else {
      showFailToast(getErrorMessage(e))
      router.back()
    }
  } finally {
    loading.value = false
  }
}

function onGenderConfirm({ selectedOptions }: { selectedOptions: { text: string; value: string }[] }) {
  gender.value = selectedOptions[0]?.value || 'unknown'
  showGenderPicker.value = false
}

async function onPickAvatar(ev: Event) {
  const input = ev.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  uploadingAvatar.value = true
  try {
    const res = await uploadMentorFile('avatar', file)
    avatarUrl.value = res.url
    showSuccessToast('头像已更新')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    uploadingAvatar.value = false
  }
}

async function onPickVideo(ev: Event) {
  const input = ev.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  if (file.size > 30 * 1024 * 1024) {
    showFailToast('视频请控制在 30MB 以内')
    return
  }
  uploadingVideo.value = true
  videoProgress.value = 0
  try {
    const res = await uploadMentorFile('intro_video', file, (pct) => {
      videoProgress.value = pct
    })
    introVideoUrl.value = res.url
    showSuccessToast('宣传视频已上传')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    uploadingVideo.value = false
    videoProgress.value = 0
  }
}

async function onSave() {
  if (!realName.value.trim() || !schoolName.value.trim() || !major.value.trim()) {
    showFailToast('请填写姓名、学校和专业')
    return
  }
  saving.value = true
  try {
    const updated = await updateMentorProfile({
      real_name: realName.value.trim(),
      school_name: schoolName.value.trim(),
      major: major.value.trim(),
      gender: gender.value,
      english_score: englishScore.value.trim(),
      bio: bio.value.trim(),
      tags: tagList.value,
      avatar_url: avatarUrl.value,
      intro_video_url: introVideoUrl.value,
    })
    fillForm(updated)
    auth.syncMentorSummary(updated)
    showSuccessToast('资料已保存')
    router.back()
  } catch (e) {
    const msg = getErrorMessage(e)
    if (isAuthError(e)) {
      await promptRelogin('登录状态已失效或尚未获得学霸权限，请重新登录。')
    } else {
      showFailToast(msg)
    }
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <NavBar title="编辑资料" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" vertical class="loading">加载中...</Loading>

    <template v-else>
      <div class="avatar-section">
        <VanImage
          round
          width="88"
          height="88"
          fit="cover"
          :src="mediaUrl(avatarUrl) || undefined"
        >
          <template #error>
            <div class="avatar-fallback">{{ realName.slice(0, 1) || '?' }}</div>
          </template>
        </VanImage>
        <label class="upload-btn">
          <input type="file" accept="image/jpeg,image/png,image/webp" hidden @change="onPickAvatar" />
          {{ uploadingAvatar ? '上传中…' : '更换头像' }}
        </label>
        <p class="hint">支持 JPG/PNG/WebP，最大 2MB</p>
      </div>

      <CellGroup inset title="基本信息">
        <Field v-model="realName" label="姓名" placeholder="真实姓名" maxlength="50" />
        <Field v-model="schoolName" label="学校" placeholder="如：清华大学" maxlength="100" />
        <Field v-model="major" label="专业" placeholder="如：计算机科学" maxlength="100" />
        <Field
          v-model="genderLabel"
          is-link
          readonly
          label="性别"
          @click="showGenderPicker = true"
        />
        <Field v-model="englishScore" label="英语成绩" placeholder="如：高考英语 148 分 / 托福 110" />
      </CellGroup>

      <CellGroup inset title="个人标签">
        <Field
          v-model="tagsText"
          rows="2"
          autosize
          type="textarea"
          label="标签"
          placeholder="用逗号分隔，如：阳光幽默，善于引导（最多 6 个）"
        />
        <div v-if="tagList.length" class="tag-preview">
          <Tag v-for="t in tagList" :key="t" plain type="primary">{{ t }}</Tag>
        </div>
      </CellGroup>

      <CellGroup inset title="个人简介">
        <Field
          v-model="bio"
          rows="3"
          autosize
          type="textarea"
          maxlength="500"
          show-word-limit
          placeholder="简单介绍你的教学风格、经历，家长预约前会看到"
        />
      </CellGroup>

      <CellGroup id="video" inset title="30 秒宣传视频">
        <div class="video-block">
          <video
            v-if="introVideoUrl"
            class="preview-video"
            :src="mediaUrl(introVideoUrl)"
            controls
            playsinline
          />
          <div v-else class="video-empty">
            <p>尚未上传宣传视频</p>
            <p class="sub">建议竖屏 30 秒内，介绍自己与教学特点</p>
          </div>

          <label class="upload-btn block">
            <input type="file" accept="video/mp4,video/webm,video/quicktime" hidden @change="onPickVideo" />
            {{ uploadingVideo ? `上传中 ${videoProgress}%` : introVideoUrl ? '重新上传视频' : '上传宣传视频' }}
          </label>
          <Progress v-if="uploadingVideo" :percentage="videoProgress" stroke-width="6" />
          <p class="hint">支持 MP4/WebM/MOV，最大 30MB</p>
        </div>
      </CellGroup>

      <div class="footer">
        <Button round block type="primary" :loading="saving" @click="onSave">保存资料</Button>
      </div>
    </template>

    <Popup v-model:show="showGenderPicker" position="bottom" round>
      <Picker :columns="genderOptions" @confirm="onGenderConfirm" @cancel="showGenderPicker = false" />
    </Popup>
  </div>
</template>

<style scoped>
.loading {
  margin-top: 80px;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 16px 8px;
  gap: 10px;
}

.avatar-fallback {
  width: 88px;
  height: 88px;
  border-radius: 50%;
  background: #dbeafe;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 600;
}

.upload-btn {
  display: inline-block;
  padding: 8px 16px;
  border-radius: 999px;
  background: #eff6ff;
  color: #2563eb;
  font-size: 14px;
  cursor: pointer;
}

.upload-btn.block {
  width: 100%;
  text-align: center;
  box-sizing: border-box;
}

.hint {
  margin: 0;
  font-size: 12px;
  color: #94a3b8;
  text-align: center;
}

.tag-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0 16px 12px;
}

.video-block {
  padding: 12px 16px 16px;
}

.preview-video {
  width: 100%;
  max-height: 220px;
  border-radius: 12px;
  background: #0f172a;
  margin-bottom: 12px;
}

.video-empty {
  text-align: center;
  padding: 28px 12px;
  background: #f1f5f9;
  border-radius: 12px;
  margin-bottom: 12px;
}

.video-empty p {
  margin: 0;
  color: #334155;
}

.video-empty .sub {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
}

.footer {
  padding: 16px 16px 32px;
}
</style>
