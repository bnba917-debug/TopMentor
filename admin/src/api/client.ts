import axios from 'axios'
import type {
  AdminLoginResult,
  ApiResponse,
  Courseware,
  FinanceSummary,
  PendingMentorApplication,
} from './types'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('tm_admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (res) => res,
  (err) => {
    if (axios.isAxiosError(err) && err.response?.status === 401) {
      localStorage.removeItem('tm_admin_token')
      localStorage.removeItem('tm_admin_user')
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  },
)

function unwrap<T>(res: ApiResponse<T>): T {
  if (res.code !== 0 || res.data === undefined) {
    throw new Error(res.msg || '请求失败')
  }
  return res.data
}

export async function adminLogin(username: string, password: string): Promise<AdminLoginResult> {
  const { data } = await http.post<ApiResponse<AdminLoginResult>>('/admin/auth/login', {
    username,
    password,
  })
  return unwrap(data)
}

export async function fetchPendingMentors(): Promise<PendingMentorApplication[]> {
  const { data } = await http.get<ApiResponse<PendingMentorApplication[]>>('/admin/mentors/pending')
  return unwrap(data)
}

export async function reviewMentor(
  mentorId: number,
  payload: { action: 'approve' | 'reject'; reject_reason?: string },
): Promise<void> {
  const { data } = await http.put<ApiResponse<unknown>>(`/admin/mentors/${mentorId}/review`, payload)
  unwrap(data)
}

export async function fetchCourseware(): Promise<Courseware[]> {
  const { data } = await http.get<ApiResponse<Courseware[]>>('/admin/courseware')
  return unwrap(data)
}

export async function createCourseware(payload: {
  title: string
  cover_url?: string
  content_url: string
  sort_order?: number
  is_active?: boolean
}): Promise<Courseware> {
  const { data } = await http.post<ApiResponse<Courseware>>('/admin/courseware', payload)
  return unwrap(data)
}

export async function updateCourseware(
  id: number,
  payload: Partial<{
    title: string
    cover_url: string
    content_url: string
    sort_order: number
    is_active: boolean
  }>,
): Promise<Courseware> {
  const { data } = await http.put<ApiResponse<Courseware>>(`/admin/courseware/${id}`, payload)
  return unwrap(data)
}

export async function deleteCourseware(id: number): Promise<void> {
  const { data } = await http.delete<ApiResponse<unknown>>(`/admin/courseware/${id}`)
  unwrap(data)
}

export async function fetchFinanceSummary(): Promise<FinanceSummary> {
  const { data } = await http.get<ApiResponse<FinanceSummary>>('/admin/finance/summary')
  return unwrap(data)
}

export function getErrorMessage(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const body = err.response?.data as ApiResponse<unknown> | undefined
    if (body?.msg) return body.msg
  }
  if (err instanceof Error) return err.message
  return '网络错误'
}

export function maskUrl(url: string): string {
  if (!url) return '-'
  if (url.length <= 12) return '***'
  return `${url.slice(0, 8)}***${url.slice(-4)}`
}
