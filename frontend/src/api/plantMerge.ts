import request from '@/utils/request'
import type { PageData } from '@/types/api'
import type { MergeTaskStatus } from '@/constants/merge'

export interface MergePrecheckStats {
  favorites_total: number
  favorites_dedupe: number
  gardens_total: number
  gardens_dedupe: number
  reminders_total: number
  reminders_merge: number
  pests_total: number
}

export interface MergePrecheck {
  keep_id: number
  keep_name: string
  source_ids: number[]
  source_names: string[]
  issues: string[]
  stats: MergePrecheckStats
}

export interface MergeResult {
  favorites_moved: number
  favorites_deduped: number
  gardens_moved: number
  gardens_deduped: number
  reminders_moved: number
  reminders_merged: number
  pests_moved: number
  retired_plants: number
}

export interface MergeTask {
  id: number
  task_key: string
  keep_id: number
  source_ids: number[]
  status: MergeTaskStatus
  result?: MergeResult
  last_error: string
  operator_id: number
  created_at: string
  updated_at: string
}

export interface MergeSubmitResponse {
  task: MergeTask
  message: string
}

export function precheckMerge(payload: { keep_id: number; source_ids: number[] }) {
  return request.post<never, MergePrecheck>('/plants/merge/precheck', payload)
}

export function submitMerge(payload: { keep_id: number; source_ids: number[] }) {
  return request.post<never, MergeSubmitResponse>('/plants/merge', payload)
}

export function getMergeTask(id: number) {
  return request.get<never, MergeTask>(`/plants/merge/tasks/${id}`)
}

export function listMergeTasks(params: { page?: number; page_size?: number }) {
  return request.get<never, PageData<MergeTask>>('/plants/merge/tasks', { params })
}
