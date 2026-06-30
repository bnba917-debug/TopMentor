import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { MentorProfile, User } from '@/api/types'
import { smsLogin, sendSms, updateProfile } from '@/api/client'

const TOKEN_KEY = 'tm_token'
const USER_KEY = 'tm_user'
const MENTOR_KEY = 'tm_mentor'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const user = ref<User | null>(readStoredUser())
  const mentor = ref<MentorProfile | null>(readStoredMentor())

  const isLoggedIn = computed(() => !!token.value)
  const isMentor = computed(() => !!mentor.value)

  function persist(loginToken: string, loginUser: User, loginMentor?: MentorProfile | null) {
    token.value = loginToken
    user.value = loginUser
    localStorage.setItem(TOKEN_KEY, loginToken)
    localStorage.setItem(USER_KEY, JSON.stringify(loginUser))

    if (loginMentor) {
      mentor.value = loginMentor
      localStorage.setItem(MENTOR_KEY, JSON.stringify(loginMentor))
    } else {
      mentor.value = null
      localStorage.removeItem(MENTOR_KEY)
    }
  }

  async function requestCode(phone: string) {
    return sendSms(phone)
  }

  async function login(phone: string, code: string) {
    const result = await smsLogin(phone, code)
    persist(result.token, result.user, result.mentor ?? null)
    return result
  }

  async function saveProfile(payload: {
    child_name: string
    child_age: number
    english_level: 'beginner' | 'intermediate' | 'advanced'
  }) {
    const updated = await updateProfile(payload)
    if (token.value) {
      persist(token.value, updated, mentor.value)
    }
    return updated
  }

  function logout() {
    token.value = ''
    user.value = null
    mentor.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(MENTOR_KEY)
  }

  function syncMentorSummary(profile: {
    id: number
    real_name: string
    school_name: string
    is_verified: number
    avatar_url?: string
  }) {
    if (!token.value) return
    const next: MentorProfile = {
      id: profile.id,
      real_name: profile.real_name,
      school_name: profile.school_name,
      is_verified: profile.is_verified,
      avatar_url: profile.avatar_url,
    }
    mentor.value = next
    localStorage.setItem(MENTOR_KEY, JSON.stringify(next))
  }

  function setAvailableLessons(count: number) {
    if (!user.value) return
    user.value = { ...user.value, available_lessons: count }
    localStorage.setItem(USER_KEY, JSON.stringify(user.value))
  }

  return {
    token,
    user,
    mentor,
    isLoggedIn,
    isMentor,
    requestCode,
    login,
    saveProfile,
    logout,
    syncMentorSummary,
    setAvailableLessons,
  }
})

function readStoredUser(): User | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

function readStoredMentor(): MentorProfile | null {
  const raw = localStorage.getItem(MENTOR_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as MentorProfile
  } catch {
    return null
  }
}
