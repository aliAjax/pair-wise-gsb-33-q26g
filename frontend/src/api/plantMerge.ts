import request from '@/utils/request'
import type { PageData } from '@/types/api'

export interface MergeStats {
  favorites_moved: number
  favorites_deduped: number
  gardens_moved: number
  gardens_deduped: number
  reminders_moved: number
  reminders_merged: number
  pests_moved: number
}

export interface MergePrecheckItem {
  source_id: number
  source_name: string
  favorites_total: number
  favorites_dedup: number
  gardens_total: number
  gardens_dedup: number
  reminders_total: number
  reminders_merge: number
  pests_total: number
  already_merged: boolean
}

export interface MergePrecheckResult {
  keep_id: number
  keep_name: string
  items: MergePrecheckItem[]
}

export type MergeTaskStatus = 'pending' | 'running' | 'done' | 'failed'

export interface PlantMergeTask {
  id: number
  keep_id: number
  source_id: number
  keep_name: string
  source_name: string
  status: MergeTaskStatus
  stage: string
  stats: MergeStats
  last_error: string
  operator_id: number
  created_at: string
  finished_at: string | null
}

export const MergeStatusMap: Record<MergeTaskStatus, string> = {
  pending: '等待中',
  running: '执行中',
  done: '已完成',
  failed: '失败',
}

export const MergeStageMap: Record<string, string> = {
  favorites: '收藏迁移',
  gardens: '我的花园迁移',
  reminders: '养护提醒归并',
  pests: '病虫害迁移',
  retire: '旧卡退出',
}

export function precheckPlantMerge(payload: { keep_id: number; source_ids: number[] }) {
  return request.post<never, MergePrecheckResult>('/plants/merge/precheck', payload)
}

export function mergePlants(payload: { keep_id: number; source_ids: number[] }) {
  return request.post<never, PlantMergeTask[]>('/plants/merge', payload)
}

export function listMergeTasks(params: { page?: number; page_size?: number }) {
  return request.get<never, PageData<PlantMergeTask>>('/plants/merges', { params })
}

export function getMergeTask(id: number) {
  return request.get<never, PlantMergeTask>(`/plants/merges/${id}`)
}
