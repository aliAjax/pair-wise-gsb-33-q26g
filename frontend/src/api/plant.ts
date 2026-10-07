import request from '@/utils/request'
import type { PlantSpecies } from '@/constants/plant'
import type { PageData } from '@/types/api'

export function listPlants(params: { page?: number; page_size?: number; type?: string; family?: string; keyword?: string }) {
  return request.get<never, PageData<PlantSpecies>>('/plants', { params })
}

export function getPlant(id: number | string) {
  return request.get<never, PlantSpecies>(`/plants/${id}`)
}

export function createPlant(payload: Partial<PlantSpecies>) {
  return request.post<never, PlantSpecies>('/plants', payload)
}

export function updatePlant(id: number, payload: Partial<PlantSpecies>) {
  return request.put<never, PlantSpecies>(`/plants/${id}`, payload)
}

export function deletePlant(id: number) {
  return request.delete<never, { deleted: boolean }>(`/plants/${id}`)
}
