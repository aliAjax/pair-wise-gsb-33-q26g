<template>
  <el-card class="question-item" shadow="hover">
    <div class="title" @click="open">{{ question.title }}</div>
    <div class="content">{{ question.content }}</div>
    <div class="meta">{{ formatDateTime(question.created_at) }} · {{ question.status === 'closed' ? '已结题' : '待解答' }}</div>
    <el-button size="small" type="primary" @click="open">查看 / 回答</el-button>
  </el-card>
</template>

<script setup lang="ts">
import type { Question } from '@/types/api'
import { formatDateTime } from '@/utils/dateFormat'

const props = defineProps<{ question: Question }>()
const emit = defineEmits<{ (e: 'open', id: number): void }>()
function open() {
  emit('open', props.question.id)
}
</script>

<style scoped>
.question-item { margin-bottom: 12px; }
.title { font-weight: 700; cursor: pointer; }
.title:hover { color: #3c8d5c; }
.content { color: #666; margin: 6px 0; }
.meta { color: #999; font-size: 12px; margin-bottom: 8px; }
</style>
