import { useAuthStore } from '@/stores/authStore'

// Guard helpers used by router/index.ts and RoleGuard.vue.
export function requireLogin(): boolean {
  const auth = useAuthStore()
  return !!auth.token
}

export function hasRole(roles: string[]): boolean {
  const auth = useAuthStore()
  if (!auth.user) return false
  return roles.length === 0 || roles.includes(auth.user.role)
}
