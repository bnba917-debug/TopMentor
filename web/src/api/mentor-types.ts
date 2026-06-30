import type { CourseOrderDetail } from './room-types'

export interface MentorProfile {
  id: number
  real_name: string
  school_name: string
  is_verified: number
  avatar_url?: string
}

export interface MentorPortalProfile {
  id: number
  phone?: string
  real_name: string
  school_name: string
  major: string
  gender: string
  english_score: string
  avatar_url: string
  bio: string
  intro_video_url: string
  tags: string[]
  is_verified: number
}

export interface UpdateMentorProfilePayload {
  real_name: string
  school_name: string
  major: string
  gender?: string
  english_score?: string
  bio?: string
  tags?: string[]
  avatar_url?: string
  intro_video_url?: string
}

export interface UploadResult {
  url: string
}

export type MentorApplyStatusValue = 'none' | 'pending' | 'rejected' | 'approved'

export interface MentorApplyDraft {
  real_name: string
  school_name: string
  major: string
  gender: string
  english_score: string
  bio: string
  avatar_url: string
  intro_video_url: string
  tags: string[]
  id_card_url: string
  student_card_url: string
  english_proof_url?: string
}

export interface MentorApplyStatus {
  status: MentorApplyStatusValue
  mentor_id?: number
  application_id?: number
  reject_reason?: string
  applied_at?: string
  profile?: MentorApplyDraft
}

export interface SubmitMentorApplyPayload {
  real_name: string
  school_name: string
  major: string
  gender?: string
  english_score?: string
  bio?: string
  tags?: string[]
  avatar_url: string
  intro_video_url: string
  id_card_url: string
  student_card_url: string
  english_proof_url?: string
}

export interface MentorOrderDetail extends CourseOrderDetail {
  child_name?: string
  user_phone?: string
  has_report?: boolean
}

export interface GrowthReport {
  id: number
  order_id: string
  mentor_id: number
  user_id: number
  speaking_score: number
  confidence_score: number
  vocabulary_score: number
  comment: string
  mentor_name?: string
  child_name?: string
  created_at: string
}

export interface SubmitReportPayload {
  order_id: string
  speaking_score: number
  confidence_score: number
  vocabulary_score: number
  comment: string
}

export interface WalletTransaction {
  id: number
  amount: number
  type: string
  balance_after: number
  remark?: string
  created_at: string
}

export interface WalletSummary {
  balance: number
  transactions: WalletTransaction[]
}

export interface WithdrawResult {
  balance: number
  mock_paid?: boolean
}

export interface SlotToggle {
  slot_date: string
  start_time: string
  available: boolean
}
