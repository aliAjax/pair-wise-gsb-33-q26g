<template>
  <el-table :data="reminders" stripe empty-text="暂无养护提醒">
    <el-table-column prop="task_title" label="任务" min-width="180" />
    <el-table-column label="提醒日期" width="120">
      <template #default="{ row }">{{ formatDate(row.remind_date) }}</template>
    </el-table-column>
    <el-table-column label="频率" width="100">
      <template #default="{ row }">{{ frequencyText(row.frequency) }}</template>
    </el-table-column>
    <el-table-column label="状态" width="110">
      <template #default="{ row }">
        <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="160">
      <template #default="{ row }">
        <el-button v-if="row.status !== 'done'" size="small" type="success" @click="$emit('done', row.id)">完成</el-button>
        <el-button size="small" type="danger" @click="$emit('remove', row.id)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { CareReminder } from '@/types/api'
import { formatDate } from '@/utils/dateFormat'

defineProps<{ reminders: CareReminder[] }>()
defineEmits<{ (e: 'done', id: number): void; (e: 'remove', id: number): void }>()

function statusText(s: string): string {
  return s === 'pending' ? '待处理' : s === 'done' ? '已完成' : '已逾期'
}
function statusType(s: string): 'warning' | 'success' | 'danger' {
  return s === 'pending' ? 'warning' : s === 'done' ? 'success' : 'danger'
}
function frequencyText(f: string): string {
  const map: Record<string, string> = { daily: '每日', weekly: '每周', monthly: '每月', yearly: '每年' }
  return map[f] || f || '-'
}
</script>
