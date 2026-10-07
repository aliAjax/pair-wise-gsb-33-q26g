<template>
  <div class="page">
    <h1>季节养护日历</h1>
    <el-row :gutter="16">
      <el-col :xs="24" :md="14">
        <el-timeline>
          <el-timeline-item v-for="t in seasonTasks" :key="t.month" :timestamp="t.label" :type="t.month === currentMonth ? 'primary' : ''">
            {{ t.task }}
            <el-tag v-if="t.month === currentMonth" size="small" type="success">本月</el-tag>
          </el-timeline-item>
        </el-timeline>
      </el-col>
      <el-col :xs="24" :md="10">
        <el-card>
          <template #header>本月养护提醒</template>
          <el-form inline>
            <el-form-item label="任务"><el-input v-model="form.task_title" placeholder="如：给月季施肥" /></el-form-item>
            <el-form-item label="日期"><el-date-picker v-model="form.remind_date" type="date" value-format="YYYY-MM-DD" /></el-form-item>
            <el-form-item><el-button type="primary" @click="create">创建提醒</el-button></el-form-item>
          </el-form>
          <ReminderList :reminders="reminders" @done="markDone" @remove="remove" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import ReminderList from '@/components/common/ReminderList.vue'
import { useReminderStore } from '@/stores/reminderStore'
import { getSeasonTasks } from '@/utils/season'
import type { CareReminder } from '@/types/api'

const store = useReminderStore()
const seasonTasks = getSeasonTasks()
const currentMonth = new Date().getMonth() + 1
const reminders = ref<CareReminder[]>([])
const form = reactive({ task_title: '', remind_date: '' })

onMounted(async () => {
  await store.load()
  reminders.value = store.reminders
})

async function create() {
  if (!form.task_title || !form.remind_date) {
    ElMessage.warning('请填写任务与日期')
    return
  }
  await store.create({ task_title: form.task_title, remind_date: form.remind_date })
  reminders.value = store.reminders
  form.task_title = ''
  form.remind_date = ''
  ElMessage.success('养护提醒已创建')
}
async function markDone(id: number) {
  await store.setStatus(id, 'done')
  reminders.value = store.reminders
}
async function remove(id: number) {
  const { deleteReminder } = await import('@/api/reminder')
  await deleteReminder(id)
  await store.load()
  reminders.value = store.reminders
  ElMessage.success('已删除提醒')
}
</script>

<style scoped>
.page { max-width: 1200px; margin: 0 auto; }
</style>
