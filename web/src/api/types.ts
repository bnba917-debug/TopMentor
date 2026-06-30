export interface ApiResponse<T> {
  code: number
  msg: string
  data?: T
}

export interface User {
  id: number
  phone: string
  child_name: string
  child_age: number
  english_level: string
  available_lessons: number
  locked_lessons: number
}

export interface LoginResult {
  token: string
  user: User
  mentor?: import('./mentor-types').MentorProfile
}

export interface SmsSendResult {
  expires_in: number
  debug_code?: string
}

export interface Mentor {
  id: number
  real_name: string
  school_name: string
  major: string
  gender: string
  english_score: string
  avatar_url?: string
  bio?: string
  intro_video_url: string
  tags: string[]
}

export interface MentorListResult {
  list: Mentor[]
  total: number
}

export interface UpdateProfilePayload {
  child_name: string
  child_age: number
  english_level: 'beginner' | 'intermediate' | 'advanced'
}

export type {
  LessonPackage,
  PaymentChannels,
  RechargeResult,
  LessonBalance,
  MentorSlot,
  BookingResult,
} from './payment-types'

export type {
  CourseOrderDetail,
  JoinRoomResult,
  HeartbeatResult,
  CompleteRoomResult,
} from './room-types'

export type {
  MentorProfile,
  MentorPortalProfile,
  MentorApplyStatus,
  SubmitMentorApplyPayload,
  UpdateMentorProfilePayload,
  UploadResult,
  MentorOrderDetail,
  GrowthReport,
  SubmitReportPayload,
  WalletSummary,
  WalletTransaction,
  WithdrawResult,
  SlotToggle,
} from './mentor-types'
