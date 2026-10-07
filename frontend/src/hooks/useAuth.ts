import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'

export function useAuth() {
  const auth = useAuthStore()
  const isLoggedIn = computed(() => auth.isLoggedIn)
  const isAdmin = computed(() => auth.isAdmin)
  const user = computed(() => auth.user)
  return { auth, isLoggedIn, isAdmin, user }
}
