import request from '@/utils/request'
import type { DiseasePest } from '@/types/api'
import type { PageData } from '@/types/api'

export function listPests(params: { page?: number; page_size?: number; keyword?: string; plant_species_id?: number }) {
  return request.get<never, PageData<DiseasePest>>('/pests', { params })
}

export function getPest(id: number | string) {
  return request.get<never, DiseasePest>(`/pests/${id}`)
}
