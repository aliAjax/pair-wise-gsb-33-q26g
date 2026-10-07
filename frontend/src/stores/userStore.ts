import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getProfile, updateProfile } from '@/api/user'
import type { UserInfo } from '@/types/api'

export const useUserStore = defineStore('user', () => {
  const profile = ref<UserInfo | null>(null)

  async function load() {
    profile.value = await getProfile()
  }

  async function save(payload: { nickname?: string; bio?: string; avatar?: string }) {
    profile.value = await updateProfile(payload)
  }

  return { profile, load, save }
})
