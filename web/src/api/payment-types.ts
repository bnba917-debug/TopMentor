export interface LessonPackage {
  id: number
  name: string
  lesson_count: number
  price_cents: number
  is_trial: boolean
}

export interface PaymentChannels {
  mode: string
  channels: string[]
}

export interface RechargeResult {
  order_id: string
  status: string
  amount_cents: number
  package_name: string
  lesson_count: number
  channel: string
  pay_url?: string
  jsapi_params?: Record<string, string>
  mock_paid?: boolean
  available_lessons?: number
}

export interface LessonBalance {
  available_lessons: number
  locked_lessons: number
}

export interface MentorSlot {
  id: number
  mentor_id: number
  slot_date: string
  start_time: string
  end_time: string
  status: number
}

export interface BookingResult {
  order: {
    id: string
    status: string
    slot_id: number
  }
  available_lessons: number
}
