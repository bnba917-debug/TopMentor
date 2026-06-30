import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { adminLogin } from '@/api/client'

const TOKEN_KEY = 'tm_admin_token'
const USER_KEY = 'tm_admin_user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) || '')
  const username = ref(localStorage.getItem(USER_KEY) || '')

  const isLoggedIn = computed(() => !!token.value)

  async function login(user: string, password: string) {
    const result = await adminLogin(user, password)
    token.value = result.token
    username.value = result.username
    localStorage.setItem(TOKEN_KEY, result.token)
    localStorage.setItem(USER_KEY, result.username)
    return result
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, username, isLoggedIn, login, logout }
})
