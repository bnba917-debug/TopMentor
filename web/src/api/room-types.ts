export interface CourseOrderDetail {
  id: string
  user_id: number
  mentor_id: number
  slot_id: number
  status: string
  agora_channel_name?: string
  actual_minutes: number
  created_at: string
  mentor_name?: string
  slot_date?: string
  start_time?: string
  end_time?: string
}

export interface JoinRoomResult {
  app_id: string
  channel: string
  token: string
  uid: number
  role: string
  mock_mode: boolean
  order_id: string
  order_status: string
  end_at: string
  duration_minutes: number
  lesson_started: boolean
  started_at?: string
  elapsed_minutes: number
  mentor_online: boolean
  user_online: boolean
  mentor_name?: string
}

export interface HeartbeatResult {
  ok: boolean
  actual_minutes: number
  lesson_started: boolean
  end_at?: string
  started_at?: string
  mentor_online: boolean
  user_online: boolean
}

export interface CompleteRoomResult {
  order_id: string
  status: string
  actual_minutes: number
}
