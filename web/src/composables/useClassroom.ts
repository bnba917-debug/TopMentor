import { computed, onUnmounted, ref } from 'vue'
import AgoraRTC, {
  type IAgoraRTCClient,
  type IAgoraRTCRemoteUser,
  type ICameraVideoTrack,
  type IMicrophoneAudioTrack,
} from 'agora-rtc-sdk-ng'
import { completeRoom, joinRoom, sendHeartbeat, getErrorMessage } from '@/api/client'
import type { HeartbeatResult, JoinRoomResult } from '@/api/types'

function canUseLocalMedia(): boolean {
  if (typeof window === 'undefined') return false
  if (!window.isSecureContext) return false
  return !!(navigator.mediaDevices && navigator.mediaDevices.getUserMedia)
}

function mediaBlockedHint(): string {
  if (typeof window !== 'undefined' && !window.isSecureContext) {
    return '当前为 HTTP 访问，手机无法开启摄像头/麦克风。可观看电脑端画面，或使用 HTTPS 域名访问后再试。'
  }
  return '本设备不支持摄像头/麦克风，已切换为仅观看模式。'
}

function parseEndTime(endAt: string, durationMinutes: number): number {
  const parsed = new Date(endAt).getTime()
  if (!Number.isFinite(parsed)) {
    return Date.now() + durationMinutes * 60 * 1000
  }
  return parsed
}

type RoomLiveSnapshot = {
  lesson_started?: boolean
  end_at?: string
  actual_minutes?: number
  elapsed_minutes?: number
  mentor_online?: boolean
  user_online?: boolean
}

