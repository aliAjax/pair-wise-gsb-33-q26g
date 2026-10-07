import request from '@/utils/request'
import type { CareReminder } from '@/types/api'

export function listReminders(status?: string) {
  return request.get<never, CareReminder[]>('/reminders', { params: { status } })
}

export function listRemindersByMonth(year: number, month: number) {
  return request.get<never, CareReminder[]>('/reminders/calendar', { params: { year, month } })
}

export function createReminder(payload: { plant_species_id?: number; task_title: string; remind_date: string; frequency?: string }) {
  return request.post<never, CareReminder>('/reminders', payload)
}

export function updateReminderStatus(id: number, status: string) {
  return request.put<never, CareReminder>(`/reminders/${id}/status`, { status })
}

export function deleteReminder(id: number) {
  return request.delete<never, { deleted: boolean }>(`/reminders/${id}`)
}
