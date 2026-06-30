import axios from 'axios'
import type {
  ApiResponse,
  BookingResult,
  CompleteRoomResult,
  CourseOrderDetail,
  GrowthReport,
  HeartbeatResult,
  JoinRoomResult,
  LessonBalance,
  LessonPackage,
  LoginResult,
  Mentor,
  MentorListResult,
  MentorOrderDetail,
  MentorPortalProfile,
  MentorApplyStatus,
  MentorSlot,
  PaymentChannels,
  RechargeResult,
  SubmitMentorApplyPayload,
  SlotToggle,
  SubmitReportPayload,
  SmsSendResult,
  UpdateMentorProfilePayload,
  UpdateProfilePayload,
  UploadResult,
  User,
  WalletSummary,
  WithdrawResult,
} from './types'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('tm_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

function unwrap<T>(res: ApiResponse<T>): T {
  if (res.code !== 0 || res.data === undefined) {
    throw new Error(res.msg || '请求失败')
  }
  return res.data
}

export async function sendSms(phone: string): Promise<SmsSendResult> {
  const { data } = await http.post<ApiResponse<SmsSendResult>>('/auth/sms/send', { phone })
  return unwrap(data)
}

export async function smsLogin(phone: string, code: string): Promise<LoginResult> {
  const { data } = await http.post<ApiResponse<LoginResult>>('/auth/sms/login', { phone, code })
  return unwrap(data)
}

export async function updateProfile(payload: UpdateProfilePayload): Promise<User> {
  const { data } = await http.put<ApiResponse<User>>('/users/profile', payload)
  return unwrap(data)
}

export async function fetchMentors(params?: {
  school?: string
  gender?: string
  tag?: string
  page?: number
}): Promise<MentorListResult> {
  const { data } = await http.get<ApiResponse<MentorListResult>>('/mentors', { params })
  return unwrap(data)
}

export async function fetchMentor(id: number): Promise<Mentor> {
  const { data } = await http.get<ApiResponse<Mentor>>(`/mentors/${id}`)
  return unwrap(data)
}

export async function fetchPackages(): Promise<LessonPackage[]> {
  const { data } = await http.get<ApiResponse<LessonPackage[]>>('/packages')
  return unwrap(data)
}

export async function fetchPaymentChannels(): Promise<PaymentChannels> {
  const { data } = await http.get<ApiResponse<PaymentChannels>>('/payment/channels')
  return unwrap(data)
}

export async function createRecharge(packageId: number, channel: string): Promise<RechargeResult> {
  const { data } = await http.post<ApiResponse<RechargeResult>>('/recharge', {
    package_id: packageId,
    channel,
  })
  return unwrap(data)
}

export async function fetchLessonBalance(): Promise<LessonBalance> {
  const { data } = await http.get<ApiResponse<LessonBalance>>('/users/lessons')
  return unwrap(data)
}

export async function fetchMentorSlots(mentorId: number): Promise<MentorSlot[]> {
  const { data } = await http.get<ApiResponse<MentorSlot[]>>(`/mentors/${mentorId}/slots`)
  return unwrap(data)
}

export async function createBooking(slotId: number): Promise<BookingResult> {
  const { data } = await http.post<ApiResponse<BookingResult>>('/bookings', { slot_id: slotId })
  return unwrap(data)
}

export async function fetchUserOrders(): Promise<CourseOrderDetail[]> {
  const { data } = await http.get<ApiResponse<CourseOrderDetail[]>>('/users/orders')
  return unwrap(data)
}

export async function fetchOrder(orderId: string): Promise<CourseOrderDetail> {
  const { data } = await http.get<ApiResponse<CourseOrderDetail>>(`/orders/${orderId}`)
  return unwrap(data)
}

export async function joinRoom(orderId: string, role: 'user' | 'mentor'): Promise<JoinRoomResult> {
  const { data } = await http.post<ApiResponse<JoinRoomResult>>(`/rooms/${orderId}/join`, { role })
  return unwrap(data)
}

export async function sendHeartbeat(orderId: string, role: 'user' | 'mentor'): Promise<HeartbeatResult> {
  const { data } = await http.post<ApiResponse<HeartbeatResult>>(
    `/rooms/${orderId}/heartbeat?role=${role}`,
  )
  return unwrap(data)
}

export async function completeRoom(orderId: string): Promise<CompleteRoomResult> {
  const { data } = await http.post<ApiResponse<CompleteRoomResult>>(`/rooms/${orderId}/complete`)
  return unwrap(data)
}

export async function fetchMentorOrders(): Promise<MentorOrderDetail[]> {
  const { data } = await http.get<ApiResponse<MentorOrderDetail[]>>('/mentor/orders')
  return unwrap(data)
}

export async function fetchMentorPortalSlots(params?: {
  from?: string
  to?: string
}): Promise<MentorSlot[]> {
  const { data } = await http.get<ApiResponse<MentorSlot[]>>('/mentor/slots', { params })
  return unwrap(data)
}

export async function setMentorSlots(slots: SlotToggle[]): Promise<{ updated: number }> {
  const { data } = await http.put<ApiResponse<{ updated: number }>>('/mentor/slots', { slots })
  return unwrap(data)
}

export async function submitMentorReport(payload: SubmitReportPayload): Promise<GrowthReport> {
  const { data } = await http.post<ApiResponse<GrowthReport>>('/mentor/reports', payload)
  return unwrap(data)
}

export async function fetchMentorWallet(): Promise<WalletSummary> {
  const { data } = await http.get<ApiResponse<WalletSummary>>('/mentor/wallet')
  return unwrap(data)
}

export async function mentorWithdraw(amountCents: number): Promise<WithdrawResult> {
  const { data } = await http.post<ApiResponse<WithdrawResult>>('/mentor/withdraw', {
    amount_cents: amountCents,
  })
  return unwrap(data)
}

export async function fetchMentorProfile(): Promise<MentorPortalProfile> {
  const { data } = await http.get<ApiResponse<MentorPortalProfile>>('/mentor/profile')
  return unwrap(data)
}

export async function updateMentorProfile(
  payload: UpdateMentorProfilePayload,
): Promise<MentorPortalProfile> {
  const { data } = await http.put<ApiResponse<MentorPortalProfile>>('/mentor/profile', payload)
  return unwrap(data)
}

export async function uploadMentorFile(
  kind: 'avatar' | 'intro_video',
  file: File,
  onProgress?: (pct: number) => void,
): Promise<UploadResult> {
  const form = new FormData()
  form.append('kind', kind)
  form.append('file', file)
  const { data } = await http.post<ApiResponse<UploadResult>>('/mentor/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (onProgress && e.total) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
  return unwrap(data)
}

export function mediaUrl(url?: string): string {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  return url
}

export async function fetchMentorApplyStatus(): Promise<MentorApplyStatus> {
  const { data } = await http.get<ApiResponse<MentorApplyStatus>>('/mentor/apply/status')
  return unwrap(data)
}

export async function submitMentorApply(payload: SubmitMentorApplyPayload): Promise<MentorApplyStatus> {
  const { data } = await http.post<ApiResponse<MentorApplyStatus>>('/mentor/apply', payload)
  return unwrap(data)
}

export async function uploadApplyFile(
  kind: 'avatar' | 'intro_video' | 'id_card' | 'student_card' | 'english_proof',
  file: File,
  onProgress?: (pct: number) => void,
): Promise<UploadResult> {
  const form = new FormData()
  form.append('kind', kind)
  form.append('file', file)
  const { data } = await http.post<ApiResponse<UploadResult>>('/mentor/apply/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: (e) => {
      if (onProgress && e.total) {
        onProgress(Math.round((e.loaded / e.total) * 100))
      }
    },
  })
  return unwrap(data)
}

export async function fetchGrowthReport(orderId: string): Promise<GrowthReport> {
  const { data } = await http.get<ApiResponse<GrowthReport>>(`/reports/${orderId}`)
  return unwrap(data)
}

export function formatYuan(cents: number): string {
  return (cents / 100).toFixed(2)
}

export function getErrorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const body = err.response?.data as ApiResponse<unknown> | undefined
    if (body?.msg) return body.msg
    if (err.response?.status === 404) {
      return '接口不存在，请确认后端已重启到最新版本'
    }
  }
  if (err instanceof Error) return err.message
  return '网络错误'
}
