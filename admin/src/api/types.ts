export interface ApiResponse<T> {
  code: number
  msg: string
  data?: T
}

export interface AdminLoginResult {
  token: string
  username: string
}

export interface PendingMentorApplication {
  mentor_id: number
  application_id: number
  real_name: string
  school_name: string
  major: string
  gender: string
  english_score: string
  intro_video_url: string
  id_card_url: string
  student_card_url: string
  english_proof_url?: string
  applied_at: string
}

export interface Courseware {
  id: number
  title: string
  cover_url?: string
  content_url: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface FinanceSummary {
  total_recharge_yuan: number
  total_withdraw_yuan: number
  total_mentor_earn_yuan: number
  unspent_lessons: number
  completed_orders: number
  pending_mentors: number
  active_courseware: number
}
