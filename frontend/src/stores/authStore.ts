import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/user'
import type { UserInfo } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('gbplantwiki_token') || '')
  const user = ref<UserInfo | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(username: string, password: string) {
    const res = await apiLogin({ username, password })
    token.value = res.token
    user.value = res.user
    localStorage.setItem('gbplantwiki_token', res.token)
  }

  async function register(payload: { username: string; email: string; password: string; nickname?: string }) {
    const res = await apiRegister(payload)
    token.value = res.token
    user.value = res.user
    localStorage.setItem('gbplantwiki_token', res.token)
  }

  async function fetchProfile() {
    if (!token.value) return
    user.value = await getProfile()
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('gbplantwiki_token')
  }

  return { token, user, isLoggedIn, isAdmin, login, register, fetchProfile, logout }
})
