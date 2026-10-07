import { computed, ref } from 'vue'
import { listReminders } from '@/api/reminder'
import type { CareReminder } from '@/types/api'

export function useReminderStats() {
  const reminders = ref<CareReminder[]>([])
  const loading = ref(false)

  async function load() {
    loading.value = true
    try {
      reminders.value = await listReminders()
    } finally {
      loading.value = false
    }
  }

  const pending = computed(() => reminders.value.filter((r) => r.status === 'pending').length)
  const done = computed(() => reminders.value.filter((r) => r.status === 'done').length)
  const overdue = computed(() => reminders.value.filter((r) => r.status === 'overdue').length)

  return { reminders, loading, load, pending, done, overdue }
}
