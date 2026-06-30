<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showFailToast, showSuccessToast } from 'vant'
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
  Steps,
  Step,
} from 'vant'
import {
  fetchMentorApplyStatus,
  submitMentorApply,
  uploadApplyFile,
  mediaUrl,
  getErrorMessage,
} from '@/api/client'

const router = useRouter()
const loading = ref(true)
const submitting = ref(false)
const activeStep = ref(0)

const realName = ref('')
const schoolName = ref('')
const major = ref('')
const englishScore = ref('')
const bio = ref('')
const tagsText = ref('')
const gender = ref('unknown')
const avatarUrl = ref('')
const introVideoUrl = ref('')
const idCardUrl = ref('')
const studentCardUrl = ref('')
const englishProofUrl = ref('')

const uploading = ref('')
const uploadProgress = ref(0)

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

function fillDraft(p: NonNullable<Awaited<ReturnType<typeof fetchMentorApplyStatus>>['profile']>) {
  if (!p) return
  realName.value = p.real_name || ''
  schoolName.value = p.school_name || ''
  major.value = p.major || ''
  gender.value = p.gender || 'unknown'
  englishScore.value = p.english_score || ''
  bio.value = p.bio || ''
  tagsText.value = (p.tags || []).join('，')
  avatarUrl.value = p.avatar_url || ''
  introVideoUrl.value = p.intro_video_url || ''
  idCardUrl.value = p.id_card_url || ''
  studentCardUrl.value = p.student_card_url || ''
  englishProofUrl.value = p.english_proof_url || ''
}

async function load() {
  loading.value = true
  try {
    const st = await fetchMentorApplyStatus()
    if (st.status === 'pending') {
      router.replace({ name: 'mentor-apply-status' })
      return
    }
    if (st.status === 'approved') {
      router.replace('/mentor')
      return
    }
    if (st.profile) fillDraft(st.profile)
  } catch (e) {
    showFailToast(getErrorMessage(e))
    router.back()
  } finally {
    loading.value = false
  }
}

function onGenderConfirm({ selectedOptions }: { selectedOptions: { text: string; value: string }[] }) {
  gender.value = selectedOptions[0]?.value || 'unknown'
  showGenderPicker.value = false
}

async function onPickFile(
  kind: 'avatar' | 'intro_video' | 'id_card' | 'student_card' | 'english_proof',
  ev: Event,
) {
  const input = ev.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  uploading.value = kind
  uploadProgress.value = 0
  try {
    const res = await uploadApplyFile(kind, file, (pct) => {
      uploadProgress.value = pct
    })
    if (kind === 'avatar') avatarUrl.value = res.url
    if (kind === 'intro_video') introVideoUrl.value = res.url
    if (kind === 'id_card') idCardUrl.value = res.url
    if (kind === 'student_card') studentCardUrl.value = res.url
    if (kind === 'english_proof') englishProofUrl.value = res.url
    showSuccessToast('上传成功')
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    uploading.value = ''
    uploadProgress.value = 0
  }
}

function validateStep(step: number): boolean {
  if (step === 0) {
    if (!realName.value.trim() || !schoolName.value.trim() || !major.value.trim()) {
      showFailToast('请填写姓名、学校和专业')
      return false
    }
  }
  if (step === 1) {
    if (!idCardUrl.value || !studentCardUrl.value) {
      showFailToast('请上传身份证和学生证照片')
      return false
    }
  }
  if (step === 2) {
    if (!avatarUrl.value) {
      showFailToast('请上传头像')
      return false
    }
    if (!introVideoUrl.value) {
      showFailToast('请上传 30 秒宣传视频')
      return false
    }
  }
  return true
}

function nextStep() {
  if (!validateStep(activeStep.value)) return
  if (activeStep.value < 3) activeStep.value += 1
}

function prevStep() {
  if (activeStep.value > 0) activeStep.value -= 1
}

