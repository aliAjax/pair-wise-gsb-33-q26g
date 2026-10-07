import request from '@/utils/request'
import type { UserGarden } from '@/types/api'

export function listGardens() {
  return request.get<never, UserGarden[]>('/gardens')
}

export function addGarden(payload: { plant_species_id: number; nickname?: string; owned_since?: string; location?: string }) {
  return request.post<never, UserGarden>('/gardens', payload)
}

export function bindReminder(id: number, careReminderId: number) {
  return request.put<never, UserGarden>(`/gardens/${id}/reminder`, { care_reminder_id: careReminderId })
}

export function removeGarden(id: number) {
  return request.delete<never, { removed: boolean }>(`/gardens/${id}`)
}
