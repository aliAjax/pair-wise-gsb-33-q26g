import request from '@/utils/request'
import type { UserInfo } from '@/types/api'

export interface LoginResult {
  token: string
  user: UserInfo
}

export function register(payload: { username: string; email: string; password: string; nickname?: string }) {
  return request.post<never, LoginResult>('/users/register', payload)
}

export function login(payload: { username: string; password: string }) {
  return request.post<never, LoginResult>('/users/login', payload)
}

export function getProfile() {
  return request.get<never, UserInfo>('/users/me')
}

export function updateProfile(payload: { nickname?: string; bio?: string; avatar?: string }) {
  return request.put<never, UserInfo>('/users/me', payload)
}

export function listUsers(params: { page?: number; page_size?: number }) {
  return request.get<never, { list: UserInfo[]; total: number }>('/users', { params })
}