async function onSubmit() {
  if (!validateStep(0) || !validateStep(1) || !validateStep(2)) {
    showFailToast('请完善所有必填项')
    return
  }
  submitting.value = true
  try {
    await submitMentorApply({
      real_name: realName.value.trim(),
      school_name: schoolName.value.trim(),
      major: major.value.trim(),
      gender: gender.value,
      english_score: englishScore.value.trim(),
      bio: bio.value.trim(),
      tags: tagList.value,
      avatar_url: avatarUrl.value,
      intro_video_url: introVideoUrl.value,
      id_card_url: idCardUrl.value,
      student_card_url: studentCardUrl.value,
      english_proof_url: englishProofUrl.value || undefined,
    })
    showSuccessToast('申请已提交')
    router.replace({ name: 'mentor-apply-status' })
  } catch (e) {
    showFailToast(getErrorMessage(e))
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page apply-page">
    <NavBar title="学霸入驻申请" left-arrow @click-left="router.back()" />

    <Loading v-if="loading" vertical class="loading">加载中...</Loading>

    <template v-else>
      <Steps :active="activeStep" active-color="#2563eb" class="steps">
        <Step>基本资料</Step>
        <Step>证件上传</Step>
        <Step>形象视频</Step>
        <Step>确认提交</Step>
      </Steps>

      <div v-show="activeStep === 0" class="step-body">
        <CellGroup inset title="基本信息">
          <Field v-model="realName" label="姓名" placeholder="真实姓名" maxlength="50" />
          <Field v-model="schoolName" label="学校" placeholder="如：清华大学" maxlength="100" />
          <Field v-model="major" label="专业" placeholder="如：计算机科学" maxlength="100" />
          <Field v-model="genderLabel" is-link readonly label="性别" @click="showGenderPicker = true" />
          <Field v-model="englishScore" label="英语成绩" placeholder="高考/四六级/托福雅思等" />
        </CellGroup>
        <CellGroup inset title="个人标签">
          <Field
            v-model="tagsText"
            rows="2"
            autosize
            type="textarea"
            label="标签"
            placeholder="逗号分隔，如：阳光幽默，善于引导"
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
            placeholder="介绍你的教学风格与经历"
          />
        </CellGroup>
      </div>

      <div v-show="activeStep === 1" class="step-body">
        <CellGroup inset title="身份与学籍证明">
          <div class="upload-card">
            <p class="upload-title">身份证正面 / 反面（可合并一张）</p>
            <img v-if="idCardUrl" :src="mediaUrl(idCardUrl)" alt="" class="preview-img" />
            <label class="upload-btn">
              <input type="file" accept="image/*,.pdf" hidden @change="onPickFile('id_card', $event)" />
              {{ uploading === 'id_card' ? '上传中…' : idCardUrl ? '重新上传' : '上传身份证' }}
            </label>
          </div>
          <div class="upload-card">
            <p class="upload-title">学生证或在读证明</p>
            <img v-if="studentCardUrl" :src="mediaUrl(studentCardUrl)" alt="" class="preview-img" />
            <label class="upload-btn">
              <input type="file" accept="image/*,.pdf" hidden @change="onPickFile('student_card', $event)" />
              {{ uploading === 'student_card' ? '上传中…' : studentCardUrl ? '重新上传' : '上传学生证' }}
            </label>
          </div>
          <div class="upload-card">
            <p class="upload-title">英语成绩证明（选填）</p>
            <img v-if="englishProofUrl" :src="mediaUrl(englishProofUrl)" alt="" class="preview-img" />
            <label class="upload-btn">
              <input type="file" accept="image/*,.pdf" hidden @change="onPickFile('english_proof', $event)" />
              {{ uploading === 'english_proof' ? '上传中…' : englishProofUrl ? '重新上传' : '上传成绩证明' }}
            </label>
          </div>
          <p class="hint">支持 JPG/PNG/PDF，单张最大 5MB</p>
        </CellGroup>
      </div>

      <div v-show="activeStep === 2" class="step-body">
        <div class="avatar-section">
          <VanImage round width="88" height="88" fit="cover" :src="mediaUrl(avatarUrl) || undefined">
            <template #error>
              <div class="avatar-fallback">{{ realName.slice(0, 1) || '?' }}</div>
            </template>
          </VanImage>
          <label class="upload-btn">
            <input type="file" accept="image/jpeg,image/png,image/webp" hidden @change="onPickFile('avatar', $event)" />
            {{ uploading === 'avatar' ? '上传中…' : '上传头像' }}
          </label>
        </div>
        <CellGroup inset title="30 秒宣传视频">
          <div class="video-block">
            <video v-if="introVideoUrl" class="preview-video" :src="mediaUrl(introVideoUrl)" controls playsinline />
            <div v-else class="video-empty">
              <p>录制或上传约 30 秒自荐视频</p>
              <p class="sub">介绍学校背景、英语能力与陪伴风格</p>
            </div>
            <label class="upload-btn block">
              <input type="file" accept="video/mp4,video/webm,video/quicktime" hidden @change="onPickFile('intro_video', $event)" />
              {{ uploading === 'intro_video' ? `上传中 ${uploadProgress}%` : introVideoUrl ? '重新上传视频' : '上传宣传视频' }}
            </label>
            <Progress v-if="uploading === 'intro_video'" :percentage="uploadProgress" stroke-width="6" />
          </div>
        </CellGroup>
      </div>

      <div v-show="activeStep === 3" class="step-body">
        <CellGroup inset title="请确认提交内容">
          <Field :model-value="realName" readonly label="姓名" />
          <Field :model-value="`${schoolName} · ${major}`" readonly label="学校专业" />
          <Field :model-value="englishScore" readonly label="英语成绩" />
          <Field :model-value="tagList.join('、')" readonly label="标签" />
          <Field :model-value="bio" readonly type="textarea" rows="2" autosize label="简介" />
        </CellGroup>
        <p class="confirm-hint">提交后进入人工审核，审核通过即可接单授课。</p>
      </div>

      <Progress v-if="uploading && activeStep !== 2" :percentage="uploadProgress" stroke-width="6" class="progress" />

      <div class="footer">
        <Button v-if="activeStep > 0" round block plain @click="prevStep">上一步</Button>
        <Button v-if="activeStep < 3" round block type="primary" @click="nextStep">下一步</Button>
        <Button v-else round block type="primary" :loading="submitting" @click="onSubmit">提交申请</Button>
      </div>
    </template>

    <Popup v-model:show="showGenderPicker" position="bottom" round>
      <Picker :columns="genderOptions" @confirm="onGenderConfirm" @cancel="showGenderPicker = false" />
    </Popup>
  </div>
</template>

<style scoped>
.apply-page {
  padding-bottom: 24px;
}

.loading {
  margin-top: 80px;
}

.steps {
  margin: 12px 0 8px;
}

.step-body {
  padding-bottom: 8px;
}

.tag-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 0 16px 12px;
}

.upload-card {
  padding: 12px 16px;
  border-bottom: 1px solid #f1f5f9;
}

.upload-card:last-child {
  border-bottom: none;
}

.upload-title {
  margin: 0 0 10px;
  font-size: 14px;
  color: #334155;
}

.preview-img {
  display: block;
  width: 100%;
  max-height: 160px;
  object-fit: contain;
  border-radius: 8px;
  margin-bottom: 10px;
  background: #f8fafc;
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
  display: block;
  text-align: center;
  width: 100%;
  box-sizing: border-box;
}

.hint {
  padding: 8px 16px 0;
  margin: 0;
  font-size: 12px;
  color: #94a3b8;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 16px;
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
  padding: 24px 12px;
  background: #f1f5f9;
  border-radius: 12px;
  margin-bottom: 12px;
}

.video-empty p {
  margin: 0;
}

.video-empty .sub {
  margin-top: 6px;
  font-size: 12px;
  color: #94a3b8;
}

.confirm-hint {
  padding: 12px 16px;
  margin: 0;
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
}

.progress {
  margin: 0 16px;
}

.footer {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