export function useClassroom(orderId: string, role: 'user' | 'mentor' = 'user') {
  const joining = ref(false)
  const inRoom = ref(false)
  const mockMode = ref(false)
  const receiveOnly = ref(false)
  const mediaHint = ref('')
  const countdown = ref('')
  const lessonStarted = ref(false)
  const elapsedMinutes = ref(0)
  const mentorOnline = ref(false)
  const userOnline = ref(false)
  const mentorName = ref('')
  const errorMsg = ref('')
  const localVideoRef = ref<HTMLDivElement | null>(null)
  const remoteVideoRef = ref<HTMLDivElement | null>(null)

  let client: IAgoraRTCClient | null = null
  let localAudio: IMicrophoneAudioTrack | null = null
  let localVideo: ICameraVideoTrack | null = null
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null
  let countdownTimer: ReturnType<typeof setInterval> | null = null
  let joinInfo: JoinRoomResult | null = null
  let joinedAtElapsed = 0

  const peerOnline = computed(() => (role === 'user' ? mentorOnline.value : userOnline.value))
  const peerLabel = computed(() => (role === 'user' ? '老师' : '学生'))

  const statusHint = computed(() => {
    if (role === 'user') {
      if (!lessonStarted.value) {
        const name = mentorName.value || '老师'
        return `等待 ${name} 进房。计时尚未开始，等待期间不会扣减课时`
      }
      if (!mentorOnline.value) {
        return elapsedMinutes.value > 0
          ? `老师暂时不在课室，课程仍在进行（已上课 ${elapsedMinutes.value} 分钟，剩余 ${countdown.value || '--:--'}）`
          : '老师暂时不在课室，请稍候…'
      }
      if (joinedAtElapsed > 0) {
        return `您晚进了约 ${joinedAtElapsed} 分钟，剩余 ${countdown.value || '--:--'}。请尽快进入学习，下次尽量提前进房`
      }
      return ''
    }

    if (!lessonStarted.value) {
      return '您进房后将自动开始计时。若学生尚未到场，可先等待或联系家长'
    }
    if (!userOnline.value) {
      if (elapsedMinutes.value <= 0) {
        return '等待学生进房中，计时已开始…'
      }
      return `学生尚未进房，已过去 ${elapsedMinutes.value} 分钟。可先准备讲解内容，或联系家长提醒进房`
    }
    return ''
  })

  const timerLabel = computed(() => {
    if (!lessonStarted.value) return '状态'
    if (role === 'user' && joinedAtElapsed > 0) return '剩余时间'
    if (role === 'mentor' && !userOnline.value && elapsedMinutes.value > 0) return '已上课 / 剩余'
    return '剩余时间'
  })

  const timerDisplay = computed(() => {
    if (!lessonStarted.value) return '等待老师进房'
    if (role === 'mentor' && !userOnline.value && elapsedMinutes.value > 0) {
      return `${elapsedMinutes.value} 分 · ${countdown.value || '--:--'}`
    }
    return countdown.value || '--:--'
  })

  function startCountdown(endAt: string, durationMinutes: number) {
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
    const end = parseEndTime(endAt, durationMinutes)
    const tick = () => {
      const diff = end - Date.now()
      if (diff <= 0) {
        countdown.value = '00:00'
        return
      }
      const m = Math.floor(diff / 60000)
      const s = Math.floor((diff % 60000) / 1000)
      countdown.value = `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
    }
    tick()
    countdownTimer = setInterval(tick, 1000)
  }

  function applyLiveSnapshot(res: RoomLiveSnapshot) {
    if (res.mentor_online !== undefined) mentorOnline.value = res.mentor_online
    if (res.user_online !== undefined) userOnline.value = res.user_online

    const elapsed = res.actual_minutes ?? res.elapsed_minutes
    if (elapsed !== undefined) elapsedMinutes.value = elapsed

    if (res.lesson_started && res.end_at) {
      if (!lessonStarted.value) {
        if (role === 'user' && (elapsed ?? 0) > 0) {
          joinedAtElapsed = elapsed ?? 0
        }
        lessonStarted.value = true
        startCountdown(res.end_at, joinInfo?.duration_minutes || 45)
      }
    }
  }

  function handleHeartbeat(res: HeartbeatResult) {
    applyLiveSnapshot(res)
  }

  function startHeartbeat() {
    const beat = () => {
      sendHeartbeat(orderId, role)
        .then(handleHeartbeat)
        .catch(() => {})
    }
    heartbeatTimer = setInterval(beat, 30000)
    beat()
  }

  async function playRemote(user: IAgoraRTCRemoteUser, mediaType: 'video' | 'audio') {
    if (!client) return
    await client.subscribe(user, mediaType)
    if (mediaType === 'video' && remoteVideoRef.value && user.videoTrack) {
      user.videoTrack.play(remoteVideoRef.value)
    }
    if (mediaType === 'audio' && user.audioTrack) {
      user.audioTrack.play()
    }
  }

  function bindRemoteSubscribe() {
    if (!client) return
    client.on('user-published', async (user, mediaType) => {
      await playRemote(user, mediaType)
    })
  }

  async function subscribeExistingRemotes() {
    if (!client) return
    for (const user of client.remoteUsers) {
      if (user.hasVideo) await playRemote(user, 'video')
      if (user.hasAudio) await playRemote(user, 'audio')
    }
  }

  async function joinAgoraLive(info: JoinRoomResult, publishLocal: boolean) {
    client = AgoraRTC.createClient({ mode: 'rtc', codec: 'vp8' })
    bindRemoteSubscribe()
    await client.join(info.app_id, info.channel, info.token, info.uid)
    await subscribeExistingRemotes()

    if (!publishLocal) {
      receiveOnly.value = true
      mediaHint.value = mediaBlockedHint()
      return
    }

    try {
      localAudio = await AgoraRTC.createMicrophoneAudioTrack()
      localVideo = await AgoraRTC.createCameraVideoTrack()
      if (localVideoRef.value) {
        localVideo.play(localVideoRef.value)
      }
      await client.publish([localAudio, localVideo])
    } catch (e) {
      receiveOnly.value = true
      mediaHint.value = mediaBlockedHint()
      const msg = e instanceof Error ? e.message : String(e)
      if (msg.includes('NOT_SUPPORTED') || msg.includes('getUserMedia')) {
        return
      }
      throw e
    }
  }

  async function enter() {
    joining.value = true
    errorMsg.value = ''
    receiveOnly.value = false
    mediaHint.value = ''
    try {
      joinInfo = await joinRoom(orderId, role)
      mentorName.value = joinInfo.mentor_name || ''
      mockMode.value = joinInfo.mock_mode
      applyLiveSnapshot(joinInfo)
      if (!lessonStarted.value) {
        countdown.value = ''
      }
      startHeartbeat()
      inRoom.value = true

      if (joinInfo.mock_mode) {
        return
      }

      const publishLocal = canUseLocalMedia()
      await joinAgoraLive(joinInfo, publishLocal)
    } catch (e) {
      const raw = e instanceof Error ? e.message : getErrorMessage(e)
      if (raw.includes('NOT_SUPPORTED') || raw.includes('getUserMedia')) {
        receiveOnly.value = true
        mediaHint.value = mediaBlockedHint()
        errorMsg.value = ''
        if (!inRoom.value && joinInfo) {
          inRoom.value = true
        }
      } else {
        errorMsg.value = getErrorMessage(e)
      }
    } finally {
      joining.value = false
    }
  }

  async function leave() {
    if (role === 'mentor') {
      try {
        await completeRoom(orderId)
      } catch {
        /* ignore */
      }
    }
    cleanup()
    inRoom.value = false
  }

  function cleanup() {
    if (heartbeatTimer) clearInterval(heartbeatTimer)
    if (countdownTimer) clearInterval(countdownTimer)
    localAudio?.close()
    localVideo?.close()
    client?.leave()
    client = null
  }

  onUnmounted(cleanup)

  return {
    joining,
    inRoom,
    mockMode,
    receiveOnly,
    mediaHint,
    countdown,
    lessonStarted,
    elapsedMinutes,
    mentorOnline,
    userOnline,
    peerOnline,
    peerLabel,
    statusHint,
    timerLabel,
    timerDisplay,
    mentorName,
    errorMsg,
    localVideoRef,
    remoteVideoRef,
    enter,
    leave,
  }
}
